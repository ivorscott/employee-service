FROM golang:1.26-alpine AS base
ENV CGO_ENABLED=0
WORKDIR /employees
COPY go.* ./
RUN go mod download && go mod verify
COPY . .

FROM base AS build-stage
RUN go build -o main ./cmd/employee

FROM alpine:3.21 AS prod
WORKDIR /employees
# the zap logger tees to ./log/out.log and panics if the path is not writable
RUN mkdir -p /employees/log
COPY --from=build-stage /employees/main main
EXPOSE 8080
CMD ["./main"]
