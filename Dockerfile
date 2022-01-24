FROM golang:1.18beta1-alpine as build

RUN apk add --no-cache make

WORKDIR ./src

COPY . ./

RUN make build && cp ./bin/news-api /usr/local/bin/ && rm -rf /go/src

FROM alpine

COPY --from=build /usr/local/bin/ /usr/local/bin/

ENTRYPOINT ["news-api"]
