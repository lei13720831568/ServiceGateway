package ActiveHttpProxy //Reverse Porxy

import (
	log "RollLoger"
	"WorkPool"
	"encoding/json"
	//	"errors"
	"net"
	"net/url"
	"time"
	//	"fmt"
	"net/http"
	//	"runtime"
	"strconv"
	"sync"
	//"time"
	"fmt"
	"io"
)

type ArRouteLoad struct {
	MaxVer int64
	Routes []*ArRoute
}

type ArRoute struct {
	PublishID   int
	Name        string
	ReqUrl      string
	ProxyToUrl  string
	IpList      string
	SecretType  string
	Encrypt     string
	MaxConnects int
	TimeOut     int64 //单位毫秒
	Ver         int64
	Status      int
	ProxyWorks  *WorkPool.WPool
}

func (ar *ArRoute) ToJson() string {
	r, err := json.Marshal(ar)
	if err != nil {
		log.Error("Marshal ArRoute To json failed;")
	}
	return string(r)
}

func (ar *ArRoute) InitProxyWorks() error {
	ar.ProxyWorks = WorkPool.NewWorkPool(ar.Name)
	return ar.ProxyWorks.SetMax(ar.MaxConnects)
}

func (ar *ArRoute) CloseProxyWorks() {
	ar.ProxyWorks.SetMax(0)
}

type ArRouteMap struct {
	mu     sync.Mutex
	Routes map[string]*ArRoute
	MaxVer int64
}

func NewArRouteMap() *ArRouteMap {
	service_routes := &ArRouteMap{}
	service_routes.Routes = make(map[string]*ArRoute)
	service_routes.MaxVer = 0
	return service_routes
}

//查找路由map找到对应的目的地址
func (arm *ArRouteMap) MatchRoute(url string) (*ArRoute, bool) {
	arm.mu.Lock()
	defer arm.mu.Unlock()
	r, ok := arm.Routes[url]
	if ok {
		return r, true
	} else {
		return nil, false
	}
}

func (arm *ArRouteMap) FindRouteByID(publishID int) (*ArRoute, bool) {
	for _, arr := range arm.Routes {
		if arr.PublishID == publishID {
			return arr, true
		}
	}
	return nil, false
}

//从json加载路由信息，返回值代表更新的路由数量
func (arm *ArRouteMap) RoadRoute(newAroute *ArRouteLoad) (result int) {
	arm.mu.Lock()
	defer arm.mu.Unlock()

	if newAroute.MaxVer < arm.MaxVer { //版本低于当前版本
		return result
	}

	for _, newroute := range newAroute.Routes { //检查路由信息并进行替换

		oldroute, ok := arm.FindRouteByID(newroute.PublishID)
		if ok {
			if newroute.Ver > oldroute.Ver {
				if newroute.Status != 0 {
					oldroute.CloseProxyWorks()          //关闭服务
					delete(arm.Routes, oldroute.ReqUrl) //已禁用从路由表内删除
					result++
					log.Info("update route to close ", oldroute.ToJson())
					continue
				}

				//更新属性
				oldroute.Encrypt = newroute.Encrypt
				oldroute.IpList = newroute.IpList
				oldroute.Name = newroute.Name
				oldroute.ProxyToUrl = newroute.ProxyToUrl
				oldroute.PublishID = newroute.PublishID
				//发布地址发生了改变
				if oldroute.ReqUrl != newroute.ReqUrl {
					delete(arm.Routes, oldroute.ReqUrl)    //从路由表内删除
					arm.Routes[newroute.ReqUrl] = oldroute //使用新的key添加
				}

				oldroute.ReqUrl = newroute.ReqUrl
				oldroute.SecretType = newroute.SecretType
				oldroute.TimeOut = newroute.TimeOut
				oldroute.Ver = newroute.Ver
				//调整连接数
				err := oldroute.ProxyWorks.SetMax(newroute.MaxConnects)
				if err == nil {
					oldroute.MaxConnects = newroute.MaxConnects
				} else {
					log.Error("连接调整失败old:", strconv.Itoa(oldroute.MaxConnects), " new:", strconv.Itoa(newroute.MaxConnects))
				}
				result++
				log.Info("update route to ", oldroute.ToJson())
			}

		} else {
			if newroute.Status != 0 {
				continue
			}
			newroute.ProxyWorks = WorkPool.NewWorkPool(newroute.Name)
			err := newroute.InitProxyWorks()
			if err != nil {
				//连接可能过大
				newroute.MaxConnects = 0
				log.Error("初始化连接失败 max:", strconv.Itoa(newroute.MaxConnects))
			}

			arm.Routes[newroute.ReqUrl] = newroute

			result++
			log.Info("create route to ", newroute.ToJson())

		}
	}

	return result
}

