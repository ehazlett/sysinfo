FROM golang:1.15-alpine as build
RUN apk add -U make
COPY . /src/app
WORKDIR /src/app
RUN make build

FROM alpine:latest
COPY --from=build /src/app/bin/sysinfo /bin/sysinfo
CMD ["/bin/sysinfo"]
