package registry

import (
	"github.com/coreos/go-etcd/etcd"
	"fmt"
	"registrator/services"
	"time"
)

func GetEtcdClient(etcdHosts []string) *etcd.Client{
	client := etcd.NewClient(etcdHosts)
	return client
}

func RegisterServices(client *etcd.Client, serviceList []services.Service, ttl time.Duration){
	for _, service := range serviceList{
		RegisterService(client, service, ttl)
	}
}

func RegisterService(client *etcd.Client, service services.Service, ttl time.Duration){
	keyPrefix := fmt.Sprintf("/services/%s/%s/", service.Service, service.Id)
	var attributes map[string]string = make(map[string]string)
	attributes["address"] = service.Address
	for key, value := range service.Labels{
		if key == "address"{
			fmt.Println("WARNING: overriding address(%s) with label: %s", service.Address, value)
		}
		attributes[key] = value
	}

	client.SetDir(keyPrefix, uint64(ttl.Seconds()))
	for key, value := range attributes{
		_, err := client.Set(keyPrefix + key, value , uint64(ttl.Seconds()))
		if err != nil{
			fmt.Println(err)
			// abort current registration
			return
		}
	}

}
