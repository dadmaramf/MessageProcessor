FROM ubuntu:latest

WORKDIR /usr/local/bin/app/

COPY build/ /usr/local/bin/app/

COPY config/ /etc/app/config/

COPY migrations/ /usr/local/bin/app/migrations/

COPY entrypoint.sh /usr/local/bin/app/entrypoint.sh

ENV PATHCONF=/etc/app/config/config.yaml

RUN chmod +x /usr/local/bin/app/*

EXPOSE 8081

# Command to run your application
ENTRYPOINT ["/usr/local/bin/app/entrypoint.sh"]