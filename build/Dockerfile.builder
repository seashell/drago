FROM golang:1.17.0-alpine3.14 as builder

ARG HOST_UID=${HOST_UID}
ARG HOST_USER=${HOST_USER}

RUN apk add nodejs npm && \
    npm install -g yarn

RUN if [ "${HOST_USER}" != "root" ]; then \
    (adduser --gecos "" --home /home/${HOST_USER} --disabled-password -u ${HOST_UID} ${HOST_USER} \
    && chown -R "${HOST_UID}:${HOST_UID}" /home/${HOST_USER}); \
    fi

USER ${HOST_USER}