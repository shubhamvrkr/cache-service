FROM golang:1.12 as build
WORKDIR /go/src/cacheservice
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
FROM gcr.io/distroless/base
COPY --from=build /go/bin/cacheservice /
COPY --from=build /go/src/cacheservice/config.yml /
COPY --from=build /go/src/cacheservice/dist /dist/
EXPOSE 8080
CMD ["/cacheservice"]
