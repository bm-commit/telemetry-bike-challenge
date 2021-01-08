FROM golang:1.15.3-alpine as builder
LABEL maintainer="Brian Marin"

RUN apk update && apk add --virtual build-dependencies build-base ca-certificates

WORKDIR /build
COPY . .

RUN go mod download

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ENV PORT ""
ENV PUBLIC_DIR ""
ENV SIMFILE_DIR ""

RUN go build -ldflags="-w -s" -o app cmd/app/http/main.go

FROM scratch

WORKDIR /opt/app

COPY --from=builder /build/app .
COPY --from=builder /build/data/simfile.json data/
COPY --from=builder /build/public/static/frontend.html public/static/

EXPOSE 3000

ENTRYPOINT [ "./app" ]