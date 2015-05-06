package ActiveHttpReverseProxy

import (
	//"strings"
	"testing"
)

func Test_getAdmin(t *testing.T) {

	arm := &ArRouteMap{}
	arm.Routes = make(map[string]*ArRoute)
	b := arm.RoadRoute("sfasf")
	if !b {
		t.Error("error")
	}

}
