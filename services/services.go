package services

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"os"
	"strings"
)

func getContainers(client *docker.Client) []*docker.Container {
	allContainers, _ := client.ListContainers(docker.ListContainersOptions{All: false})
	var containers []*docker.Container
	for _, c := range allContainers {
		container, _ := client.InspectContainer(c.ID)
		containers = append(containers, container)
	}
	return containers
}

func GetDockerClient(endpoint string) *docker.Client {
	if strings.HasPrefix(endpoint, "unix://") {
		client, _ := docker.NewClient(endpoint)
		return client
	} else if strings.HasPrefix(endpoint, "tcp://") {
		path := os.Getenv("DOCKER_CERT_PATH")
		ca := fmt.Sprintf("%s/ca.pem", path)
		cert := fmt.Sprintf("%s/cert.pem", path)
		key := fmt.Sprintf("%s/key.pem", path)
		client, _ := docker.NewTLSClient(endpoint, cert, key, ca)
		return client

	}
	return nil
}

type Service struct {
	id      string
	service string
	address string
	labels  map[string]string
}

func parsePort(port string) docker.Port {
	if port == "" {
		return docker.Port(port)
	}
	if strings.Contains(port, "/") {
		return docker.Port(port)
	} else {
		return docker.Port(port + "/tcp")
	}
}

func determineServicePort(container *docker.Container, hostIp string, portLabel string) string {
	portCount := len(container.NetworkSettings.Ports)
	containerPort := parsePort(container.Config.Labels[portLabel])
	var portMapping []docker.PortBinding
	for portKey, mappings := range container.NetworkSettings.Ports {
		if portCount == 1 {
			// use the only exposed port
			portMapping = mappings
		}
		if portKey == containerPort {
			// use (container) port defined by label
			portMapping = mappings
		}
	}

	// choose port from mappings:
	for _, mapping := range portMapping {
		if mapping.HostIP == "0.0.0.0" || mapping.HostIP == hostIp {
			return mapping.HostPort
		}
	}
	// unable to determine port
	return ""
}

func containerToService(container *docker.Container, hostIp string, portLabel string) Service {
	hostPort := determineServicePort(container, portLabel, hostIp)
	return Service{
		id: container.ID,
		service: container.Config.Image,
		address: fmt.Sprintf("%s:%s", hostIp, hostPort),
		labels: container.Config.Labels,
	}
}

func GetDiscoverableServices(dockerClient *docker.Client) []Service {
	var services []Service
	containers := getContainers(dockerClient)
	for _, container := range containers {
		service := containerToService(container, os.Getenv("HOST_IP"), os.Getenv("SERVICE_PORT_LABEL"))
		fmt.Println(container.ID)
		services = append(services, service)
	}
	return services
}

