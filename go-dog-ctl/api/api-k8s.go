package api

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//GetKubernetesNameSpace 获取kubernetes的namespqce
func (pointer *API) GetKubernetesNameSpace(ctx plugins.Context, request param.GetKubernetesNameSpaceReq) (response param.GetKubernetesNameSpaceRes, err error) {
	namespaces, e := pointer.clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, "获取kubernetes的namespqce失败")
		return
	}
	total := len(namespaces.Items)
	response.PageNo = 1
	response.PageSize = total
	response.TotalCount = total
	if total <= 0 {
		response.TotalPage = 0
		return
	}
	if total%response.PageSize > 0 {
		response.TotalPage = total/response.PageSize + 1
	}
	if total%response.PageSize < 0 {
		response.TotalPage = total / response.PageSize
	}
	for _, ns := range namespaces.Items {
		response.K8sNameSpace = append(response.K8sNameSpace, param.K8sNameSpace{
			Name:       ns.Name,
			Satus:      string(ns.Status.Phase),
			CreateTime: ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return
}

//GetKubernetesDeployments 获取kubernetes的Deployments部署
func (pointer *API) GetKubernetesDeployments(ctx plugins.Context, request param.GetKubernetesDeploymentsReq) (response param.GetKubernetesDeploymentsRes, err error) {
	deployments, e := pointer.clientSet.AppsV1().Deployments(request.NameSpace).List(metav1.ListOptions{})
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetKubernetesPodsByNamsespace, "获取kubernetes的Deployments部署失败")
		return
	}
	total := len(deployments.Items)
	response.PageNo = 1
	response.PageSize = total
	response.TotalCount = total
	if total <= 0 {
		response.TotalPage = 0
		return
	}
	if total%response.PageSize > 0 {
		response.TotalPage = total/response.PageSize + 1
	}
	if total%response.PageSize < 0 {
		response.TotalPage = total / response.PageSize
	}
	for _, deployment := range deployments.Items {
		response.Deployments = append(response.Deployments, param.Deployments{
			Name:              deployment.Name,
			NameSpace:         deployment.Namespace,
			Replicas:          deployment.Status.Replicas,
			AvailableReplicas: deployment.Status.AvailableReplicas,
			CreateTime:        deployment.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return
}

//GetKubernetesDeploymentInfoByName 根据Name获取kubernetes的Deployments部署的详情
func (pointer *API) GetKubernetesDeploymentInfoByName(ctx plugins.Context, request param.GetKubernetesDeploymentInfoByNameReq) (response param.GetKubernetesDeploymentInfoByNameRes, err error) {
	deployment, e := pointer.clientSet.AppsV1().Deployments(request.NameSpace).Get(request.Name, metav1.GetOptions{})
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetKubernetesDeploymentInfoByNameErr, "Deployments部署的详情失败")
		return
	}
	response.Name = deployment.Name
	response.NameSpace = deployment.Namespace
	response.Labels = deployment.Labels
	response.Replicas = deployment.Status.Replicas
	response.AvailableReplicas = deployment.Status.AvailableReplicas
	response.Annotations = deployment.Annotations
	response.CreateTime = deployment.CreationTimestamp.Format("2006-01-02 15:04:05")
	//生成label选择器
	var labelSelector []string
	for key, value := range deployment.Labels {
		labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", key, value))
	}
	listOptions := metav1.ListOptions{
		LabelSelector: strings.Join(labelSelector, ","),
	}
	//获取service
	if service, e := pointer.clientSet.CoreV1().Services(request.NameSpace).Get(request.Name, metav1.GetOptions{}); e == nil {
		response.Service.Labels = service.Labels
		response.Service.Annotations = service.Annotations
		response.Service.CreateTime = service.CreationTimestamp.Format("2006-01-02 15:04:05")
		response.Service.ClusterIP = service.Spec.ClusterIP
		response.Service.Selector = service.Spec.Selector
		response.Service.Type = string(service.Spec.Type)
		for _, port := range service.Spec.Ports {
			response.Service.Ports = append(response.Service.Ports, param.ServicePort{
				Name:       port.Name,
				NodePort:   port.NodePort,
				Port:       port.Port,
				TargetPort: port.TargetPort.IntVal,
				Protocol:   string(port.Protocol),
			})
		}
	}
	//获取副本集合
	replicaSets, e := pointer.clientSet.AppsV1().ReplicaSets(request.NameSpace).List(listOptions)
	if e == nil {
		for _, replicaSet := range replicaSets.Items {
			response.ReplicaSets = append(response.ReplicaSets, param.ReplicaSet{
				Name:       replicaSet.Name,
				Desired:    replicaSet.Status.Replicas,
				Current:    replicaSet.Status.AvailableReplicas,
				Ready:      replicaSet.Status.ReadyReplicas,
				CreateTime: replicaSet.CreationTimestamp.Format("2006-01-02 15:04:05"),
			})
		}
	}
	//获取pod
	pods, e := pointer.clientSet.CoreV1().Pods(request.NameSpace).List(listOptions)
	if e == nil {
		for _, pod := range pods.Items {
			//pod基本信息
			p := param.Pod{
				Name:          pod.Name,
				NameSpace:     pod.Namespace,
				CreateTime:    pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
				StartTime:     pod.Status.StartTime.Format("2006-01-02 15:04:05"),
				Phase:         string(pod.Status.Phase),
				PodIP:         pod.Status.PodIP,
				HostIP:        pod.Status.HostIP,
				Nodes:         pod.Spec.NodeName,
				RestartPolicy: string(pod.Spec.RestartPolicy),
			}
			//pod的事件
			events, e := pointer.clientSet.CoreV1().Events(request.NameSpace).List(metav1.ListOptions{
				FieldSelector: "involvedObject.name=" + pod.Name,
			})
			if e == nil {
				for _, event := range events.Items {
					p.Events = append(p.Events, param.PodEvent{
						Type:           event.Type,
						Reason:         event.Reason,
						Source:         event.Source.Component + "," + event.Source.Host,
						Message:        event.Message,
						FirstTimestamp: event.FirstTimestamp.Format("2006-01-02 15:04:05"),
						LastTimestamp:  event.LastTimestamp.Format("2006-01-02 15:04:05"),
					})
				}

			}
			//pod的运行状况
			for _, condition := range pod.Status.Conditions {
				p.Conditions = append(p.Conditions, param.PodCondition{
					Type:               string(condition.Type),
					Status:             string(condition.Status),
					LastTransitionTime: condition.LastTransitionTime.Format("2006-01-02 15:04:05"),
					Reason:             condition.Reason,
					Message:            condition.Message,
				})
			}
			//pod的容器信息
			for _, container := range pod.Spec.Containers {
				c := param.PodContainer{
					Name:  container.Name,
					Image: container.Image,
					ContainerSpec: param.ContainerSpec{
						ImagePullPolicy: string(container.ImagePullPolicy),
						Command:         container.Command,
						Args:            container.Args,
					},
				}
				for _, status := range pod.Status.ContainerStatuses {
					if container.Name == status.Name {
						c.ContainerStatus.Ready = status.Ready
						c.ContainerStatus.RestartCount = status.RestartCount
						c.ContainerStatus.Image = status.Image
						c.ContainerStatus.ImageID = status.ImageID
						c.ContainerStatus.ContainerID = status.ContainerID
						c.ContainerStatus.Started = *status.Started
					}
					if c.ContainerStatus.Started {
						c.ContainerStatus.StartedAt = status.State.Running.StartedAt.Format("2006-01-02 15:04:05")
					}
				}
				p.Containers = append(p.Containers, c)
			}
			response.Pods = append(response.Pods, p)
		}
	}
	for _, conditions := range deployment.Status.Conditions {
		response.Conditions = append(response.Conditions, param.DeploymentConditions{
			Type:           string(conditions.Type),
			Status:         string(conditions.Status),
			LastUpdateTime: conditions.LastUpdateTime.Format("2006-01-02 15:04:05"),
			Reason:         conditions.Reason,
			Message:        conditions.Message,
		})
	}
	return
}

//CreateKubernetesDeployment 创建k8s部署
func (pointer *API) CreateKubernetesDeployment(ctx plugins.Context, request param.CreateKubernetesDeploymentReq) (response param.CreateKubernetesDeploymentRes, err error) {
	deployment := new(v1.Deployment)
	deployment.Name = request.Name
	deployment.Namespace = request.NameSpace
	label := map[string]string{
		"k8s.kuboard.cn/name":  request.Name,
		"k8s.kuboard.cn/layer": "web",
	}
	deployment.Labels = label
	deployment.Spec.Replicas = &request.Replicas
	deployment.Spec.Template.Labels = label
	deployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: label,
	}
	//容器设置
	for _, container := range request.Containers {
		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, corev1.Container{
			Name:            container.Name,
			Image:           container.Image,
			ImagePullPolicy: corev1.PullPolicy(container.ContainerSpec.ImagePullPolicy),
			Command:         container.ContainerSpec.Command,
			Args:            container.ContainerSpec.Args,
		})
	}
	if _, e := pointer.clientSet.AppsV1().Deployments(request.NameSpace).Create(deployment); e != nil {
		if !k8serrors.IsAlreadyExists(e) {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.CreateKubernetesDeploymentErr, "创建部署失败")
			return
		}
		//删除部署
		if e := pointer.clientSet.AppsV1().Deployments(request.NameSpace).Delete(deployment.Name, &metav1.DeleteOptions{}); e != nil {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.CreateKubernetesDeploymentErr, "创建部署失败")
			return
		}
		//重新部署
		if _, e := pointer.clientSet.AppsV1().Deployments(request.NameSpace).Create(deployment); e != nil {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.CreateKubernetesDeploymentErr, "创建部署失败")
			return
		}
	}
	//删除服务
	if _, e := pointer.clientSet.CoreV1().Services(request.NameSpace).Get(request.Name, metav1.GetOptions{}); e == nil {
		pointer.clientSet.CoreV1().Services(request.NameSpace).Delete(request.Name, &metav1.DeleteOptions{})
	}
	//设置服务
	if len(request.Service.Ports) > 0 {
		service := new(corev1.Service)
		service.Name = request.Name
		service.Namespace = request.NameSpace
		service.Spec.Type = corev1.ServiceType(request.Service.Type)
		service.Labels = label
		service.Annotations = label
		service.Spec.Selector = label
		for _, port := range request.Service.Ports {
			service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
				Name:       port.Name,
				Protocol:   corev1.Protocol(port.Protocol),
				Port:       port.Port,
				NodePort:   port.NodePort,
				TargetPort: intstr.FromInt(int(port.TargetPort)),
			})
		}
		if _, e := pointer.clientSet.CoreV1().Services(request.NameSpace).Create(service); e != nil {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.CreateKubernetesDeploymentErr, "创建部署失败")
			return
		}
	}
	response.Success = true
	return
}

