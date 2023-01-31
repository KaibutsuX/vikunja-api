# syntax=docker/dockerfile:1
#  ┬─┐┬ ┐o┬  ┬─┐
#  │─││ │││  │ │
#  ┘─┘┘─┘┘┘─┘┘─┘

FROM techknowlogick/xgo:go-1.19.2 AS builder

RUN go install github.com/magefile/mage@latest && \
    mv /go/bin/mage /usr/local/go/bin

WORKDIR /go/src/code.vikunja.io/api
COPY . ./

ARG TARGETOS TARGETARCH TARGETVARIANT

RUN mage build:clean && \
    mage release:xgo "${TARGETOS}/${TARGETARCH}/${TARGETVARIANT}"

#  ┬─┐┬ ┐┌┐┐┌┐┐┬─┐┬─┐
#  │┬┘│ │││││││├─ │┬┘
#  ┘└┘┘─┘┘└┘┘└┘┴─┘┘└┘

# The actual image
# Note: I wanted to use the scratch image here, but unfortunatly the go-sqlite bindings require cgo and
# because of this, the container would not start when I compiled the image without cgo.
FROM alpine:3.16 AS runner
LABEL maintainer="maintainers@vikunja.io"
WORKDIR /app/vikunja
ENTRYPOINT [ "/sbin/tini", "-g", "--", "/entrypoint.sh" ]

ENV VIKUNJA_SERVICE_ROOTPATH=/app/vikunja/
ENV PUID 1000
ENV PGID 1000

RUN apk --update --no-cache add tzdata tini
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod 0755 /entrypoint.sh && mkdir files

COPY --from=builder /build/vikunja-* vikunja
