package main

import (
	"encoding/json"
	"testing"
)

func TestParseSpeedTestResult(t *testing.T) {
	b := []byte(`
	{
		"type": "result",
		"timestamp": "2020-10-24T01:32:34Z",
		"ping": {
			"jitter": 0.083000000000000004,
			"latency": 3.222
		},
		"download": {
			"bandwidth": 53038114,
			"bytes": 435642344,
			"elapsed": 8312
		},
		"upload": {
			"bandwidth": 89205892,
			"bytes": 429968780,
			"elapsed": 4808
		},
		"packetLoss": 0,
		"isp": "JPNE",
		"interface": {
			"internalIp": "10.0.0.1",
			"name": "eth0",
			"macAddr": "FF:FF:FF:FF:FF:FF",
			"isVpn": false,
			"externalIp": "1.1.1.1"
		},
		"server": {
			"id": 14623,
			"name": "IPA CyberLab",
			"location": "Bunkyo",
			"country": "Japan",
			"host": "speed.coe.ad.jp",
			"port": 8080,
			"ip": "103.95.184.74"
		},
		"result": {
			"id": "b3d6bd12-1ef1-455d-9fc0-157f49d69e14",
			"url": "https://www.speedtest.net/result/c/b3d6bd12-1ef1-455d-9fc0-157f49d69e14"
		}
	}`)

	want := SpeedTestResult{
		Timestamp: "2020-10-24T01:32:34Z", Ping: struct {
			Jitter  float64 "json:\"jitter\""
			Latency float64 "json:\"latency\""
		}{Jitter: 0.083, Latency: 3.222}, Download: struct {
			Bandwidth float64 "json:\"bandwidth\""
			Bytes     float64 "json:\"bytes\""
			Elapsed   float64 "json:\"elapsed\""
		}{Bandwidth: 53038114, Bytes: 435642344, Elapsed: 8312},
		Upload: struct {
			Bandwidth float64 "json:\"bandwidth\""
			Bytes     float64 "json:\"bytes\""
			Elapsed   float64 "json:\"elapsed\""
		}{Bandwidth: 89205892, Bytes: 429968780, Elapsed: 4808},
		PacketLoss: 0,
		Isp:        "JPNE",
		Interface: struct {
			InternalIp string "json:\"internalIp\""
			Name       string "json:\"name\""
			MacAddr    string "json:\"macAddr\""
			IsVpn      bool   "json:\"isVpn\""
			ExternalIp string "json:\"externalIp\""
		}{InternalIp: "10.0.0.1", Name: "eth0", MacAddr: "FF:FF:FF:FF:FF:FF", IsVpn: false, ExternalIp: "1.1.1.1"},
		Server: struct {
			Id       uint64 "json:\"id\""
			Name     string "json:\"name\""
			Location string "json:\"location\""
			Country  string "json:\"country\""
			Host     string "json:\"host\""
			Port     uint16 "json:\"port\""
			Ip       string "json:\"ip\""
		}{Id: 14623, Name: "IPA CyberLab", Location: "Bunkyo", Country: "Japan", Host: "speed.coe.ad.jp", Port: 8080, Ip: "103.95.184.74"},
		Result: struct {
			Id  string "json:\"id\""
			Url string "json:\"url\""
		}{Id: "b3d6bd12-1ef1-455d-9fc0-157f49d69e14", Url: "https://www.speedtest.net/result/c/b3d6bd12-1ef1-455d-9fc0-157f49d69e14"},
	}

	result := new(SpeedTestResult)

	if err := json.Unmarshal(b, result); err != nil {
		t.Fatalf("json.Unmarshal failed, %v", err)
	}

	if *result != want {
		t.Fatalf(`result = %#v, want = %+v`, *result, want)
	}
}
