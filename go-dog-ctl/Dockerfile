FROM golang

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.io"
RUN mkdir /app
ADD Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
ADD go-dog-ctl /app/go-dog-ctl
ADD config /app/config
RUN chmod -R 777 /app/go-dog-ctl
WORKDIR /app
CMD ["/app/go-dog-ctl"]