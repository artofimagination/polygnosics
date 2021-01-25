FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/github.com/artofimagination/polygnosics

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++ git bash
RUN go mod tidy

RUN chmod 0766 $GOPATH/src/github.com/artofimagination/polygnosics/scripts/init.sh $GOPATH/src/github.com/artofimagination/polygnosics/scripts/build.sh
RUN cd $GOPATH/src/github.com/artofimagination/polygnosics/ && ./scripts/build.sh polygnosics

# This container exposes port 8081 to the outside world
EXPOSE 8081

# Run the executable
ENTRYPOINT ["./scripts/init.sh"]
CMD ["polygnosics"]