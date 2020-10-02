FROM golang:1.14-alpine

WORKDIR $GOPATH/src/aiplayground

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++
RUN go mod tidy
RUN cd $GOPATH/src/aiplayground/ && go build main.go

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]