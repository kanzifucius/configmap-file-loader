package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)


type ConfigMapRetrieve struct {
	namespace string
	client    kubernetes.Interface
	options  metav1.ListOptions
}


func NewConfigMapsRetrieve(namespace , labelFilter string, client kubernetes.Interface) *ConfigMapRetrieve {
	return &ConfigMapRetrieve{
		namespace: namespace,
		client:    client,
		options: metav1.ListOptions {
			LabelSelector: labelFilter, //"promethuesType=rulefile",

		},

	}
}


func (p *ConfigMapRetrieve) GetListerWatcher() cache.ListerWatcher {

	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return p.client.CoreV1().ConfigMaps(p.namespace).List(p.options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return p.client.CoreV1().ConfigMaps(p.namespace).Watch(p.options)
		},
	}
}


func (p *ConfigMapRetrieve) GetObject() runtime.Object {
	return &corev1.ConfigMap{}
}