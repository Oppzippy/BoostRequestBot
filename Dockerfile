FROM golang:1.17.2-alpine AS build

WORKDIR /boostrequestbot

COPY . .
RUN go build -o BoostRequestBot

FROM gcr.io/distroless/static-debian11

WORKDIR /boostrequestbot
COPY --from=build /boostrequestbot/BoostRequestBot .

EXPOSE 80/tcp

USER boostrequestbot:boostrequestbot

ENTRYPOINT [ "/boostrequestbot/BoostRequestBot" ]
