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
	Deployments []Deployments `json:"data" description:"kubernetes的Deployments部署列表" type:"[]Deployments"`
	PageSize    int           `json:"pageSize" description:"一页大小" type:"int"`
	PageNo      int           `json:"pageNo" description:"当前所处的页数" type:"int"`
	TotalPage   int           `json:"totalPage" description:"总页数" type:"int"`
	TotalCount  int           `json:"totalCount" description:"总数量" type:"int"`
}

//Deployments Deployments部署
type Deployments struct {
	Name              string `json:"name" description:"名称" type:"string"`
	NameSpace         string `json:"nameSpace" description:"命名空间" type:"string"`
	CreateTime        string `json:"createTime" description:"创建时间" type:"string"`
	Replicas          int32  `json:"replicas" description:"部署的副本数量" type:"string"`
	AvailableReplicas int32  `json:"availableReplicas" description:"部署可用的副本数量" type:"string"`
}

//GetKubernetesDeploymentInfoByNameReq 根据Name获取kubernetes的Deployments部署的详情
type GetKubernetesDeploymentInfoByNameReq struct {
	Name      string `json:"name" description:"部署的名称" type:"string"`
	NameSpace string `json:"nameSpace" description:"命名空间" type:"string"`
}

//GetKubernetesDeploymentInfoByNameRes 根据Name获取kubernetes的Deployments部署的详情
type GetKubernetesDeploymentInfoByNameRes struct {
	Name              string                 `json:"name" description:"名称" type:"string"`
	NameSpace         string                 `json:"nameSpace" description:"命名空间" type:"string"`
	CreateTime        string                 `json:"createTime" description:"创建时间" type:"string"`
	Replicas          int32                  `json:"replicas" description:"部署的副本数量" type:"string"`
	AvailableReplicas int32                  `json:"availableReplicas" description:"部署可用的副本数量" type:"string"`
	Labels            map[string]string      `json:"labels" description:"标签" type:"object"`
	Annotations       map[string]string      `json:"annotations" description:"注解" type:"object"`
	Conditions        []DeploymentConditions `json:"conditions" description:"Deployment运行时候信息" type:"[]DeploymentConditions"`
	ReplicaSets       []ReplicaSet           `json:"replicaSets" description:"副本集合" type:"[]ReplicaSet"`
	Pods              []Pod                  `json:"pods" description:"pod集合" type:"[]Pod"`
}

//DeploymentConditions  Deploymen运行状态
type DeploymentConditions struct {
	Type           string `json:"type" description:"类型" type:"string"`
	Status         string `json:"status" description:"状态" type:"string"`
	LastUpdateTime string `json:"lastUpdateTime" description:"最后更新时间" type:"string"`
	Reason         string `json:"reason" description:"原因" type:"string"`
	Message        string `json:"message" description:"消息" type:"string"`
}

//ReplicaSet 副本集
type ReplicaSet struct {
	Name       string `json:"name" description:"名称" type:"string"`
	Desired    int32  `json:"desired" description:"期望" type:"int32"`
	Current    int32  `json:"current" description:"当前的值" type:"int32"`
	Ready      int32  `json:"ready" description:"已经启动的值" type:"int32"`
	CreateTime string `json:"createTime" description:"创建时间" type:"string"`
}

//Pod  Pod信息
type Pod struct {
	Name          string         `json:"name" description:"名称" type:"string"`
	NameSpace     string         `json:"nameSpace" description:"命名空间" type:"string"`
	Phase         string         `json:"phase" description:"阶段" type:"string"`
	CreateTime    string         `json:"createTime" description:"创建时间" type:"string"`
	StartTime     string         `json:"startTime" description:"启动时间" type:"string"`
	PodIP         string         `json:"podIP" description:"PodIP" type:"string"`
	HostIP        string         `json:"hostIP" description:"NodeIP" type:"string"`
	Nodes         string         `json:"nodes" description:"节点名称" type:"string"`
	RestartPolicy string         `json:"restartPolicy" description:"镜像重启方式" type:"string"`
	Events        []PodEvent     `json:"events" description:"事件" type:"[]PodEvent"`
	Conditions    []PodCondition `json:"conditions" description:"pod的运行状况" type:"[]PodCondition"`
	Containers    []PodContainer `json:"containers" description:"容器" type:"[]PodContainers"`
}

// PodCondition Pod状态
type PodCondition struct {
	Type               string `json:"type" description:"类型" type:"string"`
	Status             string `json:"status" description:"状态" type:"string"`
	LastTransitionTime string `json:"lastTransitionTime" description:"上次转换时间" type:"string"`
	Reason             string `json:"reason" description:"原因" type:"string"`
	Message            string `json:"message" description:"消息" type:"string"`
}

//Event  事件
type PodEvent struct {
	Type           string `json:"type" description:"事件类型" type:"string"`
	Source         string `json:"source" description:"来源" type:"string"`
	Message        string `json:"message" description:"消息" type:"string"`
	Reason         string `json:"reason" description:"原因" type:"string"`
	FirstTimestamp string `json:"firstTimestamp" description:"第一次时间" type:"string"`
	LastTimestamp  string `json:"lastTimestamp" description:"最后一次时间" type:"string"`
}

//PodContainers Pod里面容器定义
type PodContainer struct {
	Name            string          `json:"name" description:"名称" type:"string"`
	Image           string          `json:"image" description:"镜像" type:"string"`
	ContainerSpec   ContainerSpec   `json:"containerSpec" description:"定义容器的规格" type:"ContainerSpec"`
	ContainerStatus ContainerStatus `json:"containerStatus" description:"容器的运行状态" type:"ContainerStatus"`
}

//ContainerSpec  容器规格定义
type ContainerSpec struct {
	ImagePullPolicy string   `json:"imagePullPolicy" description:"镜像拉取方式" type:"string"`
	Command         []string `json:"command" description:"容器的运行状态" type:"[]string"`
	Args            []string `json:"args" description:"容器的运行状态" type:"[]string"`
}

// ContainerStatus 容器状态
type ContainerStatus struct {
	Ready        bool   `json:"ready" description:"是否准备成功" type:"bool"`
	StartedAt    string `json:"startedAt" description:"启动时间" type:"string"`
	RestartCount int32  `json:"restartCount" description:"重启次数" type:"int32"`
	Image        string `json:"image" description:"镜像" type:"string"`
	ImageID      string `json:"imageID" description:"镜像ID" type:"string"`
	ContainerID  string `json:"containerID" description:"docker的container id" type:"string"`
	Started      bool   `json:"started" description:"是否启动" type:"bool"`
}
