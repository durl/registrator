#!/usr/bin/env python

import logging
import os

import docker
import etcd

# configure logging:
logging.basicConfig(
    format='%(asctime)s - %(module)s - %(levelname)s: %(message)s',
    level=logging.INFO)
log = logging.getLogger(__name__)

# defaults:
DOCKER_SOCKET = 'unix:///var/run/docker.sock'
ETCD_TTL = 15


class RegistrationError(Exception):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return str(self.value)


class DockerClient(object):

    def __init__(self, base_url):
        """Creates and links a new docker client to the specified socket."""
        self.client = docker.Client(base_url=base_url)

    def get_containers(self, label_keys=[]):
        """Get all active containers which have all the specified labels."""
        containers = self.client.containers()
        # filter containers with necessary labels:
        containers = [c for c in containers
                      if all(key in c['Labels'] for key in label_keys)]
        return containers

    def get_container_ids(self, label_keys=[]):
        containers = self.get_containers(label_keys=label_keys)
        return [c['Id'] for c in containers]

    def get_container_info(self, id):
        info = self.client.inspect_container(id)
        return info


class ServiceRegistry(object):

    def __init__(self, etcd_host, label_keys, node_address, ttl, port_label):
        # init etcd client:
        host, port = etcd_host.split(":")
        self.client = etcd.Client(host=host, port=int(port))

        # set values:
        self.label_keys = label_keys
        self.node_address = node_address
        self.ttl = ttl
        self.port_label = port_label

    def register(self, container):
        self._register_service(container)
        self._register_labels(container)

    def _register_service(self, container):
        primary_port = _get_primary_port(container)
        if primary_port:
            service_address = ":".join([self.node_address, primary_port])
            self._write(container, "address", service_address)
        else:
            id = container['Id']
            raise RegistrationError("Unable to determine port for: " + id)

    def _register_labels(self, container):
        labels = {k: container['Config']['Labels'][k] for k in self.label_keys}
        for k in labels:
            self._write(container, k, labels[k])

    def _generate_key_prefix(self, container):
        service = container['Config']['Image']
        instance = container['Id']
        return "/services/{service}/{instance}/".format(service=service,
                                                        instance=instance)

    def _write(self, container, key, value):
        prefix = self._generate_key_prefix(container)
        self.client.write(prefix + key, value, ttl=self.ttl)
        # update directory ttl:
        self.client.write(prefix, None, dir=True, prevExist=True, ttl=self.ttl)


def _flatten_ports(ports):
    return {k: ports[k][0] for k in ports if ports[k]}


def _get_primary_port(container, port_label=None):
    ports = _flatten_ports(container['NetworkSettings']['Ports'])
    # only 1 exposed port:
    if len(ports) == 1:
        return ports.values()[0]['HostPort']
    # specified primary port:
    if port_label:
        labels = container['Config']['Labels']
        if port_label in labels and labels[port_label] in ports:
            return ports[labels[port_label]]['HostPort']
    # try https port:
    if '443/tcp' in ports:
        return ports['443/tcp']['HostPort']
    # try http port:
    if '80/tcp' in ports:
        return ports['80/tcp']['HostPort']
    # return None if no port was found:
    return None


def register_active_containers(docker, registry, labels):
    containers = docker.get_container_ids(label_keys=labels)
    for id in containers:
        container = docker.get_container_info(id)
        try:
            registry.register(container)
        except RegistrationError as e:
            logging.error(str(e))
        except Exception as e:
            logging.exception("Error during service registration for: " + id)
        else:
            logging.info("Registered service for: " + id)


if __name__ == "__main__":
    # read environment:
    docker_socket = os.environ.setdefault('DOCKER_SOCKET', DOCKER_SOCKET)
    etcd_host = os.environ['ETCD_HOST']
    host_ip = os.environ['HOST_IP']
    ttl = int(os.environ.setdefault('ETCD_TTL', str(ETCD_TTL)))
    labels = os.environ.setdefault('LABELS', '')
    labels = labels.split(',') if labels != '' else []
    port_label = os.environ.setdefault('PRIMARY_PORT_LABEL', None)

    # init:
    docker_client = DockerClient(docker_socket)
    registry = ServiceRegistry(etcd_host, labels, host_ip, ttl, port_label)

    # run:
    import time
    while True:
        register_active_containers(docker_client, registry, labels)
        time.sleep(ttl/2)
