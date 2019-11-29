FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git gcc
WORKDIR /go/src/github.com/PaulRaUnite/uni-db-index-task/shop-server
COPY ./shop-server/ .
COPY ./vendor ./vendor
RUN go build -o /usr/local/bin/shop-server github.com/PaulRaUnite/uni-db-index-task/shop-server/cmd/shop-server

# final stage
FROM alpine
RUN apk add --no-cache bash
WORKDIR /service
COPY --from=build-env /usr/local/bin/shop-server /usr/local/bin/shop-server
COPY ./shop-server/config.yaml /service/config.yaml
COPY ./shop-server/run.sh /service/run.sh
RUN ["chmod", "+x", "/service/run.sh"]
ENTRYPOINT /service/run.sh