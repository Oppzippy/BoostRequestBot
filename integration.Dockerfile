FROM golang:1.16

WORKDIR /go/src/BoostRequestBot
COPY . .
RUN chmod +x ./docker-integration/wait-for-it.sh ./docker-integration/test.sh
RUN go get -d -v ./...
RUN go install -v ./...


CMD ./docker-integration/wait-for-it.sh "$DB_HOST:$DB_PORT" -- ./docker-integration/test.sh
