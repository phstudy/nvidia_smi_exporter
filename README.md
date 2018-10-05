# Nvidia SMI Exporter

Nvidia SMI Exporter for GPU statistics from nvidia-smi written in Go.

Supports multiple GPUs.

# Building and running

Prerequisites:

* [Go compiler](https://golang.org/dl/)

Building:

    go get github.com/phstudy/nvidia_smi_exporter/src
    cd ${GOPATH-$HOME/go}/src/github.com/phstudy/nvidia_smi_exporter
    go build -v -o bin/nvidia_smi_exporter src/nvidia_smi_exporter.go

Running with default port:

    bin/nvidia_smi_exporter

Running with specified port:

    LISTEN_ADDRESS=9201 bin/nvidia_smi_exporter

Running in test mode:

    TEST_MODE=1 bin/nvidia_smi_exporter


# Using Nvidia Docker
The `nvidia_smi_exporter` depends on [NVIDIA/nvidia-docker](https://github.com/NVIDIA/nvidia-docker).

    $ docker run -d \
          --runtime=nvidia \
          --name=nvidia_smi_exporter \
          -p 9202:9202 \
          study/nvidia_smi_exporter

You can check the output of `nvidia_smi_exporter` with:

    $ docker logs -f nvidia_smi_exporter