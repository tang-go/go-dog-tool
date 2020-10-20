FROM busybox:latest
RUN mkdir /app
ADD Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
ADD go-dog-find-linux /app/go-dog-find-linux
ADD config /app/config
RUN chmod -R 777 /app/go-dog-find-linux
WORKDIR /app
CMD ["/app/go-dog-find-linux"]