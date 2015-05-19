package ActiveHttpProxy //Reverse Porxy

import (
	log "RollLoger"
	"WorkPool"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	//	"strings"
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
	"StoppableListener"
	//	"container/list"
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
	TimeOut     int64 //连接排队超时单位毫秒
	Ver         int64
	Status      int
	SecKey      string
	MatchMode   string // MatchDir ,MatchFile
	WaitTimeOut int64
	transport   *http.Transport
	ProxyWorks  *WorkPool.WPool
}

func (ar *ArRoute) InitTransport() {
	ar.transport = &http.Transport{DisableKeepAlives: false, DisableCompression: false}
	ar.transport.Dial = ar.dialTimeout
}
func (ar *ArRoute) dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(ar.WaitTimeOut)*time.Millisecond)
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
	mu          sync.Mutex
	CacheRoutes map[string]*ArRoute
	Routes      []*ArRoute
	MaxVer      int64
}

func NewArRouteMap() *ArRouteMap {
	service_routes := &ArRouteMap{}
	service_routes.CacheRoutes = make(map[string]*ArRoute)
	service_routes.Routes = []*ArRoute{}
	service_routes.MaxVer = 0
	return service_routes
}

//查找路由map找到 完全匹配对应的目的地址
func (arm *ArRouteMap) MatchRoute(url string) (*ArRoute, bool) {
	arm.mu.Lock()
	defer arm.mu.Unlock()

	r, ok := arm.CacheRoutes[url] //优先查找缓存
	if ok {                       //找到缓存按缓存处理
		return r, true
	} else { //找不到
		ar, _, find := arm.FindRouteByReqUrl(url)

		if find {
			arm.CacheRoutes[url] = ar
			return ar, true
		} else {
			ar, find = arm.MatchRouteByDir(url)
			if find {
				return ar, true
			}
			return nil, false
		}
	}

}

//查找目录匹配
func (arm *ArRouteMap) MatchRouteByDir(url string) (*ArRoute, bool) {
	ar, _, find := arm.FindRouteByReqDir(url)
	if find {
		arm.CacheRoutes[url] = ar //加入缓存
		return ar, true
	} else {
		return nil, false
	}
}

func (arm *ArRouteMap) FindRouteByReqUrl(r string) (*ArRoute, int, bool) {

	for i, arr := range arm.Routes {
		if arr.ReqUrl == r {
			return arr, i, true
		}
	}

	return nil, 0, false
}

//根据目录前缀匹配
func (arm *ArRouteMap) FindRouteByReqDir(r string) (*ArRoute, int, bool) {
	for i, arr := range arm.Routes {
		if strings.Index(r, arr.ReqUrl) == 0 {
			return arr, i, true
		}
	}
	return nil, 0, false
}

func (arm *ArRouteMap) FindRouteByID(publishID int) (*ArRoute, int, bool) {

	for i, arr := range arm.Routes {
		if arr.PublishID == publishID {
			return arr, i, true
		}
	}

	return nil, 0, false
}

func (arm *ArRouteMap) ClearRouteCache() {
	for key, _ := range arm.CacheRoutes {
		delete(arm.CacheRoutes, key)
	}
}

