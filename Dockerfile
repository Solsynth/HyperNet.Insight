# Building Backend
FROM golang:alpine as insight-server

WORKDIR /source
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs -o /dist ./pkg/main.go

# Runtime
FROM golang:alpine

COPY --from=insight-server /dist /insight/server

EXPOSE 8445

CMD ["/insight/server"]
