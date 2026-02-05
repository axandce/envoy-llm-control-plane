FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/mock-upstream .

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/mock-upstream /mock-upstream
EXPOSE 9000
ENTRYPOINT ["/mock-upstream"]
