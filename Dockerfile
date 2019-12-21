FROM golang:1.12-alpine as builder
RUN apk add git
COPY . /go/src/shuSemester
ENV GO111MODULE on
WORKDIR /go/src/shuSemester/cli
RUN go get && go build
WORKDIR /go/src/shuSemester/web
RUN go get && go build

FROM alpine
MAINTAINER longfangsong@icloud.com
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
COPY --from=builder /go/src/shuSemester/web/web /
COPY --from=builder /go/src/shuSemester/cli/cli /
WORKDIR /
CMD ./web
ENV PORT 8000
EXPOSE 8000