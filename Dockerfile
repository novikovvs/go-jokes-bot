FROM golang

WORKDIR /app

COPY ./src .
RUN go install

RUN go install github.com/githubnemo/CompileDaemon@latest

ENTRYPOINT exec CompileDaemon -build="go build -o /usr/local/bin/analytic-pusher main.go" -command="analytic-pusher"
