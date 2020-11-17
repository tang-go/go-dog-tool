FROM busybox:latest
RUN mkdir /app
ADD Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
ADD go-dog-auth /app/go-dog-auth
ADD config /app/config
RUN chmod -R 777 /app/go-dog-auth
WORKDIR /app
CMD ["/app/go-dog-auth"]