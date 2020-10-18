FROM golang:1.14-alpine

WORKDIR $GOPATH/src/polygnosics

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++
RUN go mod tidy
RUN cd $GOPATH/src/polygnosics/ && go build main.go

# This container exposes port 8081 to the outside world
EXPOSE 8081

# Run the executable
CMD ["./main"]