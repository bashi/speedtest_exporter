# speedtest_exporter

A [prometetheus](https://prometheus.io/) exporter that runs [Ookla's speedtest](https://www.speedtest.net/apps/cli) periodically.

## Setup

```sh
$ docker-compose build
$ docker-compose up -d
```

The exporter exports metrics via `localhost:9300/metrics`.
