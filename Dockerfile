FROM golang:1.12-alpine as builder
RUN apk add git
COPY . /go/src/shuSemester
WORKDIR /go/src/shuSemester/cli
RUN go get && go build
WORKDIR /go/src/shuSemester/web
RUN go get && go build

FROM alpine
MAINTAINER longfangsong@icloud.com
COPY --from=builder /go/src/shuSemester/web/web /
COPY --from=builder /go/src/shuSemester/cli/cli /
WORKDIR /
CMD ./web
RUN export PORT=8000
EXPOSE 8000