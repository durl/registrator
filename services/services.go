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
	Id      string
	Service string
	Address string
	Labels  map[string]string
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
		if portKey.Port() == containerPort.Port() && portKey.Proto() == containerPort.Proto(){
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

func filterLabels(labels map[string]string, labelKeys []string) map[string]string {
	var filteredLabels map[string]string = make(map[string]string)
	for _, key := range labelKeys {
		if labels[key] != "" {
			filteredLabels[key] = labels[key]
		}
	}
	return filteredLabels
}

func containerToService(container *docker.Container, hostIp string, portLabel string, labelKeys []string) *Service {
	hostPort := determineServicePort(container, hostIp, portLabel)
	labels := filterLabels(container.Config.Labels, labelKeys)
	if hostPort == "" {
		// no viable port detected
		return nil
	}
	if len(labels) != len(labelKeys) {
		// labels are missing
		return nil
	}
	return &Service{
		Id:      container.ID,
		Service: container.Config.Image,
		Address: fmt.Sprintf("%s:%s", hostIp, hostPort),
		Labels:  labels,
	}
}

func GetDiscoverableServices(dockerClient *docker.Client, hostIp string, portLabel string, labelKeys []string) []Service {
	var services []Service
	containers := getContainers(dockerClient)
	for _, container := range containers {
		service := containerToService(container, hostIp, portLabel, labelKeys)
		if service != nil {
			services = append(services, *service)
		}

	}
	return services
}
