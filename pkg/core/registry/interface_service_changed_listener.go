package registry

import (
	"reflect"
	"sync"

	"dubbo.apache.org/dubbo-go/v3/common"
	dubboconstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/metadata/service/local"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	"github.com/apache/dubbo-kubernetes/pkg/core/consts"
	"github.com/apache/dubbo-kubernetes/pkg/core/logger"
	gxset "github.com/dubbogo/gost/container/set"
	"github.com/dubbogo/gost/gof/observer"
)

// DubboISDNotifyListener The Service Discovery Changed  Event Listener
type DubboISDNotifyListener struct {
	serviceNames       *gxset.HashSet
	listeners          map[string]registry.NotifyListener
	serviceUrls        map[string][]*common.URL
	revisionToMetadata map[string]*common.MetadataInfo
	allInstances       map[string][]registry.ServiceInstance

	mutex sync.Mutex
}

func NewDubboISDNotifyListener(services *gxset.HashSet) registry.ServiceInstancesChangedListener {
	return &DubboISDNotifyListener{
		serviceNames:       services,
		listeners:          make(map[string]registry.NotifyListener),
		serviceUrls:        make(map[string][]*common.URL),
		revisionToMetadata: make(map[string]*common.MetadataInfo),
		allInstances:       make(map[string][]registry.ServiceInstance),
	}
}

// OnEvent on ServiceInstancesChangedEvent the service instances change event
func (lstn *DubboISDNotifyListener) OnEvent(e observer.Event) error {
	ce, ok := e.(*registry.ServiceInstancesChangedEvent)
	if !ok {
		return nil
	}
	var err error

	lstn.mutex.Lock()
	defer lstn.mutex.Unlock()

	lstn.allInstances[ce.ServiceName] = ce.Instances
	revisionToInstances := make(map[string][]registry.ServiceInstance)
	newRevisionToMetadata := make(map[string]*common.MetadataInfo)
	localServiceToRevisions := make(map[*common.ServiceInfo]*gxset.HashSet)
	protocolRevisionsToUrls := make(map[string]map[*gxset.HashSet][]*common.URL)
	newServiceURLs := make(map[string][]*common.URL)

	logger.Infof("Received instance notification event of service %s, instance list size %s", ce.ServiceName, len(ce.Instances))

	for _, instances := range lstn.allInstances {
		for _, instance := range instances {
			metadataInstance := ConvertToMetadataInstance(instance)
			if metadataInstance.GetMetadata() == nil {
				logger.Warnf("Instance metadata is nil: %s", metadataInstance.GetHost())
				continue
			}
			revision := metadataInstance.GetMetadata()[dubboconstant.ExportedServicesRevisionPropertyName]
			if "0" == revision {
				logger.Infof("Find instance without valid service metadata: %s", metadataInstance.GetHost())
				continue
			}
			subInstances := revisionToInstances[revision]
			if subInstances == nil {
				subInstances = make([]registry.ServiceInstance, 8)
			}
			revisionToInstances[revision] = append(subInstances, metadataInstance)
			metadataInfo := lstn.revisionToMetadata[revision]
			if metadataInfo == nil {
				metadataInfo, err = GetMetadataInfo(metadataInstance, revision)
				if err != nil {
					return err
				}
			}
			metadataInstance.SetServiceMetadata(metadataInfo)
			for _, service := range metadataInfo.Services {
				if localServiceToRevisions[service] == nil {
					localServiceToRevisions[service] = gxset.NewSet()
				}
				localServiceToRevisions[service].Add(revision)
			}

			newRevisionToMetadata[revision] = metadataInfo
		}
		lstn.revisionToMetadata = newRevisionToMetadata

		for serviceInfo, revisions := range localServiceToRevisions {
			revisionsToUrls := protocolRevisionsToUrls[serviceInfo.Protocol]
			if revisionsToUrls == nil {
				protocolRevisionsToUrls[serviceInfo.Protocol] = make(map[*gxset.HashSet][]*common.URL)
				revisionsToUrls = protocolRevisionsToUrls[serviceInfo.Protocol]
			}
			urls := revisionsToUrls[revisions]
			if urls != nil {
				newServiceURLs[serviceInfo.Name] = urls
			} else {
				urls = make([]*common.URL, 0, 8)
				for _, v := range revisions.Values() {
					r := v.(string)
					for _, i := range revisionToInstances[r] {
						if i != nil {
							urls = append(urls, i.ToURLs(serviceInfo)...)
						}
					}
				}
				revisionsToUrls[revisions] = urls
				newServiceURLs[serviceInfo.Name] = urls
			}
		}
		lstn.serviceUrls = newServiceURLs

		for key, notifyListener := range lstn.listeners {
			urls := lstn.serviceUrls[key]
			events := make([]*registry.ServiceEvent, 0, len(urls))
			for _, url := range urls {
				url.SetParam(consts.RegistryType, consts.RegistryInstance)
				events = append(events, &registry.ServiceEvent{
					Action:  remoting.EventTypeAdd,
					Service: url,
				})
			}
			notifyListener.NotifyAll(events, func() {})
		}
	}
	return nil
}

