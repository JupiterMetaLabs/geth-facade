FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /rpc-facade

FROM gcr.io/distroless/base-debian12
COPY --from=build /rpc-facade /usr/local/bin/rpc-facade
EXPOSE 8545 8546
ENTRYPOINT ["/usr/local/bin/rpc-facade"]
