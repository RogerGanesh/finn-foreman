FROM alpine:latest 

RUN mkdir /app 

COPY foreman /app

CMD ["/app/foreman"]