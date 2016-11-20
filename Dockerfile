FROM bentou/tensorflowgo

# maintener info
MAINTAINER Lukasz Pyrzyk <lukasz.pyrzyk@gmail.com>, Jakub Bentkowski <bentkowski.jakub@gmail.com>

#COPY ./Libs /lib
#RUN echo $(ls /lib | grep tensorflow)

RUN echo $(echo $GOPATH)

# copy all files
COPY ./Src/Server /go/src/Cerber

# install go application and its dependencies
RUN go get github.com/op/go-logging & go get github.com/streadway/amqp & go get github.com/golang/protobuf/proto & go get gopkg.in/mgo.v2 & go get -d github.com/tensorflow/tensorflow/tensorflow/go

#ENV LDFLAGS "-lstdc++ -lm -lgnustl_static"

# Patch motherfuckers
COPY ./Patches $GOPATH/src/github.com/tensorflow/tensorflow/tensorflow/go

#RUN echo $(uname -a)

#RUN echo $(ls $GOPATH/src/github.com/tensorflow/tensorflow/tensorflow/go/#genop/internal)

#ENV GOGCCFLAGS "-lstdc++ -lm -lgnustl_static -lpthreads -lc -lgcc -lgc #GOGCCFLAGS"

#RUN echo $(go env)

#RUN echo $(g++ --version)
#RUN echo $(g++ -v)

# add repo with protobuf 3
RUN add-apt-repository ppa:maarten-fonville/protobuf
RUN apt-get update && apt-get install -y --allow-unauthenticated protobuf-compiler python-protobuf

#RUN bazel build -c opt @protobuf//:protoc
RUN echo $(which protoc)

#-linkshared -ldflags=-lstdc++ 
RUN go get github.com/tensorflow/tensorflow/tensorflow/go/op & go generate -v -x github.com/tensorflow/tensorflow/tensorflow/go/op

RUN go install -v -x Cerber

# set entrypoint to the docker run
ENTRYPOINT ["/go/bin/Cerber"]
