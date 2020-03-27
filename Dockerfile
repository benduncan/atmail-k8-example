FROM golang:1.11-alpine AS build-env

# Build phase
RUN apk add build-base git

ENV API_PORT 8001
ENV GOPATH /workspace/atmail-rbl/
ENV GOBIN /workspace/atmail-rbl/bin

ADD ./ $GOPATH/
WORKDIR $GOPATH/

RUN make clean
RUN make build

RUN apk del build-base git

# Next, just copy the golang binary, create a lightweight environment

FROM alpine
WORKDIR /workspace/atmail-rbl/
RUN apk add ca-certificates

COPY --from=build-env /workspace/atmail-rbl/bin/ /workspace/atmail-rbl/bin/

EXPOSE $API_PORT
ENTRYPOINT ["/workspace/atmail-rbl/bin/atmail-rbl"]