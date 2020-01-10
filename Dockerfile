FROM alpine

RUN apk add --no-cache ca-certificates

COPY command-bot /home/bot/

WORKDIR /home/bot/

CMD ["./command-bot"]