//从json加载路由信息，返回值代表更新的路由数量
func (arm *ArRouteMap) RoadRoute(newAroute *ArRouteLoad) (result int) {
	arm.mu.Lock()
	defer arm.mu.Unlock()

	if newAroute.MaxVer <= arm.MaxVer { //版本低于当前版本
		return result
	}

	arm.ClearRouteCache() //清理匹配缓存

	for _, newroute := range newAroute.Routes { //检查路由信息并进行替换

		oldroute, oldrouteIndex, ok := arm.FindRouteByID(newroute.PublishID)
		if ok {
			if newroute.Ver > oldroute.Ver {
				if newroute.Status != 0 {
					oldroute.CloseProxyWorks() //关闭服务

					//删除
					arm.Routes = append(arm.Routes[:oldrouteIndex], arm.Routes[oldrouteIndex+1:]...)

					//delete(arm.Routes, oldroute.ReqUrl)
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
				//				//发布地址发生了改变
				//				if oldroute.ReqUrl != newroute.ReqUrl {
				//					delete(arm.Routes, oldroute.ReqUrl) //从路由表内删除
				//					arm.Routes.PushBack(oldroute)
				//					//arm.Routes[newroute.ReqUrl] = oldroute //使用新的key添加
				//				}

				oldroute.ReqUrl = newroute.ReqUrl
				oldroute.SecretType = newroute.SecretType
				oldroute.TimeOut = newroute.TimeOut
				oldroute.Ver = newroute.Ver
				oldroute.WaitTimeOut = newroute.WaitTimeOut
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

			newroute.InitTransport()

			//arm.Routes[newroute.ReqUrl] = newroute
			arm.Routes = append(arm.Routes, newroute)

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
	//Timeout   int64
}

func NewProxyWork(writer http.ResponseWriter, req *http.Request, desturl string, t *http.Transport) *ProxyWork {
	pw := &ProxyWork{}
	pw.w = writer
	pw.DestUrl = desturl
	pw.r = req
	pw.Transport = t
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
	defer resp.Body.Close()
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			pw.rspStatus = 408 //超时
		} else {
			log.Debug("503 ", err.Error())
			pw.rspStatus = 503 //其他错误
		}
		return err
	} else {

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
	ln             *StoppableListener.StoppableListener
	ch             chan bool //关闭用
	dbLogger       *ServiceGatewayLogger
}

func NewArProxy(port string, r RouteReader, dlogger *ServiceGatewayLogger) *ArProxy {
	arp := &ArProxy{}
	arp.port = port
	arp.service_queue = make(map[int]chan byte)
	arp.service_routes = NewArRouteMap()
	arp.wa = NewWatcher(r)
	arp.ch = make(chan bool)
	arp.dbLogger = dlogger
	return arp
}

func (arp *ArProxy) checkRequest(r *http.Request, route *ArRoute) error {
	if route.SecretType == "nil" {
		return nil
	}

	//签名
	if route.SecretType == "sig" {

		h := r.Header.Get("SecHeader")
		//log.Debug("header:", h)
		if h != "" {
			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				r.Body = ioutil.NopCloser(bytes.NewReader(body))
				md := md5.New()
				_, err = io.Copy(md, bytes.NewReader([]byte(string(body)+route.SecKey)))

				cMd5 := fmt.Sprintf("%x", md.Sum(nil))

				if cMd5 == h {
					return nil
				} else {
					log.Error("sig error:", r.URL.Path, " reqbody:", string(body), " reqsig:", h, " checkmd5:", cMd5)
					return errors.New("sig error:" + h)
				}
			} else {
				log.Error("Read Body error:", err.Error())
				return err
			}

		} else {
			log.Error("sig SecHeader error:" + h)
			return errors.New("sig SecHeader error:" + h)
		}

	}

	//3DES
	if route.SecretType == "encrypt" {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {

			od, err := TripleDesDecrypt(body, []byte(route.SecKey))
			if err == nil {
				r.Body = ioutil.NopCloser(bytes.NewReader(od))
				return nil
			} else {
				log.Error("Decrypt error: ", r.URL.Path, " reqbody:", string(body), ";", err.Error())
				return errors.New("Decrypt error")
			}

		} else {
			log.Error("Read Body error:", err.Error())
			return err
		}
	}

	return nil
}

func (arp *ArProxy) handleService(w http.ResponseWriter, r *http.Request) {

	log.Debug("receive request ", r.URL.Path)
	rurl := r.URL.Path
	rhost := r.URL.Host
	begintime := time.Now()

	ar, ok := arp.service_routes.MatchRoute(r.URL.Path)
	if ok {
		log.Debug("match request ", r.URL.Path, " ", ar.ReqUrl)

		err := arp.checkRequest(r, ar) //安全检查
		if err != nil {
			w.Write([]byte("Secret Error " + err.Error()))
			log.Info("route to err:", ar.Name, "$$", rurl, "$$", "", "$$", "400")
			arp.dbLogger.AddLog(ar, rurl, "", begintime, time.Now(), 404, err.Error(), rhost)
			return
		}

		var desturl string
		if ar.MatchMode == "MatchDir" { //发现是目录匹配,进行路径处理

			desturl = strings.Replace(rurl, ar.ReqUrl, ar.ProxyToUrl, 1)
		} else {
			desturl = ar.ProxyToUrl
		}

		pw := NewProxyWork(w, r, desturl, ar.transport)
		//pw.Timeout = ar.WaitTimeOut
		err = ar.ProxyWorks.PutWork(pw, time.Duration(ar.TimeOut)*time.Millisecond)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "work Error: %v", err)
			log.Info("route to err:", ar.Name, "$$", rurl, "$$", desturl, "$$", strconv.Itoa(http.StatusServiceUnavailable))
			arp.dbLogger.AddLog(ar, rurl, desturl, begintime, time.Now(), http.StatusServiceUnavailable, err.Error(), rhost)
		} else {
			arp.dbLogger.AddLog(ar, rurl, desturl, begintime, time.Now(), pw.rspStatus, "", rhost)
		}
	} else {
		http.ServeFile(w, r, "NotFindService.html")
		log.Info("route Match:", "", "$$", r.URL.Path, "$$", "", "$$", "404")
		arp.dbLogger.AddLog(nil, rurl, "NotFindService.html", begintime, time.Now(), 404, "", rhost)
	}

	//w.Write([]byte("Hello"))

}

func (arp *ArProxy) handleReload(w http.ResponseWriter, r *http.Request) {
	Arrl, err := arp.wa.reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	rc := arp.service_routes.RoadRoute(Arrl)
	log.Info("reload routeInfo done. update routes count:", strconv.Itoa(rc))
	w.Write([]byte("Reload ok"))
}

func (arp *ArProxy) handleLog(w http.ResponseWriter, r *http.Request) {
	arp.dbLogger.FlushLog()
	log.Info("Flush log done.")
	w.Write([]byte("Flush log ok"))
}

func (arp *ArProxy) Stop() {
	arp.wa.StopWatch()
	arp.ln.Stop()
	<-arp.ch //等待所有完成
}

func (arp *ArProxy) Start() {
	Arrl, err := arp.wa.reader.Read()

	if err != nil {
		log.Fatal(err)
	}

	//	jj, err := json.Marshal(Arrl)
	//log.Debug("dd", string(jj))

	rc := arp.service_routes.RoadRoute(Arrl)

	log.Info("init routeInfo done. routes count:", strconv.Itoa(rc))

	go arp.wa.StartWatch(arp)

	originalListener, err := net.Listen("tcp", ":"+arp.port)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("start listen tcp ", arp.port)

	sl, err := StoppableListener.New(originalListener)
	if err != nil {
		log.Fatal(err)
	}

	arp.ln = sl
	http.HandleFunc("/", arp.handleService)                  //设置访问的路由
	http.HandleFunc("/Config/ReLoad.html", arp.handleReload) //重新加载
	http.HandleFunc("/Config/Log.html", arp.handleLog)       //刷新日志
	server := &http.Server{}

	log.Debug("start http Serve ", arp.port)

	go func() {
		server.Serve(sl)
		arp.ch <- true
	}()
	//	//err = http.ListenAndServe(":"+arp.port, nil) //设置监听的端口
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Info("start listen port:", arp.port)

}
