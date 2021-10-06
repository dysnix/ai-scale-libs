package cloud_watcher

import (
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func (cw *K8sCloudWatcher) getClient(pathToCfg string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if pathToCfg == "" {
		cw.logger.Info("Using in cluster config")
		config, err = rest.InClusterConfig()
		// in cluster access
	} else {
		cw.logger.Info("Using out of cluster config")
		config, err = clientcmd.BuildConfigFromFlags("", pathToCfg)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func (cw *K8sCloudWatcher) convertResourceToObject(resource string) runtime.Object {
	return map[string]runtime.Object{
		"nodes": &api.Node{},
	}[resource]
}
