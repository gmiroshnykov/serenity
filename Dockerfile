FROM golang:1.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /serenity

EXPOSE 8080

CMD [ "/serenity" ]
