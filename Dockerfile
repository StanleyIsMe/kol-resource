ARG GO_VERSION=1.23.0

###########
# MODULES #
###########

FROM golang:${GO_VERSION} AS modules

WORKDIR /src

COPY ./go.mod ./go.sum ./

# install storage.googleapis.com certs
RUN apt-get update && apt-get install -y ca-certificates openssl

# Get certificate from "storage.googleapis.com"
RUN openssl s_client -showcerts -connect storage.googleapis.com:443 </dev/null 2>/dev/null| openssl x509 -outform PEM >  \
    /usr/local/share/ca-certificates/googleapis.crt

# Update certificates
RUN update-ca-certificates

RUN go mod download

###########
# BUILDER #
###########

FROM golang:${GO_VERSION} AS builder

COPY --from=modules /go/pkg /go/pkg

RUN useradd -u 10001 nonroot

WORKDIR /src

COPY ./ ./
 
ARG GLOBAL_VAR_PKG
ARG GO_GOOS=linux
ARG GO_GOARCH=arm64
ARG LAST_MAIN_COMMIT_HASH
ARG LAST_MAIN_COMMIT_TIME

ENV FLAG="-X ${GLOBAL_VAR_PKG}.CommitTime=${LAST_MAIN_COMMIT_TIME}"
ENV FLAG="$FLAG -X ${GLOBAL_VAR_PKG}.CommitHash=${LAST_MAIN_COMMIT_HASH}"

RUN CGO_ENABLED=0 GOOS=${GO_GOOS} GOARCH=${GO_GOARCH} go build \
    -installsuffix 'static' \
    -ldflags "-s -w $FLAG" \
    -buildvcs=true \
    -o /app ./cmd/api/*.go

#########
# FINAL #
#########

FROM scratch AS final

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG BINARY_NAME 

COPY ./config/api /config/api

COPY --from=builder /app /app

USER nonroot

CMD ["/app"]
