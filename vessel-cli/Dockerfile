FROM alpine:latest

RUN mkdir -p /app  
WORKDIR /app

ADD vessel.json /app/vessel.json  
ADD vessel-cli /app/vessel-cli

CMD ["./vessel-cli"]