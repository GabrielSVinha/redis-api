FROM ubuntu

COPY redis-api /

CMD [ "/redis-api" ]
