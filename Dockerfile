FROM golang:1.15 AS build1
ENV GOPROXY "https://goproxy.cn"
USER root
WORKDIR /root
COPY go.mod go.sum ./
RUN go mod download
ADD . ./

WORKDIR /root/cmd/label
RUN go build -o label .

FROM golang:1.15 AS build2
ENV GOPROXY "https://goproxy.cn"
USER root

WORKDIR /src
# enable modules caching in separate layer
COPY go.mod go.sum ./
RUN go mod download
COPY ./eth-bee/. ./

RUN make binary

FROM debian:10.9-slim

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*; \
    groupadd -r bee --gid 999; \
    useradd -r -g bee --uid 999 --no-log-init -m bee;

RUN apt-get install wget

# make sure mounted volumes have correct permissions
RUN mkdir -p /home/bee/.bee && chown 999:999 /home/bee/.bee

COPY --from=build2 /src/dist/bee /usr/local/bin/bee
COPY --from=build1 /root/cmd/label /usr/local/bin/label

EXPOSE 1633 1634 1635

WORKDIR /home/bee
VOLUME /home/bee/.bee
