ARG GO_VERSION=1.23.0

###########
# MODULES #
###########

FROM golang:${GO_VERSION} AS modules

WORKDIR /src

COPY ./go.mod ./go.sum ./

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

ARG LAST_MAIN_COMMIT_HASH
ARG LAST_MAIN_COMMIT_TIME

ENV FLAG="-X ${GLOBAL_VAR_PKG}.CommitTime=${LAST_MAIN_COMMIT_TIME}"
ENV FLAG="$FLAG -X ${GLOBAL_VAR_PKG}.CommitHash=${LAST_MAIN_COMMIT_HASH}"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -installsuffix 'static' \
    -ldflags "-s -w $FLAG" \
    -buildvcs=true \
    -o /app ./cmd/api/*.go

#########
# FINAL #
#########

FROM scratch AS final

COPY --from=builder /etc/passwd /etc/passwd

ARG BINARY_NAME 

COPY ./config/api /config/api

COPY --from=builder /app /app

USER nonroot

CMD ["/app"]
