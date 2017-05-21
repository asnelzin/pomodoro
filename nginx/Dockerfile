FROM nginx:stable-alpine

ARG TZ=Europe/Moscow

ADD conf/nginx.conf /etc/nginx/nginx.conf
ADD conf/pomodoro-nginx.conf /etc/nginx/conf.d/pomodoro.conf
ADD conf/robots.txt /robots.txt

ADD scripts/entrypoint.sh /entrypoint.sh
ADD scripts/le.sh /le.sh

RUN \
 chmod +x /entrypoint.sh && \
 chmod +x /le.sh && \
 apk add  --update letsencrypt tzdata openssl && \
 cp /usr/share/zoneinfo/$TZ /etc/localtime && \
 echo $TZ > /etc/timezone && \
 rm -rf /var/cache/apk/*

CMD ["/entrypoint.sh"]