FROM mysql:8.0

ENV MYSQL_ROOT_PASSWORD=root123456 \
    MYSQL_DATABASE=weixin \
    MYSQL_USER=weixin \
    MYSQL_PASSWORD=123456

COPY my.cnf /etc/mysql/conf.d/

RUN echo "character-set-server=utf8mb4" >> /etc/mysql/conf.d/docker.cnf \
    && echo "collation-server=utf8mb4_unicode_ci" >> /etc/mysql/conf.d/docker.cnf

HEALTHCHECK --interval=10s --timeout=5s --retries=5 \
  CMD mysqladmin ping -h localhost -u weixin -p123456 