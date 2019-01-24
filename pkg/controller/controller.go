package controller

import (
"context"

"github.com/spotahome/kooper/operator/controller"
"k8s.io/apimachinery/pkg/runtime"
"k8s.io/client-go/kubernetes"


"github.com/kanzifucius/configmap_file_loader/pkg/log"
"github.com/kanzifucius/configmap_file_loader/pkg/service"
)

// Controller is a controller that echoes pod events.
type Controller struct {
	controller.Controller

	config Config
	logger log.Logger
}

// New returns a new Echo controller.
func New(config Config, k8sCli kubernetes.Interface, logger log.Logger) (*Controller, error) {

	ret := NewConfigMapsRetrieve(config.Namespace,config.LabelFilter, k8sCli)
	configMapWriterSrv := service.NewConfigMapFileWriter(logger,config.Directory,config.KeyNameFilter,config.Webhook,config.WebhookMethod,config.WebhookStatusCode)
	handler := &handler{configMapWriterSrv: configMapWriterSrv}

	ctrl := controller.NewSequential(config.ResyncPeriod, handler, ret, nil, logger)

	return &Controller{
		Controller: ctrl,
		config:     config,
		logger:     logger,
	}, nil
}


type handler struct {
	configMapWriterSrv service.ConfigMapFileWriter
}

func (h *handler) Add(_ context.Context, obj runtime.Object) error {
	err := h.configMapWriterSrv.WriteConfigMap(obj)
	return err
}
func (h *handler) Delete(_ context.Context, objKey string) error {
	err:= h.configMapWriterSrv.DeleteConfigMap(objKey)
	return err
}