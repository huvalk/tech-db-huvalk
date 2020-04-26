FROM golang:1.11-stretch AS build

ADD ./ /opt/build/tech-db-huvalk/
WORKDIR /opt/build/tech-db-huvalk/
RUN go get github.com/gobuffalo
RUN go get github.com/mailru/easyjson
RUN go build main.go

FROM ubuntu:18.04 AS release

ENV PGVER 10
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres
ENV PGPASSWORD '12345'


WORKDIR /opt/build/tech-db-huvalk/
COPY --from=build /opt/build/tech-db-huvalk/ ./
RUN ls -al
RUN /etc/init.d/postgresql start &&\
    psql --command "ALTER USER postgres WITH SUPERUSER PASSWORD '12345';" &&\
    createdb -E utf8 -T template0 -O postgres base  &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5\n\
host all all 0.0.0.0/0  md5\
" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN /etc/init.d/postgresql start &&\
    psql base postgres -h 127.0.0.1 -f api/sql/base_V4.sql &&\
    psql base postgres -h 127.0.0.1 -f api/sql/indexes.sql &&\
    /etc/init.d/postgresql stop
RUN echo "listen_addresses='*'\n\
synchronous_commit = off\n\
fsync = off\n\
shared_buffers = 256MB\n\
effective_cache_size = 512MB\n\
full_page_writes = off\n\
fsync = off " >> /etc/postgresql/$PGVER/main/postgresql.conf
#wal_compression = on
EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

EXPOSE 5000

ENV DB_USER 'postgres'
ENV DB_PASSWORD '12345'
ENV DB_NAME 'base'
ENV DB_HOST '/var/run/postgresql/'

CMD service postgresql start && ./main
