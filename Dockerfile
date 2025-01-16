FROM ubuntu:24.10
FROM node:22-slim AS node
FROM php:8.3-cli as php

RUN apt update && apt install gcc g++ libonig-dev libargon2-0 -y

COPY --from=node /usr/lib /usr/lib
COPY --from=node /usr/local/lib /usr/local/lib
COPY --from=node /usr/local/include /usr/local/include
COPY --from=node /usr/local/bin /usr/local/bin

RUN node -v

COPY --from=php /usr/local/bin/php /usr/local/bin/php
COPY --from=php /usr/local/lib/php /usr/local/lib/php

RUN php -v

WORKDIR /app

COPY ./bin/linux_amd64/app ./
COPY ./web/dist ./web/dist
COPY .env.prod .env.prod

EXPOSE 9001

CMD ["./app","-mode","prod"]