//DeleteKubernetesDeployment 删除k8s部署
func (pointer *API) DeleteKubernetesDeployment(ctx plugins.Context, request param.DeleteKubernetesDeploymentReq) (response param.DeleteKubernetesDeploymentRes, err error) {
	//删除服务
	if _, e := pointer.clientSet.CoreV1().Services(request.NameSpace).Get(request.Name, metav1.GetOptions{}); e == nil {
		pointer.clientSet.CoreV1().Services(request.NameSpace).Delete(request.Name, &metav1.DeleteOptions{})
	}
	//删除部署
	if e := pointer.clientSet.AppsV1().Deployments(request.NameSpace).Delete(request.Name, &metav1.DeleteOptions{}); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DeleteKubernetesDeploymentErr, "删除k8s部署失败")
		return
	}
	response.Success = true
	return
}

//GetKubernetesPodLog 获取kubernetes的pod日志
func (pointer *API) GetKubernetesPodLog(ctx plugins.Context, request param.GetKubernetesPodLogReq) (response param.GetKubernetesPodLogRes, err error) {
	key := fmt.Sprintf("k8s-pod-%s-%s", request.NameSpace, request.Name)
	count, e := pointer.cache.GetCache().SCard(key)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetKubernetesPodLogErr, "获取日志失败")
		return
	}
	if count > 0 {
		if _, e := pointer.cache.GetCache().Sadd(key, ctx.GetToken()); e != nil {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.GetKubernetesPodLogErr, "获取日志失败")
			return
		}
		response.Success = true
		return
	}
	if request.TailLines <= 0 {
		request.TailLines = 100
	}
	logs, e := pointer.clientSet.CoreV1().Pods(request.NameSpace).GetLogs(request.Name, &corev1.PodLogOptions{
		TailLines: &request.TailLines,
		Follow:    true,
	}).Stream()
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetKubernetesPodLogErr, "获取日志失败")
		return
	}
	//添加到缓存队列
	if _, e := pointer.cache.GetCache().Sadd(key, ctx.GetToken()); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetKubernetesPodLogErr, "获取日志失败")
		logs.Close()
		return
	}
	go func() {
		scanner := bufio.NewScanner(logs)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		fmt.Println("退出")
	}()
	response.Success = true
	return
}
