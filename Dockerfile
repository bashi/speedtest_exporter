FROM golang:1.15-alpine

WORKDIR /go/src/app
COPY . .

# TODO: Check "$(apk --print-arch)" and download an appropriate binary.
RUN mkdir -p /opt/speedtest \
    && wget -O speedtest.tgz https://bintray.com/ookla/download/download_file?file_path=ookla-speedtest-1.0.0-x86_64-linux.tgz \
    && tar -C /opt/speedtest -xzf speedtest.tgz \
    && rm speedtest.tgz

ENV PATH /opt/speedtest:$PATH

RUN go install

CMD ["speedtest_exporter"]
