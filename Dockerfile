FROM ubuntu:24.10
FROM node:22 AS node

RUN apt update && apt install gcc g++

COPY --from=node /usr/lib /usr/lib
COPY --from=node /usr/local/lib /usr/local/lib
COPY --from=node /usr/local/include /usr/local/include
COPY --from=node /usr/local/bin /usr/local/bin

RUN node -v

WORKDIR /app

COPY ./bin/linux_amd64/app ./
COPY ./web/dist ./web/dist
COPY .env.prod .env.prod

EXPOSE 9001

CMD ["./app","-mode","prod"]
