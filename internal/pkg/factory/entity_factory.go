package factory

import (
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
	plugin "gitlab.intsig.net/textin-gateway/pkg/apisix/plugin"
	"gitlab.intsig.net/textin-gateway/pkg/md5"
)

func GeneUpstream(name string, description string, nodes map[string]int) *apisix.Upstream {
	var data apisix.Upstream
	var idleTimeout = float32(60)
	data.UpstreamDef = apisix.UpstreamDef{
		Name:   name,
		Desc:   description,
		Type:   "roundrobin",
		Scheme: "http",
		Nodes:  nodes,
		Timeout: &apisix.Timeout{
			Connect: 6,
			Read:    6,
			Send:    6,
		},
		PassHost: "pass",
		KeepalivePool: &apisix.UpstreamKeepalivePool{
			IdleTimeout: (*apisix.TimeoutValue)(&idleTimeout),
			Requests:    1000,
			Size:        320,
		},
	}

	data.BaseInfo.Creating()
	return &data
}

func GeneRoute(name string, description string, nodes map[string]int) *apisix.Upstream {
	var data apisix.Upstream
	var idleTimeout = float32(60)
	data.UpstreamDef = apisix.UpstreamDef{
		Name:   name,
		Desc:   description,
		Type:   "roundrobin",
		Scheme: "http",
		Nodes:  nodes,
		Timeout: &apisix.Timeout{
			Connect: 6,
			Read:    6,
			Send:    6,
		},
		PassHost: "pass",
		KeepalivePool: &apisix.UpstreamKeepalivePool{
			IdleTimeout: (*apisix.TimeoutValue)(&idleTimeout),
			Requests:    1000,
			Size:        320,
		},
	}

	data.BaseInfo.Creating()
	return &data
}

func GeneTestRoute(name, upstream_id string) *apisix.Route {
	var (
		data       apisix.Route
		vars       []interface{}
		plugins    = make(map[string]interface{})
		api_key    = md5.GetMD5Encode(apisix.GetFlakeUidStr())
		api_secret = md5.GetMD5Encode(apisix.GetFlakeUidStr())
	)
	vars = append(vars, []interface{}{"http_App_key", "==", api_key})
	vars = append(vars, []interface{}{"http_App_secret", "==", api_secret})

	plugins["cors"] = plugin.Cors{
		AllowCredential: false,
		AllowHeaders:    "*",
		AllowOrigins:    "*",
		AllowMethods:    "*",
		MaxAge:          5,
		Disable:         false,
	}
	plugins["request-id"] = plugin.RequestId{
		Disable: false,
	}

	data = apisix.Route{
		Name: name,
		Methods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
		},
		Plugins:    plugins,
		Status:     1,
		UpstreamID: upstream_id,
		Vars:       vars,
		URI:        "/*",
	}
	data.BaseInfo.Creating()
	return &data
}
