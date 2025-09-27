FROM golang:1.25-alpine AS builder

EXPOSE 8080 6969

ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o geomatch-api ./cmd


FROM gcr.io/distroless/static
USER nonroot:nonroot

WORKDIR /app

COPY --from=builder /app/geomatch-api .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/resources ./resources

CMD ["./geomatch-api"]
