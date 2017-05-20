FROM golang:alpine

ARG TZ=Europe/Moscow

# add user and set TZ
RUN \
 apk add --update tzdata && \
 adduser -s /bin/bash -D -u 1001 pomodoro && \
 mkdir -p /srv && chown -R pomodoro:pomodoro /srv && \
 cp /usr/share/zoneinfo/$TZ /etc/localtime && \
 echo $TZ > /etc/timezone && \
 rm -rf /var/cache/apk/*

# build service
COPY app /go/src/github.com/asnelzin/pomodoro/app
RUN \
 apk add --update --no-progress git && \
 cd /go/src/github.com/asnelzin/pomodoro/app && \
 go get -v && \
 go build -ldflags "-X main.revision=$(date +%Y%m%d-%H%M%S)" -o /srv/pomodoro && \
 apk del git && rm -rf /go/src/* && rm -rf /var/cache/apk/*

EXPOSE 8080
USER pomodoro
WORKDIR /srv
ENTRYPOINT ["/srv/pomodoro"]