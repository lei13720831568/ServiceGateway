package ActiveHttpProxy //Reverse Porxy

import (
	log "RollLoger"
	"WorkPool"
	"encoding/json"
	//	"fmt"
	"net/http"
	//	"runtime"
	"strconv"
	"sync"
	//"time"
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
	ProxyWorks  *WorkPool.WPool
}

func (ar *ArRoute) ToJson() string {
	r, err := json.Marshal(ar)
	if err != nil {
		log.Error("Marshal ArRoute To json failed;")
	}
	return string(r)
}

func (ar *ArRoute) InitProxyWorks() {
	ar.ProxyWorks = WorkPool.NewWorkPool(ar.Name)
	ar.ProxyWorks.SetMax(ar.MaxConnects)
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

//从json加载路由信息，返回值代表更新的路由数量
func (arm *ArRouteMap) RoadRoute(jsonstr string) (result int) {

	result = 0
	var newAroute ArRouteLoad
	err := json.Unmarshal([]byte(jsonstr), &newAroute) //反序列化json
	if err != nil {
		log.Error("error Unmarshal route json:", err.Error())
	}

	arm.mu.Lock()
	defer arm.mu.Unlock()

	if newAroute.MaxVer <= arm.MaxVer { //版本低于当前版本
		return result
	}

	for _, newroute := range newAroute.Routes { //检查路由信息并进行替换

		oldroute, ok := arm.Routes[newroute.ReqUrl]
		if ok {
			if newroute.Ver > oldroute.Ver {
				//更新属性
				oldroute.Encrypt = newroute.Encrypt
				oldroute.IpList = newroute.IpList
				oldroute.Name = newroute.Name
				oldroute.ProxyToUrl = newroute.ProxyToUrl
				oldroute.PublishID = newroute.PublishID
				oldroute.ReqUrl = newroute.ReqUrl
				oldroute.SecretType = newroute.SecretType
				oldroute.TimeOut = newroute.TimeOut
				oldroute.Ver = newroute.Ver
				//调整连接数
				err = oldroute.ProxyWorks.SetMax(newroute.MaxConnects)
				if err == nil {
					oldroute.MaxConnects = newroute.MaxConnects
				} else {
					log.Error("连接调整失败old:", strconv.Itoa(oldroute.MaxConnects), " new:", strconv.Itoa(newroute.MaxConnects))
				}
				result++
				log.Info("update route to ", oldroute.ToJson())
			}

		} else {
			newroute.ProxyWorks = WorkPool.NewWorkPool(newroute.Name)
			err = newroute.ProxyWorks.SetMax(newroute.MaxConnects)
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

type ArProxy struct {
	port           string
	service_queue  map[int]chan byte
	service_routes *ArRouteMap
}

func NewArProxy(port string) *ArProxy {
	arp := &ArProxy{}
	arp.port = port
	arp.service_queue = make(map[int]chan byte)
	arp.service_routes = NewArRouteMap()
	return arp
}

func (arp *ArProxy) HandleService(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello"))

}

func (arp *ArProxy) Start() {
	http.HandleFunc("/", arp.HandleService)       //设置访问的路由
	err := http.ListenAndServe(":"+arp.port, nil) //设置监听的端口
	if err != nil {
		log.Fatal(err)
	}

	log.Info("start listen port:", arp.port)
}
