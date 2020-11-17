FROM busybox:latest
RUN mkdir /app
ADD Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
ADD go-dog-gw /app/go-dog-gw
ADD config /app/config
RUN chmod -R 777 /app/go-dog-gw
WORKDIR /app
CMD ["/app/go-dog-gw"]