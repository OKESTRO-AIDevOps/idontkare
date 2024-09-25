FROM ubuntu:20.04

# Change apt mirror to a faster one to speed up the build process
RUN sed -i 's/http:\/\/archive.ubuntu.com/http:\/\/mirror.kakao.com/g' /etc/apt/sources.list && \
    sed -i 's/http:\/\/security.ubuntu.com/http:\/\/mirror.kakao.com/g' /etc/apt/sources.list

RUN apt-get update && apt-get install -y git vim nano make curl wget

# Install Golang
RUN wget https://go.dev/dl/go1.21.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz && \
    rm go1.21.3.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

# Install Kind
RUN curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64 && \
    chmod +x ./kind && \
    mv ./kind /usr/local/bin/kind

ARG USERNAME=worker
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Set up user and group
RUN apt-get install -y sudo && \
    groupadd --gid $USER_GID $USERNAME && \
    useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    echo "$USERNAME ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

USER $USERNAME
