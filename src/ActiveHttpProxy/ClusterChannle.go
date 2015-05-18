package ActiveHttpProxy

import ()

type ClusterChannle struct {
	ch        chan bool
	exit      chan bool
	OtherHost map[string]string
	Host      string
}
