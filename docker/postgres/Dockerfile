FROM postgres:latest

ADD ./migrations /docker-entrypoint-initdb.d/

EXPOSE 5432

CMD ["postgres"]
