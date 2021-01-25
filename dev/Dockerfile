FROM golang:1.15

RUN go get github.com/githubnemo/CompileDaemon

RUN mkdir -p /go/src/github.com/artofimagination/polygnosics
WORKDIR /go/src/github.com/artofimagination/polygnosics

VOLUME /go/src/github.com/artofimagination/polygnosics

COPY ./entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["run"]

EXPOSE 80