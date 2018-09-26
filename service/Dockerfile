FROM golang
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get -d -v && go build -o main .
CMD ["/app/main"]
