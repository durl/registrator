FROM debian:7.8

MAINTAINER David Url <david@x00.at>

# set environment for apt
ENV TERM linux
ENV DEBIAN_FRONTEND noninteractive

# install, build, cleanup
ENV BUILD_DEPS python-pip=1.1-3 python-dev=2.7.3-4+deb7u1 libssl-dev=1.0.1e-2+deb7u17 libffi-dev=3.0.10-3
ENV RUN_DEPS python=2.7.3-4+deb7u1
RUN apt-get update \
    && apt-cache showpkg $RUN_DEPS $BUILD_DEPS \
    && apt-get install -y python $BUILD_DEPS \
    && pip install docker-py==1.3.0 \
    && pip install python-etcd==0.3.3 \
    && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false -o APT::AutoRemove::SuggestsImportant=false $BUILD_DEPS \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/* \
    && unlink /usr/bin/python \
    && ln -s /usr/bin/python2.7 /usr/bin/python

# install application:
ADD registrator.py /registrator.py

VOLUME /var/run/docker.sock

CMD ["python", "/registrator.py"]
