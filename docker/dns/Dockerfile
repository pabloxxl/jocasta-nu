FROM golang:latest

ENV GOPATH=/dns
ENV PORT=8090
ENV RESOLVER_IP=1.1.1.1
ENV RESOLVER_PORT=53
EXPOSE 8090/udp
COPY cmd/dns/main.go /dns/
COPY pkg/dns /dns/src/github.com/pabloxxl/jocasta-nu/pkg/dns
RUN go get golang.org/x/net/dns/dnsmessage
RUN go build -o /dns/dns /dns/main.go
RUN rm /dns/main.go
CMD ["/dns/dns"]
