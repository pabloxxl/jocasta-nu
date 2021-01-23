FROM golang:latest

EXPOSE 8090/udp
COPY cmd/dns/main.go /dns/
RUN go get golang.org/x/net/dns/dnsmessage
RUN go build -o /dns/dns /dns/main.go
RUN rm /dns/main.go
CMD ["/dns/dns"]
