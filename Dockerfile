FROM debian:7.8

MAINTAINER David Url <david@x00.at>

# set environment for apt
ENV TERM linux
ENV DEBIAN_FRONTEND noninteractive

# install, build, cleanup
RUN apt-get update \
    && apt-get install -y python python-pip python-dev libssl-dev libffi-dev \
    && pip install docker-py==1.3.0 \
    && pip install python-etcd==0.3.3 \
    && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false -o APT::AutoRemove::SuggestsImportant=false python-pip python-dev libssl-dev libffi-dev \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/* \
    && unlink /usr/bin/python \
    && ln -s /usr/bin/python2.7 /usr/bin/python

# install application:
ADD registrator.py /registrator.py

VOLUME /var/run/docker.sock

CMD ["python", "/registrator.py"]
