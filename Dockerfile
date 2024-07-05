FROM golang:latest as builder
ADD . /go/src/muzz
WORKDIR /go/src/muzz
RUN go get -d -v muzz
RUN go build -o muzz cmd/muzz/main.go

FROM golang:latest 
# FROM scratch

# security
RUN addgroup --system limited
RUN adduser --system --disabled-password --ingroup limited --home /app appuser

COPY --from=builder /go/src/muzz/muzz /app/muzz 
COPY --from=builder /go/src/muzz/.env /app/.env

# RUN chown appuser /app
# RUN chown appuser /app/*

USER appuser
WORKDIR /app

ENTRYPOINT ["/app/muzz"]
