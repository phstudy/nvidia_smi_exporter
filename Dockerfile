FROM nvidia/cuda:9.0-base

MAINTAINER Study Hsueh <ph.study@gmail.com>

RUN apt-get update && \
    apt-get -y install golang --no-install-recommends && \
    rm -r /var/lib/apt/lists/*

WORKDIR /go

COPY . .

RUN go build -v -o bin/nvidia_smi_exporter src/nvidia_smi_exporter.go

EXPOSE 9202

CMD ["./bin/nvidia_smi_exporter"]