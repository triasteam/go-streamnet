FROM gorocksdb:1.15
WORKDIR $GOPATH/src/github.com/triasteam/go-streamnet
COPY ./ $GOPATH/src/github.com/triasteam/go-streamnet
RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct
ENV HOST_IP false
ENTRYPOINT ["sh","/src/github.com/triasteam/go-streamnet/entrypoint.sh"]
