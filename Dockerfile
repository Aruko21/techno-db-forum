FROM golang:1.13.4-stretch AS build

WORKDIR /usr/src/techno-db

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

FROM ubuntu:18.04 AS release

MAINTAINER Aleksandr Kosenkov

#
# Установка postgresql
#
ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER

USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER technopark WITH SUPERUSER PASSWORD 'park';" &&\
    createdb -O technopark db-forum &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "include_dir='conf.d'" >> /etc/postgresql/$PGVER/main/postgresql.conf
ADD ./postgresql.conf /etc/postgresql/$PGVER/main/conf.d

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# EXPOSE the server port
EXPOSE 5000

COPY --from=build /usr/src/techno-db/ .

# Launch PostgreSQL and server
CMD service postgresql start && ./db-forum-kosenkov --scheme=http --port=5000 --host=0.0.0.0