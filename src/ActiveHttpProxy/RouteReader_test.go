package ActiveHttpProxy

import (
	"testing"
)

func Test_ReadRouteFromDB(t *testing.T) {
	r := &ReaderFromDB{"driver={sql server};server=192.168.1.100;port=1433;uid=sa;pwd=654321;database=ADC3CoreDB"}
	arrl, err := r.Read()
	if err != nil {
		t.Error(err)
	} else if arrl == nil {
		t.Error(err)
	} else {
		t.Log(arrl)
	}
}

func Test_ReadRouteFormJson(t *testing.T) {
	//str := `{"MaxVer":0,"Routes":[{"PublishID":1,"Name":"服务1","ReqUrl":"/s1/test1.html","ProxyToUrl":"http://192.168.1.129/test/test1.html","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":10,"TimeOut":10000,"Ver":1879541,"Status":0,"ProxyWorks":null},{"PublishID":2,"Name":"服务2","ReqUrl":"/s2/test2.html","ProxyToUrl":"http://192.168.1.100:8088/ForOMP/ADCInterfaceForOMP.asmx","IpList":"nil","SecretType":"nil","Encrypt":"nil","MaxConnects":300,"TimeOut":10000,"Ver":1879554,"Status":0,"ProxyWorks":null}]}`
	r := &ReaderFromJsonFile{"testRoutes.json"}
	arrl, err := r.Read()

	if err != nil {
		t.Error(err)
	} else if arrl == nil {
		t.Error(err)
	} else {
		t.Log(arrl)
	}
}
