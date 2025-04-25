FROM ubuntu:20.04

# configure timezone
ENV TZ=Europe/Stockholm
RUN ln -snf "/usr/share/zoneinfo/$TZ" /etc/localtime
RUN echo "$TZ" > /etc/timezone

RUN apt-get update
RUN apt install -y git curl wget cmake pkg-config
RUN apt-get install build-essential -y

# install docker
RUN apt-get install ca-certificates curl &&  install -m 0755 -d /etc/apt/keyrings \
    && curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc \
    && chmod a+r /etc/apt/keyrings/docker.asc
RUN echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
    tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update && apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y

RUN wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz && rm go1.23.6.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# install libfuse
RUN apt-get install libfuse3-dev -y


# Clone the BAFFS repository
RUN git clone https://github.com/negativa-ai/BLAFS.git && cd BLAFS && make install

ENTRYPOINT ["/sbin/init"]

CMD ["dockerd"]
