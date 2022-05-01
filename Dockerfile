FROM golang:1.18-alpine AS build

WORKDIR /boostrequestbot

COPY . .
RUN CGO_ENABLED=0 go build -o BoostRequestBot
RUN chmod +x BoostRequestBot

FROM gcr.io/distroless/static-debian11

WORKDIR /
COPY --from=build /boostrequestbot/BoostRequestBot .

EXPOSE 80/tcp

USER nonroot:nonroot

ENTRYPOINT [ "/BoostRequestBot" ]
