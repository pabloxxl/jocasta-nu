FROM golang:latest

ENV GOPATH=/rest
EXPOSE 8080
COPY cmd/rest/main.go /rest/
COPY pkg/rest /rest/src/github.com/pabloxxl/jocasta-nu/pkg/rest
COPY blocked.txt /
RUN go get github.com/gorilla/mux
RUN go build -o /rest/rest /rest/main.go
RUN rm /rest/main.go
CMD ["/rest/rest"]
