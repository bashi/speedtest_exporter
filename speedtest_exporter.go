package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Flags

	speedTestBinary = flag.String("speedTestBinary", "speedtest", "Path to the `speedtest` command")

	addr = flag.String("addr", ":9300", "Listen address")

	intervalString = flag.String("interval", "30m", "Perform speedtest every interval")
	interval       = 30 * time.Minute

	useMockResponse = flag.Bool("useMockResponse", false, "Use mock response")

	// Prometheus metrics

	pingJitter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_ping_jitter",
		Help: "SpeedTest ping jitter",
	})
	pingLatency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_ping_latency",
		Help: "SpeedTest ping latency",
	})

	downloadBandwidth = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_download_bandwidth",
		Help: "SpeedTest download bandwidth",
	})
	downloadBytes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_download_bytes",
		Help: "SpeedTest download bytes",
	})
	downloadElapsed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_download_elapsed",
		Help: "SpeedTest download elapsed",
	})

	uploadBandwidth = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_upload_bandwidth",
		Help: "SpeedTest upload bandwidth",
	})
	uploadBytes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_upload_bytes",
		Help: "SpeedTest upload bytes",
	})
	uploadElapsed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_upload_elapsed",
		Help: "SpeedTest upload elapsed",
	})
)

func runSpeedTest() ([]byte, error) {
	if *useMockResponse {
		out := []byte(`{"type":"result","timestamp":"2020-10-24T01:32:34Z","ping":{"jitter":0.083000000000000004,"latency":3.222},"download":{"bandwidth":53038114,"bytes":435642344,"elapsed":8312},"upload":{"bandwidth":89205892,"bytes":429968780,"elapsed":4808},"packetLoss":0,"isp":"JPNE","interface":{"internalIp":"192.168.100.11","name":"eth0","macAddr":"00:15:5D:0B:0C:15","isVpn":false,"externalIp":"106.72.179.96"},"server":{"id":14623,"name":"IPA CyberLab","location":"Bunkyo","country":"Japan","host":"speed.coe.ad.jp","port":8080,"ip":"103.95.184.74"},"result":{"id":"ee0d6f53-dbed-44d1-a231-cad85743cde3","url":"https://www.speedtest.net/result/c/ee0d6f53-dbed-44d1-a231-cad85743cde3"}}`)
		log.Printf("%s", out)
		return out, nil
	}

	out, err := exec.Command(*speedTestBinary, "--accept-license", "--precision=0", "--format=json", "--progress=no").Output()
	log.Printf("%s", out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func recordSpeedTest() {
	go func() {
		for {
			out, err := runSpeedTest()
			if err != nil {
				log.Fatal(err)
			}

			result, err := parseSpeedTestResult(out)
			if err != nil {
				result, err := parseSpeedTestErrorResult(out)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("ERROR: %v", result)
				time.Sleep(interval)
				continue
			}

			pingJitter.Set(result.Ping.Jitter)
			pingLatency.Set(result.Ping.Latency)

			downloadBandwidth.Set(result.Download.Bandwidth)
			downloadBytes.Set(result.Download.Bytes)
			downloadElapsed.Set(result.Download.Elapsed)

			uploadBandwidth.Set(result.Upload.Bandwidth)
			uploadBytes.Set(result.Upload.Bytes)
			uploadElapsed.Set(result.Upload.Elapsed)

			time.Sleep(interval)
		}
	}()
}

type SpeedTestResult struct {
	Timestamp string `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
	} `json:"ping"`
	Download struct {
		Bandwidth float64 `json:"bandwidth"`
		Bytes     float64 `json:"bytes"`
		Elapsed   float64 `json:"elapsed"`
	} `json:"download"`
	Upload struct {
		Bandwidth float64 `json:"bandwidth"`
		Bytes     float64 `json:"bytes"`
		Elapsed   float64 `json:"elapsed"`
	} `json:"upload"`
	PacketLoss float64 `json:"packetLoss"`
	Isp        string  `json:"isp"`
	Interface  struct {
		InternalIp string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVpn      bool   `json:"isVpn"`
		ExternalIp string `json:"externalIp"`
	} `json:"interface"`
	Server struct {
		Id       uint64 `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		Host     string `json:"host"`
		Port     uint16 `json:"port"`
		Ip       string `json:"ip"`
	} `json:"server"`
	Result struct {
		Id  string `json:"id"`
		Url string `json:"url"`
	} `json:"result"`
}

func parseSpeedTestResult(b []byte) (*SpeedTestResult, error) {
	result := new(SpeedTestResult)
	if err := json.Unmarshal(b, result); err != nil {
		return nil, err
	}
	return result, nil
}

type SpeedTestErrorResult struct {
	Error string `json:"error"`
}

func parseSpeedTestErrorResult(b []byte) (*SpeedTestErrorResult, error) {
	result := new(SpeedTestErrorResult)
	if err := json.Unmarshal(b, result); err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	log.SetPrefix("speedtest_exporter: ")

	flag.Parse()

	parsedInterval, err := time.ParseDuration(*intervalString)
	if err != nil {
		log.Fatal(err)
	}
	interval = parsedInterval

	recordSpeedTest()

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %v, interval %v", *addr, interval)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
