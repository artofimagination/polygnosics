FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/dummy-resourcedb

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++ git curl lsof
RUN go mod tidy

RUN cd $GOPATH/src/dummy-resourcedb/ && go build main.go

# This container exposes port 8183 to the outside world
EXPOSE 8182

# Run the executable
CMD ["./main"]