package ActiveHttpProxy

import (
	//	"strconv"
	//	"time"
	//	"fmt"
	//	"math/rand"
	//	"runtime"
	//	"strconv"
	"testing"
	//	"time"
	//_ "code.google.com/p/odbc"
	//	"database/sql"
)

func Test_WatcherGetRoute(t *testing.T) {
	//	w := &Watcher{}
	//	js, err := w.ReadDBRoute()
	//	if err != nil {
	//		t.Error(err)
	//	} else {
	//		t.Log(len(js.Routes))
	//	}
}

type teststr struct {
	name  string
	value int
}

func Test_ArProxyReload(t *testing.T) {
	//	//	w := NewWatcher()
	//	ap := NewArProxy("12001")
	//	json1 := `{"MaxVer":13,"Routes":[{"PublishID":1,"Name":"service1","ReqUrl":"//testservice//1.html","ProxyToUrl":"//destservice//1.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":1,"ProxyWorks":null},{"PublishID":2,"Name":"service2","ReqUrl":"//testservice//2.html","ProxyToUrl":"//destservice//2.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":2,"ProxyWorks":null},{"PublishID":3,"Name":"service3","ReqUrl":"//testservice//3.html","ProxyToUrl":"//destservice//3.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":3,"ProxyWorks":null},{"PublishID":4,"Name":"service4","ReqUrl":"//testservice//4.html","ProxyToUrl":"//destservice//4.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":4,"ProxyWorks":null},{"PublishID":5,"Name":"service5","ReqUrl":"//testservice//5.html","ProxyToUrl":"//destservice//5.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":5,"ProxyWorks":null},{"PublishID":6,"Name":"service6","ReqUrl":"//testservice//6.html","ProxyToUrl":"//destservice//6.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":6,"ProxyWorks":null},{"PublishID":7,"Name":"service7","ReqUrl":"//testservice//7.html","ProxyToUrl":"//destservice//7.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":7,"ProxyWorks":null},{"PublishID":8,"Name":"service8","ReqUrl":"//testservice//8.html","ProxyToUrl":"//destservice//8.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":8,"ProxyWorks":null},{"PublishID":9,"Name":"service9","ReqUrl":"//testservice//9.html","ProxyToUrl":"//destservice//9.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":9,"ProxyWorks":null},{"PublishID":10,"Name":"service10","ReqUrl":"//testservice//10.html","ProxyToUrl":"//destservice//10.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":10,"ProxyWorks":null},{"PublishID":11,"Name":"service11","ReqUrl":"//testservice//11.html","ProxyToUrl":"//destservice//11.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":11,"ProxyWorks":null},{"PublishID":12,"Name":"service12","ReqUrl":"//testservice//12.html","ProxyToUrl":"//destservice//12.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":12,"ProxyWorks":null},{"PublishID":13,"Name":"service13","ReqUrl":"//testservice//13.html","ProxyToUrl":"//destservice//13.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":13,"ProxyWorks":null}]}`
	//	json2 := `{"MaxVer":13,"Routes":[{"PublishID":1,"Name":"service1","ReqUrl":"//testservice//1.html","ProxyToUrl":"//destservice//16.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":16,"ProxyWorks":null},{"PublishID":2,"Name":"service2","ReqUrl":"//testservice//2.html","ProxyToUrl":"//destservice//2.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":2,"ProxyWorks":null},{"PublishID":3,"Name":"service3","ReqUrl":"//testservice//3.html","ProxyToUrl":"//destservice//3.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":3,"ProxyWorks":null},{"PublishID":4,"Name":"service4","ReqUrl":"//testservice//4.html","ProxyToUrl":"//destservice//4.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":4,"ProxyWorks":null},{"PublishID":5,"Name":"service5","ReqUrl":"//testservice//5.html","ProxyToUrl":"//destservice//5.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":5,"ProxyWorks":null},{"PublishID":6,"Name":"service6","ReqUrl":"//testservice//6.html","ProxyToUrl":"//destservice//6.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":6,"ProxyWorks":null},{"PublishID":7,"Name":"service7","ReqUrl":"//testservice//7.html","ProxyToUrl":"//destservice//7.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":7,"ProxyWorks":null},{"PublishID":8,"Name":"service8","ReqUrl":"//testservice//8.html","ProxyToUrl":"//destservice//8.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":8,"ProxyWorks":null},{"PublishID":9,"Name":"service9","ReqUrl":"//testservice//9.html","ProxyToUrl":"//destservice//9.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":9,"ProxyWorks":null},{"PublishID":10,"Name":"service10","ReqUrl":"//testservice//10.html","ProxyToUrl":"//destservice//10.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":10,"ProxyWorks":null},{"PublishID":11,"Name":"service11","ReqUrl":"//testservice//11.html","ProxyToUrl":"//destservice//11.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":11,"ProxyWorks":null},{"PublishID":12,"Name":"service12","ReqUrl":"//testservice//12.html","ProxyToUrl":"//destservice//12.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":12,"ProxyWorks":null},{"PublishID":13,"Name":"service13","ReqUrl":"//testservice//13.html","ProxyToUrl":"//destservice//18.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000000000,"Ver":18,"ProxyWorks":null}]}`
	//	ap.service_routes.RoadRoute(json1)

	//	t1, err1 := ap.service_routes.Routes["//testservice//1.html"]

	//	if !err1 {
	//		t.Error("not find ", t1)
	//	} else {
	//		if t1.ProxyToUrl != "//destservice//1.html" {
	//			t.Error("ProxyToUrl value error ", t1)
	//		}
	//	}

	//	ap.service_routes.RoadRoute(json2)

	//	t2, err2 := ap.service_routes.Routes["//testservice//1.html"]

	//	if !err2 {
	//		t.Error("not find ", t2)
	//	} else {
	//		if t2.ProxyToUrl != "//destservice//16.html" {
	//			t.Error("ProxyToUrl value error ", t2)
	//		}
	//	}
	//ap.service_routes.RoadRoute(json2)
}

//func Test_WatcherTestWatcher(t *testing.T) {
//	w := NewWatcher()
//	ap := NewArProxy("12001")
//	go w.StartWatch(ap)
//	time.Sleep(30 * time.Second)
//	w.StopWatch()
//}
