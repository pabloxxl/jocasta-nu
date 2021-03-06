FROM golang:latest as build

ARG TARGETOS
ARG TARGETARCH

ENV WORKDIR ${GOPATH}/app
ENV CGO_ENABLED=0

WORKDIR $WORKDIR

COPY go.mod $WORKDIR/
COPY go.sum $WORKDIR/
RUN go mod download
COPY . $WORKDIR

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/dns cmd/dns/main.go
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/rest cmd/rest/main.go

FROM scratch AS dns
COPY --from=build /app/dns /

FROM scratch AS rest
COPY --from=build /app/rest /
