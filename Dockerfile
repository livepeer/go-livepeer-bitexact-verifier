FROM yondon/livepeer-ubuntu1604:latest
MAINTAINER Yondon Fu "yondon@livepeer.org"

ENV PKG_CONFIG_PATH /root/compiled/lib/pkgconfig

ADD . $GOPATH/src/github.com/livepeer/go-livepeer-bitexact-verifier
RUN cd $GOPATH/src/github.com/livepeer/go-livepeer-bitexact-verifier && \
    go build -v . && \
    ln -s $GOPATH/src/github.com/livepeer/go-livepeer-bitexact-verifier/go-livepeer-bitexact-verifier /usr/bin/verifier

ENTRYPOINT verifier $ARG0 $ARG1
