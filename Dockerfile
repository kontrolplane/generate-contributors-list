FROM golang:1.23.2 AS build

WORKDIR /action

COPY . .

RUN CGO_ENABLED=0 go build -o generate-contributors-list

FROM alpine:latest 

COPY --from=build /action/generate-contributors-list /generate-contributors-list

ENTRYPOINT ["/generate-contributors-list"]
