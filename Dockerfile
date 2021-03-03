FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/polygnosics

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++ git curl lsof
RUN git clone https://github.com/artofimagination/mysql-user-db-go-interface /tmp/mysql-user-db-go-interface && \
  cp -r /tmp/mysql-user-db-go-interface/db $GOPATH/src/polygnosics && \
  rm -fr /tmp/mysql-user-db-go-interface
RUN go mod tidy

RUN cd $GOPATH/src/polygnosics/ && go build main.go
RUN chmod 0766 $GOPATH/src/polygnosics/scripts/init.sh

# This container exposes port 8081 to the outside world
EXPOSE 8081

# Run the executable
CMD ["./scripts/init.sh"]