FROM xhumiq/gorocksdb
WORKDIR $GOPATH/src/github.com/triasteam/go-streamnet
COPY ./ $GOPATH/src/github.com/triasteam/go-streamnet
RUN go build -o build/gsn
WORKDIR $GOPATH/src/github.com/triasteam/go-streamnet/build
ENTRYPOINT ["./gsn"]
