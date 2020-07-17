package client

import (
	"time"
)

const (
	networksPath               = "/api/networks"
	hostsPath                  = "/api/hosts"
	interfacesPath             = "/api/interfaces"
	linksPath                  = "/api/links"
	tokensPath                 = "/api/tokens"
	labelQueryParamKey         = "label"
	nameQueryParamKey          = "name"
	idQueryParamKey            = "id"
	hostIDQueryParamKey        = "hostId"
	networkIDQueryParamKey     = "networkId"
	pageQueryParamKey          = "page"
	perPageQueryParamKey       = "perPage"
	defaultTokenHeader         = "X-Drago-Token"
	contentTypeHeader          = "Content-Type"
	contentTypeApplicationJson = "application/json"
)

type DragoAPIClient struct {
	Timeout   time.Duration
	Token     string
	ServerURL string
}
