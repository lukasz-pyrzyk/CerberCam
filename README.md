# CerberCam [![Docker Stars](https://img.shields.io/docker/stars/lukaszpyrzyk/cerbercam.svg)](https://hub.docker.com/r/lukaszpyrzyk/cerbercam/)

## Build:
`go build *.go`

## Usage:
CerberCam supports idea of entry commands. For now, there is only one:

`./main.go -command=receive`

## Used open source projects:

- Go-logging https://github.com/op/go-logging
- MongoDB https://github.com/mongodb/mongo
- MongoDB GO http://labix.org/mgo, made by Gustavo Niemeyer
- MongoDB docker image https://github.com/dockerfile/mongodb, made by Docker team and community
- Protobuf https://github.com/google/protobuf
- Protobuf GO https://github.com/mgravell/protobuf-net
- RabbitMQ https://github.com/rabbitmq/rabbitmq-server
- RabbitMQ GO https://github.com/streadway/amqp
- RabbitMQ docker image https://github.com/docker-library/rabbitmq
- TensorFlow https://github.com/tensorflow/tensorflow
- TensorFlow GO https://github.com/tensorflow/tensorflow/tree/master/tensorflow/go
- TensorFlow docker image https://github.com/tensorflow/tensorflow/tree/master/tensorflow/tools/docker
- Yaml for GO https://github.com/go-yaml/yaml


### Planned ###
- Zstandard https://github.com/facebook/zstd
- Zstandard GO https://github.com/DataDog/zstd
