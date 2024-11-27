FROM golang:1.23.2 AS build
WORKDIR /action
COPY . .
RUN CGO_ENABLED=0 go build -o contributors

FROM alpine:latest 
COPY --from=build /action/contributors /contributors
ENTRYPOINT ["/contributors"]
