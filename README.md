資料夾結構

```bash
prometheus
   ├── docker-compose.yml
   ├── golang
   │   └── Dockerfile
   ├── promApp
   │   ├── go.mod
   │   ├── go.sum
   │   └── main.go
   └── prometheus.yml
```

```docker
#docker-compose.yml
version: '3.2'
services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
    - 9090:9090
    command:
    - --config.file=/etc/prometheus/prometheus.yml
    volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
    environment:
      TZ: Asia/Taipei
  golang:
    build:
      context: ./
      dockerfile: ./golang/Dockerfile
    entrypoint: ./exec
    ports:
      - '2112:2112'
```

```docker
#/golang/Dockerfile
FROM golang:alpine as builder
WORKDIR /app
COPY ./promApp .
RUN go mod download
RUN go build -o exec

FROM alpine
WORKDIR /usr/bin
COPY --from=builder /app/exec .
RUN apk update && apk add tzdata
```

```yaml
//prometheus.yml
scrape_configs:
- job_name: myapp
  scrape_interval: 1s
  static_configs:
  - targets:
    - 192.168.10.100:2112
- job_name: docker desktop
  scrape_interval: 1s
  static_configs:
  - targets:
    - 192.168.10.101:9323
```

```go
// /promApp/main.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}
func switchFlip() {
	go func() {
		i := 0
		for {
			sw.Set(float64(i % 2))
			i++
			if i == 10 {
				i = 0
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	sw = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "on_off_flip",
		Help:        "開開關關翻翻樂",
		ConstLabels: map[string]string{},
	})
)

func main() {
	recordMetrics()
	switchFlip()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}
```

![1673104871078.jpg](Docker%2037983de704874fa6b1dcc661a56e0858/1673104871078.jpg)

![1673104946958.jpg](Docker%2037983de704874fa6b1dcc661a56e0858/1673104946958.jpg)

## AlertManager

資料夾結構

```
.
├── alertmanager
│   └── config.yml <- add
├── config.monitoring
├── docker-compose.yml
├── golang
│   └── Dockerfile
├── grafana
│   ├── dashboards
│   └── provisioning
├── notifyHook <- api http://192.168.10.100:8001/notify
│   ├── api
│   ├── go.mod
│   ├── jsonWrap
│   ├── line
│   ├── main.go
│   ├── test
│   └── tools
├── promApp <- alertname: goApp_flip, goApp_processed_ops_total
│   ├── go.mod
│   ├── go.sum
│   └── main.go
└── prometheus
    ├── alert.rules <- add
    └── prometheus.yml
```

alert.rules

```yaml
groups:
  - name: example
    rules:
      - alert: goApp-down
        expr: up == 0
        for: 1s
        labels:
          severity: critical
        annotations:
          title: Node {{ $labels.instance }} is down
          description: Failed to scrape {{ $labels.job }} on {{ $labels.instance }} for more than 3 minutes. Node seems down.
      - alert: goApp-flip
        expr: goApp_flip > 5
        labels:
          severity: warm
        annotations:
          title: flip {{ $labels.instance }} is more than 5 value is {{ $value }}
          description: system make {{ $labels.instance }} some error for test {{ $labels.job }}
```

config.yml

```yaml
global:
  resolve_timeout: 5s # 未收到標記告警通知，等待 timeout 時間之後事件標記為 resolve。
route:
  receiver: default-receiver
  group_wait: 10s # 告警群組訊息建立後的等待時間
  group_interval: 10s # 上下两组发送告警的间隔时间
  repeat_interval: 20s # 重新發送相同告警訊息的間隔時間
  group_by: [cluster, alertname]
  routes:
    - receiver: "Line-Notfiy" #利用alertname來取路徑至receivers
      group_wait: 1s
      match:
        alertname: goApp-flip
      #continue: true
    - receiver: "Line-Logger"
      group_wait: 1s
      match:
        alertname: goApp-logger
      #continue: true

# 設定通知管道
receivers:
  - name: Line-Notfiy
    webhook_configs:
      - url: "http://192.168.10.100:8001/notify" #notifyHook app
        http_config:
          authorization:
            type: Bearer
            credentials: "XXXX"
        #send_resolved: true
  - name: Line-Logger
    webhook_configs:
      - url: "http://192.168.10.100:8001/logger" #notifyHook app
        http_config:
          authorization:
            type: Bearer
            credentials: "XXXX"
        #send_resolved: true
  - name: default-receiver
    webhook_configs:
      - url: "http://192.168.10.100:8001/XXX"
        http_config:
          authorization:
            type: Bearer
            credentials: "XXXX"
        #send_resolved: true
```

![1674028560278.jpg](Docker%2037983de704874fa6b1dcc661a56e0858/1674028560278.jpg)

用curl測試Alertmanager

