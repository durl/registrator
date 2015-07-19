package main

import (
	"fmt"
	"time"
	"registrator/services"
	"registrator/registry"
	"github.com/docopt/docopt-go"
)

const Version string = "0.1.0"

func main() {
	usage := `Register running services in etcd.

Usage:
  registrator --hostip=<ip> --etcd=<ip:port> [--docker=<url>] [--ttl=<ttl>]
  registrator -h | --help
  registrator --version

Options:
  -h --help             Show this screen.
  --version             Show version.
  --host=<ip>           IP address of the host.
  --etcd=<url>          URL of the etcd host.
  --docker=<url>        URL of the docker daemon endpoint. [default: unix:///var/run/docker.sock]
  --ttl=<ttl>           How long the service registration is valid. [default: 15s]
`

	// parse arguments:
	args, _ := docopt.Parse(usage, nil, true, Version, false)
	fmt.Println(args)

	// read arguments:
	dockerEndpoint := args["--docker"].(string)
	ttl, _ := time.ParseDuration(string(args["--ttl"].(string)))
	machines := []string{string(args["--etcd"].(string))}

	// setup clients:
	dockerClient := services.GetDockerClient(dockerEndpoint)
	etcdClient := registry.GetEtcdClient(machines)


	// run main loop:
	for {
		services := services.GetDiscoverableServices(dockerClient)
		fmt.Println(services)
		time.Sleep(ttl/2)
		registry.RegisterServices(etcdClient, services, ttl)
		//registerServices(services)
	}

}
