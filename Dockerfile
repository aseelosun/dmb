FROM reg.1cb.kz/dockerhub/library/alpine:3.13

RUN apk add --update --no-cache tzdata ca-certificates && update-ca-certificates
ENV TZ=Asia/Almaty
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 8080

COPY ./bin/binary /binary
COPY config.json /config.json
ENTRYPOINT ["/binary"]
