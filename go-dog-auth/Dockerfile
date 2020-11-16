FROM busybox:latest
RUN mkdir /app
ADD Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
ADD auth /app/auth
ADD config /app/config
RUN chmod -R 777 /app/auth
WORKDIR /app
CMD ["/app/auth"]