```bash
curl --location --request POST 'http://192.168.10.101:9093/api/v2/alerts' \
--header 'Content-Type: application/json' \
--data-raw '[
    {
      "status": "resolve",
      "labels": {
      "alertname": "goApp_api",
      "instance": "localhost:8080",
      "job": "node",
      "severity": "critical"
    },
      "annotations": {
      "summary": "測試"
    },"generatorURL": "http://localhost:9090/graph"
    }
  ]'
```

在alertmanager:9093看到的alert

![1674032752022.jpg](Docker%2037983de704874fa6b1dcc661a56e0858/1674032752022.jpg)

notify上面看到的alert

![1674032768583.jpg](Docker%2037983de704874fa6b1dcc661a56e0858/1674032768583.jpg)

## NotifyHook API

問題: Alertmanager 提供的方式種共有以下這幾種方式

```
# The unique name of the receiver.
name: <string>

# Configurations for several notification integrations.
email_configs:
  [ - <email_config>, ... ]
opsgenie_configs:
  [ - <opsgenie_config>, ... ]
pagerduty_configs:
  [ - <pagerduty_config>, ... ]
pushover_configs:
  [ - <pushover_config>, ... ]
slack_configs:
  [ - <slack_config>, ... ]
sns_configs:
  [ - <sns_config>, ... ]
victorops_configs:
  [ - <victorops_config>, ... ]
webhook_configs:
  [ - <webhook_config>, ... ]
wechat_configs:
  [ - <wechat_config>, ... ]
telegram_configs:
  [ - <telegram_config>, ... ]
webex_configs:
  [ - <webex_config>, ... ]
```

最符合我使用的應該只有webhook

```yaml
receivers:
  - name: Line-Notfiy
    webhook_configs:
      - url: "https://notify-api.line.me/api/notify?message=test" #每次呼叫只能固定的值不太符合使用,要硬幹的話可能會有好幾個route跟receiver
        http_config:
          authorization:
            type: Bearer
            credentials: "XXXX"
        #send_resolved: true
```

所以只能自己寫API來處理

資料夾結構

```
.
├── alertmanager
│   └── config.yml
├── config.monitoring
├── docker-compose.yml
├── grafana
│   ├── dashboards
│   └── provisioning
├── notifyHook #src
│   ├── api
│   ├── go.mod
│   ├── jsonWrap
│   ├── line
│   ├── main.go
│   ├── test
│   └── tools
├── notifyHookD #Dockerfile
│   └── Dockerfile
├── promApp
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── promAppD
│   └── Dockerfile
└── prometheus
    ├── alert.rules
    └── prometheus.yml
```

Dockerfile

```docker
FROM golang:alpine as builder
WORKDIR /app
COPY ./notifyHook .
RUN go mod download
RUN go build -o exec

FROM alpine
WORKDIR /usr/bin
COPY --from=builder /app/exec .
RUN apk update && apk add tzdata
```

docker-compose

```docker
notifyHook:
    build: #假如dockerfile不再根目錄下必須使用這個方法
      context: ./
      dockerfile: ./notifyHookD/Dockerfile
    container_name: notfiyHook
    entrypoint: ./exec
    ports:
      - '8001:8001'
```

用curl測試notifyHook

```docker
curl --location --request POST '192.168.10.101:8001/notify' \
--header 'Authorization: Bearer msMRG9L842SX2HOLTWO9KcykxlgmKZmXpgx8nkSPUhU' \
--header 'Content-Type: application/json' \
--data-raw '{
    "receiver": "api-receiver",
    "status": "firing",
    "alerts": [
        {
            "status": "firing",
            "labels": {
                "alertname": "goApp flip",
                "instance": "golang:9001",
                "job": "goApp",
                "severity": "warm"
            },
            "annotations": {
                "description": "Failed to scrape goApp on golang:9001 for more than 3 minutes. Node seems down.",
                "title": "flip golang:9001 is more 5"
            },
            "startsAt": "2023-01-12T08:09:31.656Z",
            "endsAt": "2023-01-12T08:10:01.656Z",
            "generatorURL": "http://07712047c817:9090/graph?g0.expr=goApp_flip+%3E+5\u0026g0.tab=1",
            "fingerprint": "989b0a76f5d22d7f"
        }
    ],
    "groupLabels": {
        "alertname": "goApp flip"
    },
    "commonLabels": {
        "alertname": "goApp flip",
        "instance": "golang:9001",
        "job": "goApp",
        "severity": "warm"
    },
    "commonAnnotations": {
        "description": "Failed to scrape goApp on golang:9001 for more than 3 minutes. Node seems down.",
        "title": "flip golang:9001 is more 5"
    },
    "externalURL": "http://28fdaf41f3ea:9093",
    "version": "4",
    "groupKey": "{}/{alertname=\"goApp flip\"}:{alertname=\"goApp flip\"}",
    "truncatedAlerts": 0
}'
```

![1674031308712.jpg](Docker%2037983de704874fa6b1dcc661a56e0858/1674031308712.jpg)

![1674031476240.jpg](Docker%2037983de704874fa6b1dcc661a56e0858/1674031476240.jpg)
