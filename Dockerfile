FROM xhumiq/gorocksdb
WORKDIR $GOPATH/src/github.com/triasteam/go-streamnet
COPY ./ $GOPATH/src/github.com/triasteam/go-streamnet
WORKDIR $GOPATH/src/github.com/triasteam/go-streamnet/build
COPY ./config.yml  config.yml
WORKDIR $GOPATH/src/github.com/triasteam/go-streamnet
RUN go get gopkg.in/yaml.v2 \
    && go build -o build/gsn
WORKDIR $GOPATH/src/github.com/triasteam/go-streamnet/build
ENTRYPOINT ["./gsn"]
