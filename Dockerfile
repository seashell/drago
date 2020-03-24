FROM golang:1.13.8-stretch as debian

RUN curl -sS http://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb http://dl.yarnpkg.com/debian/ stable main" |  tee /etc/apt/sources.list.d/yarn.list && \
    curl -sL http://deb.nodesource.com/setup_13.x | bash - && \
    apt-get install -y nodejs && \
    apt-get update && \
    apt-get remove cmdtest && \
    apt-get install -y yarn 

RUN apt-get install -y gcc-arm-linux-gnueabihf libc6-dev-armhf-cross \
                       gcc-aarch64-linux-gnu libc6-dev-arm64-cross

FROM debian
ARG HOST_UID=${HOST_UID:-4000}
ARG HOST_USER=${HOST_USER:-nodummy}

RUN if [ "${HOST_USER}" != "root" ]; then \
    (adduser -q --gecos "" --home /home/${HOST_USER} --disabled-password -u ${HOST_UID} ${HOST_USER} \
    && chown -R "${HOST_UID}:${HOST_UID}" /home/${HOST_USER}); \
    fi
USER ${HOST_USER}