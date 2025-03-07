FROM golang AS build
WORKDIR /proxy
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/proxy ./

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/proxy /bin/proxy
ENTRYPOINT [ "/bin/proxy" ]