# syntax=docker/dockerfile:1

FROM golang:1.17-alpine
#RUN apk update && apk add
RUN apk update && apk add curl
WORKDIR /app

COPY go.mod ./
#COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /docker-gs-ping
RUN go install .

EXPOSE 8080

CMD [ "/docker-gs-ping" ]