======
README
======

:code:`registrator` is a side-kick registration process which registers running docker containers in :code:`etcd` to be discovered by a discovery service.

Usage
=====

Run a registrator container::

    $ docker build -t registrator .
    $ docker docker run -d -v /var/run/docker.sock:/var/run/docker.sock -e HOST_IP=10.0.0.10 -e ETCD_HOST=10.0.0.11:4001 registrator
    
By default registrator will look for containers on the docker host which expose only one port.
If such a container is found it will be published to the specified :code:`etcd` instance.
If a container has multiple ports published, registrator will try to register the container with port :code:`443` or :code:`80` if present.

Control port selection:
    To get better control over which port gets registered in :code:`etcd', a label can be specified which defines the internal port where the service is listening.
    This should not be necessary if each container only exposes a single port.
    ::
        
        -e PRIMARY_PORT_LABEL=publish_port
    
    A "discoverable" container may use this label to get registered::
    
        $ docker run -d -l publish_port=1000/tcp -p 1000 -p 2000 -p 3000 service-a

Control which containers will be registered:
    To filter "discoverable" containers, you can specify mandatory labels::
    
        -e LABELS=foo,bar
    
    Only containers with all of this labels will be registered.
    The specified labels also get published to :code:`etcd`.

Other Paramaters:
    Set a different ttl for registered services::
    
        -e ETCD_TTL 15
        
    Specify the path to the docker socket:
    
        -e DOCKER_SOCKET unix:///var/run/docker.sock

Key Structure
=============

Containers are grouped by the service they are running.
Currently the image name is used as the service name.

::
    
    {
        "action":"get",
        "node":{
            "key":"/services/service-a/69b49c31688a0bed047aa4ec44305e5c82d80f0121bf97da1b96f06fbe80cd78",
            "dir":true,
            "expiration":"2015-07-15T23:27:08.404316306Z",
            "ttl":14,
            "nodes":[
                {
                    "key":"/services/service-a/69b49c31688a0bed047aa4ec44305e5c82d80f0121bf97da1b96f06fbe80cd78/address",
                    "value":"192.168.59.103:32838",
                    "expiration":"2015-07-15T23:27:08.204150652Z",
                    "ttl":14,
                    "modifiedIndex":3047,
                    "createdIndex":3047
                },
                {
                    "key":"/services/service-a/69b49c31688a0bed047aa4ec44305e5c82d80f0121bf97da1b96f06fbe80cd78/foo",
                    "value":"xyz",
                    "expiration":"2015-07-15T23:27:08.268121624Z",
                    "ttl":14,
                    "modifiedIndex":3049,
                    "createdIndex":3049
                },
                {
                    "key":"/services/registrator/69b49c31688a0bed047aa4ec44305e5c82d80f0121bf97da1b96f06fbe80cd78/bar",
                    "value":"abc",
                    "expiration":"2015-07-15T23:27:08.369949396Z",
                    "ttl":14,
                    "modifiedIndex":3051,
                    "createdIndex":3051
                }
            ],
            "modifiedIndex":3011,
            "createdIndex":3011
        }
    }
