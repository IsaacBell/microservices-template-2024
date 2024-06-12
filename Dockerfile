FROM golang:1.22.2 AS builder

COPY . /src
WORKDIR /src

# use the GOPROXY flag to avoid restrictions
# RUN GOPROXY=https://goproxy.cn make build

RUN make compile

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./server", "-conf", "/data/conf"]
