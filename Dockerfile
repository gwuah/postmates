FROM golang:1.14-alpine AS main-env
# install gcc for uber/h3-go. see https://github.com/uber/h3/issues/354
RUN apk add build-base


RUN mkdir /app
ARG PORT=8080
ENV PORT=${PORT}
ADD . /app/


WORKDIR /app
RUN cd /app
# Attempt to cache the module retrieval
RUN go mod download
RUN go build -o postmates-app

FROM alpine

WORKDIR /app
COPY --from=main-env /app/postmates-app /app
COPY --from=main-env /app/database /app/database
COPY .env /app/.env

# Add docker-compose-wait tool -------------------
ENV WAIT_VERSION 2.7.3
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait
EXPOSE $PORT

CMD ["postmates-app"]