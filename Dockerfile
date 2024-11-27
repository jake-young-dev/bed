FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/

FROM golang:1.22
COPY --from=build /app/main ./
ENTRYPOINT [ "./main" ]