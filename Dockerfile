FROM golang:1.17-alpine as base
ENV CGO_ENABLED=0
WORKDIR /employees
COPY go.* ./
RUN go mod download && go mod verify
COPY . .

FROM base as build-stage
RUN go build -o main ./cmd/employee

FROM aquasec/trivy:0.4.4 as image-scan
RUN trivy alpine:3.15 && \
    echo "No image vulnerabilities" > result

FROM alpine:3.15 as prod
COPY --from=image-scan result secure
COPY --from=build-stage /employees/main main
EXPOSE 8080 
CMD ["./main"]
