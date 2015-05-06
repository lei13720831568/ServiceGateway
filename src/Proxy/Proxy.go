package ActiveHttpReverseProxy

import (
	log "RollLoger"
	"WorkPool"
	"encoding/json"
	"net/http"
	//"strconv"
	"runtime"
	"sync"
	"time"
)

type ArRouteLoad struct {
	MaxVer int64
	Routes []ArRoute
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

//func SetRoute(r *ArRoute) *ArRoute {
//	arr := new(ArRoute)
//	arr.Encrypt = r.Encrypt
//	arr.IpList = r.IpList
//	arr.MaxConnects = r.MaxConnects
//	arr.ProxyToUrl = r.ProxyToUrl
//	arr.PublishID = r.PublishID
//	arr.ReqUrl = r.ReqUrl
//	arr.SecretType = r.SecretType
//	arr.TimeOut = r.TimeOut
//	arr.Ver = r.Ver
//	return arr
//}

func (ar *ArRoute) InitProxyWorks() {
	ar.ProxyWorks = WorkPool.NewWorkPool(ar.Name, ar.MaxConnects, ar.TimeOut*time.Millisecond)
	ar.ProxyWorks.Start()
}

func (ar *ArRoute) CloseProxyWorks() {
	ar.ProxyWorks.Stop()
}

type ArRouteMap struct {
	mu     sync.Mutex
	Routes map[string]*ArRoute
	MaxVer int64
}

func (arm *ArRouteMap) MatchRoute(url string) (*ArRoute, bool) {
	arm.mu.Lock()
	defer arm.mu.Unlock()
	r, ok := arm.Routes[url]
	if ok {
		return copyRoute(r), true
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
				go oldroute.CloseProxyWorks() //关闭旧有的线程池
				oldroute = &newroute
				oldroute.InitProxyWorks()
				oldroute.ProxyWorks.Start() //启动工作线程池
				result++
			}
		} else {
			arm.Routes[newroute.ReqUrl] = &newroute
			result++
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
	return arp
}

func (arp *ArProxy) HandleService(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello"))

}

func (arp *ArProxy) Start() {
	arp.service_queue = make(map[int]chan byte)
	arp.service_routes = &ArRouteMap{}
	arp.service_routes.Routes = make(map[string]*ArRoute)

	http.HandleFunc("/", arp.HandleService)       //设置访问的路由
	err := http.ListenAndServe(":"+arp.port, nil) //设置监听的端口
	if err != nil {
		log.Fatal(err)
	}

	log.Info("start listen ")
}
