package cloud_watcher

import (
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/dysnix/ai-scale-libs/external/configs"
)

type K8sCloudWatcher struct {
	clientSet   *kubernetes.Clientset
	logger      *zap.SugaredLogger
	conf        *configs.K8sCloudWatcher
	closeChList []chan struct{}
}

func NewK8sCloudWatcher(logger *zap.SugaredLogger, conf *configs.K8sCloudWatcher) (cw *K8sCloudWatcher, err error) {
	cw = &K8sCloudWatcher{
		logger: logger,
		conf:   conf,
	}

	cw.clientSet, err = cw.getClient(conf.CtxPath)
	if err != nil {
		return nil, err
	}

	return cw, nil
}

func (cw *K8sCloudWatcher) Run() {
	for _, informerConf := range cw.conf.Informers {
		//Regular informer Nodes
		watchList := cache.NewListWatchFromClient(cw.clientSet.RESTClient(), informerConf.Resource, v1.NamespaceAll, fields.Everything())
		// Shared informer example
		informer := cache.NewSharedIndexInformer(
			watchList,
			cw.convertResourceToObject(informerConf.Resource),
			informerConf.Interval,
			cache.Indexers{},
		)

		// More than one handler can be added...
		//informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		//	AddFunc:    handleNodeAdd,
		//	UpdateFunc: handleNodeUpdate,
		//	DeleteFunc: handleNodeDelete,
		//})

		tmpCloseCh := make(chan struct{}, 1)

		cw.closeChList = append(cw.closeChList, tmpCloseCh)

		informer.Run(tmpCloseCh)
	}
}

func (cw *K8sCloudWatcher) Close() {
	for _, ch := range cw.closeChList {
		ch <- struct{}{}
		close(ch)
	}
}
