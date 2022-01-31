FROM golang:1.18beta1-alpine as build

RUN apk add --no-cache make
ADD https://raw.githubusercontent.com/eficode/wait-for/v2.2.2/wait-for /usr/local/bin
RUN chmod +x /usr/local/bin/wait-for

WORKDIR ./src

COPY . ./

RUN make build && cp ./bin/news-api /usr/local/bin/ && rm -rf /go/src

FROM alpine

COPY --from=build /usr/local/bin/ /usr/local/bin/
