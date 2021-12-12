FROM golang:1.17-alpine

ENV LISTEN_ADDRESS=":9145"
ENV METRICS_PATH="/metrics"
ENV MONITOR_DIRECTORY="files/"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build

CMD ./Blaze --fs.monitor-directory ${MONITOR_DIRECTORY} --web.listen-address ${LISTEN_ADDRESS} --web.metrics-path ${METRICS_PATH}