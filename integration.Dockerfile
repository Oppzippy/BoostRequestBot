FROM golang:1.17

WORKDIR /go/src/BoostRequestBot
COPY . .
RUN chmod +x ./docker-integration/wait-for-it.sh ./docker-integration/test.sh
RUN go get -d -v ./...
RUN go install -v ./...


CMD ./docker-integration/wait-for-it.sh "$DB_ADDRESS" -- ./docker-integration/test.sh