// AddListenerAndNotify add notify listener and notify to listen service event
func (lstn *DubboISDNotifyListener) AddListenerAndNotify(serviceKey string, notify registry.NotifyListener) {
	lstn.listeners[serviceKey] = notify
	urls := lstn.serviceUrls[serviceKey]
	for _, url := range urls {
		url.SetParam(consts.RegistryType, consts.RegistryInstance)
		notify.Notify(&registry.ServiceEvent{
			Action:  remoting.EventTypeAdd,
			Service: url,
		})
	}
}

// RemoveListener remove notify listener
func (lstn *DubboISDNotifyListener) RemoveListener(serviceKey string) {
	delete(lstn.listeners, serviceKey)
}

// GetServiceNames return all listener service names
func (lstn *DubboISDNotifyListener) GetServiceNames() *gxset.HashSet {
	return lstn.serviceNames
}

// Accept return true if the name is the same
func (lstn *DubboISDNotifyListener) Accept(e observer.Event) bool {
	if ce, ok := e.(*registry.ServiceInstancesChangedEvent); ok {
		return lstn.serviceNames.Contains(ce.ServiceName)
	}
	return false
}

// GetPriority returns -1, it will be the first invoked listener
func (lstn *DubboISDNotifyListener) GetPriority() int {
	return -1
}

// GetEventType returns ServiceInstancesChangedEvent
func (lstn *DubboISDNotifyListener) GetEventType() reflect.Type {
	return reflect.TypeOf(&registry.ServiceInstancesChangedEvent{})
}

// GetMetadataInfo get metadata info when MetadataStorageTypePropertyName is null
func GetInterfaceMetadataInfo(instance registry.ServiceInstance, revision string) (*common.MetadataInfo, error) {
	var metadataStorageType string
	var metadataInfo *common.MetadataInfo
	if instance.GetMetadata() == nil {
		metadataStorageType = dubboconstant.DefaultMetadataStorageType
	} else {
		metadataStorageType = instance.GetMetadata()[dubboconstant.MetadataStorageTypePropertyName]
	}
	if metadataStorageType == dubboconstant.RemoteMetadataStorageType {
		remoteMetadataServiceImpl, err := extension.GetRemoteMetadataService()
		if err != nil {
			return &common.MetadataInfo{}, err
		}
		metadataInfo, err = remoteMetadataServiceImpl.GetMetadata(instance)
		if err != nil {
			return &common.MetadataInfo{}, err
		}
	} else {
		var err error
		proxyFactory := extension.GetMetadataServiceProxyFactory(dubboconstant.DefaultKey)
		metadataService := proxyFactory.GetProxy(instance)
		defer metadataService.(*local.MetadataServiceProxy).Invoker.Destroy()
		metadataInfo, err = metadataService.GetMetadataInfo(revision)
		if err != nil {
			return &common.MetadataInfo{}, err
		}
	}
	return metadataInfo, nil
}
