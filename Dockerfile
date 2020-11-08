FROM golang:1.15-alpine

WORKDIR /go/src/app
COPY . .

RUN set -eux; \
    apkArch="$(apk --print-arch)"; \
    case "$apkArch" in \
      x86_64) arch=x86_64 ;; \
      armv7) arch=arm ;; \
    esac; \
    mkdir -p /opt/speedtest \
        && wget -O speedtest.tgz https://bintray.com/ookla/download/download_file?file_path=ookla-speedtest-1.0.0-"$arch"-linux.tgz \
        && tar -C /opt/speedtest -xzf speedtest.tgz \
        && rm speedtest.tgz

ENV PATH /opt/speedtest:$PATH

RUN go install

CMD ["speedtest_exporter"]
