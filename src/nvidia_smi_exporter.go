package main

import (
    "io"
    "net/http"
    "encoding/xml"
    "os/exec"
    "log"
    "os"
    "regexp"
    "strconv"
)

const LISTEN_ADDRESS = ":9202"
const NVIDIA_SMI_PATH = "/usr/bin/nvidia-smi"

var testMode string;
var listenAddress string;

type NvidiaSmiLog struct {
    DriverVersion string `xml:"driver_version"`
    AttachedGPUs string `xml:"attached_gpus"`
    GPUs []struct {
        ProductName string `xml:"product_name"`
        ProductBrand string `xml:"product_brand"`
        UUID string `xml:"uuid"`
        FanSpeed string `xml:"fan_speed"`
        PCI struct {
            PCIBus string `xml:"pci_bus"`
        } `xml:"pci"`
        FbMemoryUsage struct {
            Total string `xml:"total"`
            Used string `xml:"used"`
            Free string `xml:"free"`
        } `xml:"fb_memory_usage"`
        Utilization struct {
            GPUUtil string `xml:"gpu_util"`
            MemoryUtil string `xml:"memory_util"`
        } `xml:"utilization"`
        Temperature struct {
            GPUTemp string `xml:"gpu_temp"`
            GPUTempMaxThreshold string `xml:"gpu_temp_max_threshold"`
            GPUTempSlowThreshold string `xml:"gpu_temp_slow_threshold"`
        } `xml:"temperature"`
        PowerReadings struct {
            PowerDraw string `xml:"power_draw"`
            PowerLimit string `xml:"power_limit"`
        } `xml:"power_readings"`
        Clocks struct {
            GraphicsClock string `xml:"graphics_clock"`
            SmClock string `xml:"sm_clock"`
            MemClock string `xml:"mem_clock"`
            VideoClock string `xml:"video_clock"`
        } `xml:"clocks"`
        MaxClocks struct {
            GraphicsClock string `xml:"graphics_clock"`
            SmClock string `xml:"sm_clock"`
            MemClock string `xml:"mem_clock"`
            VideoClock string `xml:"video_clock"`
        } `xml:"max_clocks"`
    } `xml:"gpu"`
}

func formatValue(key string, meta string, value string) string {
    result := key;
    if (meta != "") {
        result += "{" + meta + "}";
    }
    return result + " " + value +"\n"
}

func filterNumber(value string) string {
    r := regexp.MustCompile("[^0-9.]")
    return r.ReplaceAllString(value, "")
}

