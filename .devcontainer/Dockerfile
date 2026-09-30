
FROM rockylinux/rockylinux:10 AS base

ARG REPO_NAME=go-play
ARG USERNAME=guido
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Add ca certs so we can build the image on vpn
COPY cas.pem /etc/pki/ca-trust/source/anchors/custom_cas.pem 
RUN update-ca-trust

RUN set -eux -o pipefail && \
    dnf install --setopt=***.sslverify=false -y ca-certificates curl-minimal git jq procps sudo wget which rsync

    # envs
ENV INSTALL_PATH="/usr/local"
ENV PATH=${HOME}/bin:${PATH}

# Create the user
RUN groupadd --gid ${USER_GID} ${USERNAME} && \
  useradd --uid ${USER_UID} --gid ${USER_GID} -m ${USERNAME} -s /bin/zsh && \
  echo ${USERNAME} ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/${USERNAME} && \
  chmod 0440 /etc/sudoers.d/${USERNAME}

########### NON ROOT FROM HERE DOWN ###########
FROM base AS user

USER ${USERNAME}

RUN set -eux && \
    sudo mkdir -p /workdir/${REPO_NAME} && \
    sudo chown ${USERNAME}:${USERNAME} /workdir

# envs
ENV HOME="/home/${USERNAME}"
ENV SDK_HOME=${HOME}/sdk
ENV PATH=${HOME}/bin:${PATH}

# go
ENV GO_VERSION="1.26.5"
ENV GOPATH=${SDK_HOME}/go
ENV PATH=${GOPATH}/bin:${HOME}/sdk/go:${PATH}
RUN set -eux && \
    export GO_ARCH="$(arch | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')" && \
    cd /tmp && mkdir -p ${SDK_HOME} && \
    curl -k -o go.tar.gz -L https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz && \
    tar xzf go.tar.gz && mv go ${SDK_HOME} && rm /tmp/go.tar.gz && \
    cd ${SDK_HOME} && mv go go${GO_VERSION} && ln -s go${GO_VERSION} go
   
# starship — config lives outside ~/.config because that dir is bind-mounted from the host
RUN set -eux && \
    mkdir -p "${HOME}/bin" && \
    export ASSET="starship-$(uname -m)-unknown-linux-musl" && \
    curl -fsSL -o "/tmp/${ASSET}.tar.gz" "https://github.com/starship/starship/releases/latest/download/${ASSET}.tar.gz" && \
    tar xzf "/tmp/${ASSET}.tar.gz" -C "${HOME}/bin" starship && \
    chmod +x "${HOME}/bin/starship" && \
    rm "/tmp/${ASSET}.tar.gz" && \
    starship --version
ENV STARSHIP_CONFIG="${HOME}/.starship.toml"
COPY --chown=${USERNAME}:${USERNAME} starship.toml ${HOME}/.starship.toml

# pretty bash
COPY --chown=${USERNAME}:${USERNAME} .bash_pretty ${HOME}/.bash_pretty
RUN set -eux && \
    echo "source ${HOME}/.bash_pretty" >> ${HOME}/.bashrc

