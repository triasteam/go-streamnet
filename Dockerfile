FROM gorocksdb:1.13
WORKDIR $GOPATH/src/github.com/triasteam/go-streamnet
COPY ./ $GOPATH/src/github.com/triasteam/go-streamnet
ENV HOST_IP false
ENTRYPOINT ["sh","/go/src/github.com/triasteam/go-streamnet/entrypoint.sh"]
