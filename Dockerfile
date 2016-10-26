FROM golang

# maintener info
MAINTAINER Lukasz Pyrzyk <lukasz.pyrzyk@gmail.com>

# copy all files
COPY ./Src/Server /go/src/Cerber

# install go application and its dependencies
RUN go get github.com/op/go-logging & go get github.com/streadway/amqp & go get github.com/golang/protobuf/proto
RUN go install Cerber

# set entrypoint to the docker run
ENTRYPOINT ["/go/bin/Cerber"]