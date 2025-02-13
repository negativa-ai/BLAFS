FROM ubuntu:20.04

# configure timezone
ENV TZ=Europe/Stockholm
RUN ln -snf "/usr/share/zoneinfo/$TZ" /etc/localtime
RUN echo "$TZ" > /etc/timezone

RUN apt-get update
RUN apt install -y git curl wget vim cmake pkg-config
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

# install on-my-zsh
RUN apt-get install zsh -y
RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

RUN wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz && rm go1.23.6.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# install libfuse
RUN apt-get install libfuse3-dev -y

# install python3.10
RUN apt install software-properties-common -y
RUN add-apt-repository ppa:deadsnakes/ppa -y
RUN apt update
RUN apt install python3.10 python3.10-venv python3.10-dev -y
RUN ls -la /usr/bin/python3
RUN rm /usr/bin/python3
RUN ln -s python3.10 /usr/bin/python3
RUN curl -sS https://bootstrap.pypa.io/get-pip.py | python3.10
RUN pip install redis
RUN pip install pyyaml

WORKDIR /home/ubuntu/repos/BAFFS

ENTRYPOINT ["/sbin/init"]

CMD ["dockerd"]
