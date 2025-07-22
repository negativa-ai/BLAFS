#
# Go build stage
#
FROM gcr.io/k8s-minikube/kicbase:v0.0.47 AS builder
WORKDIR /build


RUN apt-get update
RUN apt install -y git wget cmake pkg-config 
RUN apt-get install build-essential libfuse3-dev -y

RUN wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz && rm go1.23.6.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

#
# Clone and make install
# 
RUN git clone https://github.com/negativa-ai/BLAFS
WORKDIR BLAFS
COPY Makefile .
RUN make install 

#
# use the original Minikube image
#
FROM gcr.io/k8s-minikube/kicbase:v0.0.47

RUN apt update

RUN apt install libfuse3-3 -y 

#
# copy the tools build
#

COPY --from=builder /usr/bin/baffs /usr/bin/baffs
COPY --from=builder /usr/bin/debloated_fs /usr/bin/debloated_fs

#
# configure timezone
#
ENV TZ=Europe/Berlin
RUN ln -snf "/usr/share/zoneinfo/$TZ" /etc/localtime
RUN echo "$TZ" > /etc/timezone

