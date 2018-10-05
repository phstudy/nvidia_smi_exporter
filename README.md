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

# Using Nvidia Docker v2
The `nvidia_smi_exporter` depends on [NVIDIA/nvidia-docker](https://github.com/NVIDIA/nvidia-docker).

    $ docker run -d \
          --runtime=nvidia \
          --name=nvidia_smi_exporter \
          -p 9202:9202 \
          study/nvidia_smi_exporter

You can check the output of `nvidia_smi_exporter` with:

    $ docker logs -f nvidia_smi_exporter

# Metrics Example

    nvidia_driver_info{version="390.87"} 1
    nvidia_device_count 1
    nvidia_fanspeed{minor="0"} 0
    nvidia_info{index="0", minor="0", uuid="GPU-7d57c83c-be3c-fc52-9fe2-6eac2185265a", name="GeForce GTX 1080 Ti"} 1
    nvidia_memory_total{minor="0"} 11177
    nvidia_memory_used{minor="0"} 0
    nvidia_memory_free{minor="0"} 11177
    nvidia_utilization_gpu{minor="0"} 3
    nvidia_utilization_memory{minor="0"} 0
    nvidia_temperatures{minor="0"} 55
    nvidia_temperatures_max{minor="0"} 96
    nvidia_temperatures_slow{minor="0"} 93
    nvidia_power_usage{minor="0"} 25.69
    nvidia_power_limit{minor="0"} 280.00
    nvidia_clock_graphics{minor="0"} 139
    nvidia_clock_graphics_max{minor="0"} 1999
    nvidia_clock_sm{minor="0"} 139
    nvidia_clock_sm_max{minor="0"} 1999
    nvidia_clock_mem{minor="0"} 810
    nvidia_clock_mem_max{minor="0"} 5505
    nvidia_clock_video{minor="0"} 734
    nvidia_clock_video_max{minor="0"} 1620

# Grafana Dashboard
https://grafana.com/dashboards/6387