package api

import (
	"fmt"
	"strings"

	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	var labelSelector []string
	for key, value := range deployment.Labels {
		labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", key, value))
	}
	listOptions := metav1.ListOptions{
		LabelSelector: strings.Join(labelSelector, ","),
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
