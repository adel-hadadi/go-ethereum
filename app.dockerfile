FROM alpine:latest

WORKDIR /app

COPY webApp .

CMD [ "./webApp" ]