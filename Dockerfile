FROM golang:1.11-stretch AS build

ADD ./ /opt/build/tech-db-huvalk/
WORKDIR /opt/build/tech-db-huvalk/
RUN go mod tidy && go build main.go

FROM ubuntu:18.04 AS release

ENV PGVER 12
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get -y update && apt-get -y install gnupg && apt-get -y install wget
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" | tee  /etc/apt/sources.list.d/pgdg.list
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
    psql base postgres -h 127.0.0.1 -f api/sql/base_V5.sql &&\
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

CMD service postgresql start &&\
     ./main
