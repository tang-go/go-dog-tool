package api

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog/plugins"
)

//GetServiceList 获取服务列表
func (pointer *API) GetServiceList(ctx plugins.Context, req param.GetServiceReq) (res param.GetServiceRes, err error) {
	services := pointer.service.GetClient().GetAllRPCService()
	mp := make(map[string]*param.Services)
	for _, service := range services {
		if service.Name == define.SvcController {
			continue
		}
		info := &param.ServiceInfo{
			Key:       service.Key,
			Name:      service.Name,
			Address:   service.Address,
			Port:      service.Port,
			Explain:   service.Explain,
			Longitude: service.Longitude,
			Latitude:  service.Latitude,
			Time:      service.Time,
		}
		for _, method := range service.Methods {
			info.Methods = append(info.Methods, &param.Method{
				Name:     method.Name,
				Level:    method.Level,
				Request:  method.Request,
				Response: method.Response,
				Explain:  method.Explain,
				IsAuth:   method.IsAuth,
			})
		}
		s, ok := mp[service.Name]
		if !ok {
			mp[service.Name] = &param.Services{
				Name:    service.Name,
				Explain: service.Explain,
				Info: []*param.ServiceInfo{
					info,
				},
			}
		} else {
			s.Info = append(s.Info, info)
		}
	}
	for _, s := range mp {
		res.List = append(res.List, s)
	}
	return
}
