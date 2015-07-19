package main // import "registrator"

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"registrator/registry"
	"registrator/services"
	"strings"
	"time"
)

const Version string = "0.1.0"

func main() {
	usage := `Register running services in etcd.

Usage:
  registrator --hostip=<ip> --etcd=<ip:port> [--docker=<url>] [--ttl=<ttl>] [--labels=<key,...>] [--portlabel=<key>]
  registrator -h | --help
  registrator --version

Options:
  -h --help             Show this screen.
  --version             Show version.
  --hostip=<ip>         IP address of the host.
  --etcd=<url>          URL of the etcd host.
  --docker=<url>        URL of the docker daemon endpoint. [default: unix:///var/run/docker.sock]
  --ttl=<ttl>           How long the service registration is valid. [default: 15s]
  --labels=<key,...>    Required labels to get registered (comma separated keys).
  --portlabel=<key>     Label which specifies the internal port on which the service is listening.
`

	// parse arguments:
	args, _ := docopt.Parse(usage, nil, true, Version, false)
	fmt.Println(args)

	// read arguments:
	dockerEndpoint := args["--docker"].(string)
	ttl, _ := time.ParseDuration(string(args["--ttl"].(string)))
	machines := []string{string(args["--etcd"].(string))}
	hostIp := args["--hostip"].(string)
	labelKeys := strings.Split(string(args["--labels"].(string)), ",")
	portLabel := string(args["--portlabel"].(string))

	// setup clients:
	dockerClient := services.GetDockerClient(dockerEndpoint)
	etcdClient := registry.GetEtcdClient(machines)

	// run main loop:
	for {
		services := services.GetDiscoverableServices(dockerClient, hostIp, portLabel, labelKeys)
		fmt.Println(services)
		time.Sleep(ttl / 2)
		registry.RegisterServices(etcdClient, services, ttl)
		//registerServices(services)
	}

}
