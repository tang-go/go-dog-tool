package param

//GetKubernetesNameSpaceReq 获取kubernetes的namespace请求
type GetKubernetesNameSpaceReq struct {
}

//GetKubernetesNameSpaceReq 获取kubernetes的namespace请求
type GetKubernetesNameSpaceRes struct {
	K8sNameSpace []K8sNameSpace `json:"data" description:"kubernetes的namespace" type:"[]K8sNameSpace"`
	PageSize     int            `json:"pageSize" description:"一页大小" type:"int"`
	PageNo       int            `json:"pageNo" description:"当前所处的页数" type:"int"`
	TotalPage    int            `json:"totalPage" description:"总页数" type:"int"`
	TotalCount   int            `json:"totalCount" description:"总数量" type:"int"`
}

//K8sNameSpace kubernetes的namespace
type K8sNameSpace struct {
	Name       string `json:"name" description:"名称" type:"string"`
	CreateTime string `json:"createTime" description:"创建时间" type:"string"`
	Satus      string `json:"status" description:"状态" type:"string"`
}

//GetKubernetesDeploymentsReq 获取kubernetes的Deployments部署
type GetKubernetesDeploymentsReq struct {
	NameSpace string `json:"nameSpace" description:"命名空间" type:"string"`
}

//GetKubernetesDeploymentsRes 获取kubernetes的Deployments部署
type GetKubernetesDeploymentsRes struct {
	Pods       []Pod `json:"data" description:"kubernetes的namespace" type:"[]Pods"`
	PageSize   int   `json:"pageSize" description:"一页大小" type:"int"`
	PageNo     int   `json:"pageNo" description:"当前所处的页数" type:"int"`
	TotalPage  int   `json:"totalPage" description:"总页数" type:"int"`
	TotalCount int   `json:"totalCount" description:"总数量" type:"int"`
}

//Pods Pods信息
type Pod struct {
	Name                string `json:"name" description:"名称" type:"string"`
	NameSpace           string `json:"nameSpace" description:"命名空间" type:"string"`
	CreateTime          string `json:"createTime" description:"创建时间" type:"string"`
	ObservedGeneration  int64  `json:"observedGeneration" description:"部署控制器观察到的生成的数量" type:"string"`
	Replicas            int32  `json:"replicas" description:"部署的副本数量" type:"string"`
	UpdatedReplicas     int32  `json:"updatedReplicas" description:"创建时间" type:"string"`
	ReadyReplicas       int32  `json:"readyReplicas" description:"部署Ready的副本数量" type:"string"`
	AvailableReplicas   int32  `json:"availableReplicas" description:"部署可用的副本数量" type:"string"`
	UnavailableReplicas int32  `json:"unavailableReplicas" description:"部署不可用副本的数量" type:"string"`
}