type ProxyWork struct {
	w         http.ResponseWriter
	r         *http.Request
	Transport *http.Transport
	DestUrl   string
	rspStatus int
}

func NewProxyWork(writer http.ResponseWriter, req *http.Request, desturl string) *ProxyWork {
	pw := &ProxyWork{}
	pw.w = writer
	pw.DestUrl = desturl
	pw.r = req
	pw.Transport = &http.Transport{DisableKeepAlives: false, DisableCompression: false}

	return pw
}

func (pw *ProxyWork) PHandle() error {
	newUrlstr := pw.DestUrl + "?" + pw.r.URL.RawQuery
	newUrl, pe := url.Parse(newUrlstr)
	if pe != nil {
		return pe
	}
	pw.r.URL = newUrl

	resp, err := pw.Transport.RoundTrip(pw.r)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			pw.rspStatus = 408 //超时
		} else {
			pw.rspStatus = 503 //其他错误
		}
		return err
	} else {
		defer resp.Body.Close()
		for k, v := range resp.Header {
			for _, vv := range v {
				pw.w.Header().Add(k, vv)
			}
		}
		pw.rspStatus = resp.StatusCode
		pw.w.WriteHeader(resp.StatusCode)
		io.Copy(pw.w, resp.Body)
	}

	return nil
}

type ArProxy struct {
	port           string
	service_queue  map[int]chan byte
	service_routes *ArRouteMap
	wa             *Watcher
}

func NewArProxy(port string, r RouteReader) *ArProxy {
	arp := &ArProxy{}
	arp.port = port
	arp.service_queue = make(map[int]chan byte)
	arp.service_routes = NewArRouteMap()
	arp.wa = NewWatcher(r)
	return arp
}

func (arp *ArProxy) handleService(w http.ResponseWriter, r *http.Request) {

	log.Debug("receive request ", r.URL.Path)
	rurl := r.URL.Path
	ar, ok := arp.service_routes.MatchRoute(r.URL.Path)
	if ok {
		log.Debug("match request ", r.URL.Path, " to ", ar.ProxyToUrl)
		pw := NewProxyWork(w, r, ar.ProxyToUrl)
		err := ar.ProxyWorks.PutWork(pw, time.Duration(ar.TimeOut/2))
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "work Error: %v", err)
			log.Info("route to err:", ar.Name, "$$", rurl, "$$", ar.ProxyToUrl, "$$", strconv.Itoa(pw.rspStatus))
		} else {
			log.Info("route to :", ar.Name, "$$", rurl, "$$", ar.ProxyToUrl, "$$", strconv.Itoa(pw.rspStatus))
		}
	} else {
		http.ServeFile(w, r, "NotFindService.html")
		log.Info("route Match:", "", "$$", r.URL.Path, "$$", "", "$$", "404")
	}

	//w.Write([]byte("Hello"))

}

func (arp *ArProxy) Start() {
	Arrl, err := arp.wa.reader.Read()

	if err != nil {
		log.Fatal(err)
	}
	rc := arp.service_routes.RoadRoute(Arrl)

	log.Info("init routeInfo done. routes count:", strconv.Itoa(rc))

	go arp.wa.StartWatch(arp)
	http.HandleFunc("/", arp.handleService) //设置访问的路由

	err = http.ListenAndServe(":"+arp.port, nil) //设置监听的端口
	if err != nil {
		log.Fatal(err)
	}
	log.Info("start listen port:", arp.port)

}
