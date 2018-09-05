FROM golang:1.11

# copy project
WORKDIR /go/src/github.com/tusupov/go-exercises
COPY . ./

# run test
RUN go test -v ./...
RUN go test -bench=. -v ./...
RUN rm -rf *_test.go

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bankaccount .

FROM alpine:latest

# copy
WORKDIR /app/
COPY --from=0 /go/src/github.com/tusupov/go-exercises/bankaccount .

CMD ./bankaccount -p $PORT
