FROM golang:1.26-alpine AS base
ENV CGO_ENABLED=0
WORKDIR /employees
COPY go.* ./
RUN go mod download && go mod verify
COPY . .

FROM base AS build-stage
RUN go build -o main ./cmd/employee

FROM base AS build-verification
RUN go build -o main ./cmd/verification

FROM alpine:3.21 AS prod
WORKDIR /employees
# the zap logger tees to ./log/out.log and panics if the path is not writable
RUN mkdir -p /employees/log
COPY --from=build-stage /employees/main main
EXPOSE 8080
CMD ["./main"]

# verification-service: standalone downstream hop called by employee-service,
# see cmd/verification. Kept in this Dockerfile since it shares the module.
FROM alpine:3.21 AS verification
WORKDIR /verification
COPY --from=build-verification /employees/main main
EXPOSE 8090
CMD ["./main"]
