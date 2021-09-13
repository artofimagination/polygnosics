FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/dummy-resourcedb
ARG SERVER_PORT

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++ git curl lsof
RUN go mod tidy

RUN cd $GOPATH/src/dummy-resourcedb/ && go build main.go

# This application is exposed through SERVER_PORT to the outside
# See .env to change the value.
EXPOSE $SERVER_PORT

# Run the executable
CMD ["./main"]