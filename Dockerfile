FROM bentou/tensorflowgo

# "Design software to be easy to use correctly and hard to use incorrectly"

# maintener info
MAINTAINER Lukasz Pyrzyk <lukasz.pyrzyk@gmail.com>, Jakub Bentkowski <bentkowski.jakub@gmail.com>

# copy all files
COPY ./Src/Server /go/src/Cerber

# install go application and its dependencies
RUN go get github.com/op/go-logging & go get github.com/streadway/amqp & go get github.com/golang/protobuf/proto & go get gopkg.in/mgo.v2 & go get gopkg.in/yaml.v2 & go get -d github.com/tensorflow/tensorflow/tensorflow/go

# Apply patches
COPY ./Patches $GOPATH/src/github.com/tensorflow/tensorflow/tensorflow/go

# Add repo with protobuf 3 and install it
RUN add-apt-repository ppa:maarten-fonville/protobuf
RUN apt-get update && apt-get install -y --allow-unauthenticated protobuf-compiler python-protobuf

# Generate Tensorflow go bindings
RUN go get github.com/tensorflow/tensorflow/tensorflow/go/op & go generate -v -x github.com/tensorflow/tensorflow/tensorflow/go/op

# Install Cerber
RUN go install Cerber

# set entrypoint to the docker run
ENTRYPOINT ["/go/bin/Cerber"]
