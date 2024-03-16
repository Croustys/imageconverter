FROM golang:1.21.4-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o imageconverter .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/imageconverter .
COPY --from=build /app/templates ./templates

ENV PORT=8080

EXPOSE $PORT

CMD ["./imageconverter"]
