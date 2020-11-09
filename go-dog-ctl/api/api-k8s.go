package api

import (
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
	for _, pod := range deployments.Items {
		response.Pods = append(response.Pods, param.Pod{
			Name:                pod.Name,
			NameSpace:           pod.Namespace,
			ObservedGeneration:  pod.Status.ObservedGeneration,
			Replicas:            pod.Status.Replicas,
			UpdatedReplicas:     pod.Status.UpdatedReplicas,
			ReadyReplicas:       pod.Status.ReadyReplicas,
			AvailableReplicas:   pod.Status.AvailableReplicas,
			UnavailableReplicas: pod.Status.UnavailableReplicas,
			CreateTime:          pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return
}