func metrics(w http.ResponseWriter, r *http.Request) {
    log.Print("Serving /metrics")

    var cmd *exec.Cmd
    if (testMode == "1") {
        dir, err := os.Getwd()
        if err != nil {
            log.Fatal(err)
        }
        cmd = exec.Command("/bin/cat", dir + "/test.xml")
    } else {
        cmd = exec.Command(NVIDIA_SMI_PATH, "-q", "-x")
    }

    // Execute system command
    stdout, err := cmd.Output()
    if err != nil {
        println(err.Error())
        return
    }

    // Parse XML
    var xmlData NvidiaSmiLog
    xml.Unmarshal(stdout, &xmlData)

    // Output
    io.WriteString(w, formatValue("nvidia_driver_info", "version=\"" + xmlData.DriverVersion + "\"", filterNumber("1")))
    io.WriteString(w, formatValue("nvidia_device_count", "", filterNumber(xmlData.AttachedGPUs)))
    for i, GPU := range xmlData.GPUs {
        var idx = strconv.Itoa(i)
        io.WriteString(w, formatValue("nvidia_fanspeed", "minor=\"" + idx + "\"", filterNumber(GPU.FanSpeed)))
        io.WriteString(w, formatValue("nvidia_info", "index=\"" + idx + "\", minor=\"" + idx + "\", uuid=\"" + GPU.UUID + "\", name=\"" + GPU.ProductName + "\"", filterNumber("1")))
        io.WriteString(w, formatValue("nvidia_memory_total", "minor=\"" + idx + "\"", filterNumber(GPU.FbMemoryUsage.Total)))
        io.WriteString(w, formatValue("nvidia_memory_used", "minor=\"" + idx + "\"", filterNumber(GPU.FbMemoryUsage.Used)))
        io.WriteString(w, formatValue("nvidia_memory_free", "minor=\"" + idx + "\"", filterNumber(GPU.FbMemoryUsage.Free)))
        io.WriteString(w, formatValue("nvidia_utilization_gpu", "minor=\"" + idx + "\"", filterNumber(GPU.Utilization.GPUUtil)))
        io.WriteString(w, formatValue("nvidia_utilization_memory", "minor=\"" + idx + "\"", filterNumber(GPU.Utilization.MemoryUtil)))
        io.WriteString(w, formatValue("nvidia_temperatures", "minor=\"" + idx + "\"", filterNumber(GPU.Temperature.GPUTemp)))
        io.WriteString(w, formatValue("nvidia_temperatures_max", "minor=\"" + idx + "\"", filterNumber(GPU.Temperature.GPUTempMaxThreshold)))
        io.WriteString(w, formatValue("nvidia_temperatures_slow", "minor=\"" + idx + "\"", filterNumber(GPU.Temperature.GPUTempSlowThreshold)))
        io.WriteString(w, formatValue("nvidia_power_usage", "minor=\"" + idx + "\"", filterNumber(GPU.PowerReadings.PowerDraw)))
        io.WriteString(w, formatValue("nvidia_power_limit", "minor=\"" + idx + "\"", filterNumber(GPU.PowerReadings.PowerLimit)))
        io.WriteString(w, formatValue("nvidia_clock_graphics", "minor=\"" + idx + "\"", filterNumber(GPU.Clocks.GraphicsClock)))
        io.WriteString(w, formatValue("nvidia_clock_graphics_max", "minor=\"" + idx + "\"", filterNumber(GPU.MaxClocks.GraphicsClock)))
        io.WriteString(w, formatValue("nvidia_clock_sm", "minor=\"" + idx + "\"", filterNumber(GPU.Clocks.SmClock)))
        io.WriteString(w, formatValue("nvidia_clock_sm_max", "minor=\"" + idx + "\"", filterNumber(GPU.MaxClocks.SmClock)))
        io.WriteString(w, formatValue("nvidia_clock_mem", "minor=\"" + idx + "\"", filterNumber(GPU.Clocks.MemClock)))
        io.WriteString(w, formatValue("nvidia_clock_mem_max", "minor=\"" + idx + "\"", filterNumber(GPU.MaxClocks.MemClock)))
        io.WriteString(w, formatValue("nvidia_clock_video", "minor=\"" + idx + "\"", filterNumber(GPU.Clocks.VideoClock)))
        io.WriteString(w, formatValue("nvidia_clock_video_max", "minor=\"" + idx + "\"", filterNumber(GPU.MaxClocks.VideoClock)))
    }
}

func index(w http.ResponseWriter, r *http.Request) {
    log.Print("Serving /index")
    html := `<!doctype html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Nvidia SMI Exporter</title>
    </head>
    <body>
        <h1>Nvidia SMI Exporter</h1>
        <p><a href="/metrics">Metrics</a></p>
    </body>
</html>`
    io.WriteString(w, html)
}

func main() {
    testMode = os.Getenv("TEST_MODE")
    if (testMode == "1") {
        log.Print("Test mode is enabled")
    }

    listenAddress = os.Getenv("LISTEN_ADDRESS")
    if (len(listenAddress) == 0) {
        listenAddress = LISTEN_ADDRESS
    }

    log.Print("Nvidia SMI exporter listening on " + listenAddress)
    http.HandleFunc("/", index)
    http.HandleFunc("/metrics", metrics)
    http.ListenAndServe(listenAddress, nil)
}
