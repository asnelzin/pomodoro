FROM golang:alpine

# build service
COPY app /go/src/github.com/asnelzin/pomodoro/app
RUN \
 apk add --update --no-progress git && \
 cd /go/src/github.com/asnelzin/pomodoro/app && \
 go get -v && \
 go build -ldflags "-X main.revision=$(date +%Y%m%d-%H%M%S)" -o /srv/pomodoro && \
 apk del git && rm -rf /go/src/* && rm -rf /var/cache/apk/*

RUN \
 echo "#!/bin/sh" > /srv/exec.sh && \
 echo "tail -F /srv/logs/pomodoro.log &" /srv/exec.sh && \
 echo "/srv/pomodoro >> /srv/logs/pomodoro.log" >> /srv/exec.sh && \
 chmod +x /srv/exec.sh


EXPOSE 8080
#USER pomodoro
WORKDIR /srv
VOLUME ["/srv/logs"]
ENTRYPOINT ["/srv/exec.sh"]