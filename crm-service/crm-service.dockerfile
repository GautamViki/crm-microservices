FROM alpine:latest

RUN mkdir /app

COPY crmApp /app

CMD [ "/app/crmApp" ]