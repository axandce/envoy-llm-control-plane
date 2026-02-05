FROM golang:1.22 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/policy-stub .

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/policy-stub /policy-stub
EXPOSE 9002
ENTRYPOINT ["/policy-stub"]
