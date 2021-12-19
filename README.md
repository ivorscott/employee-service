# Employee Service

Test Go Service

## Usage

```
go run ./cmd/employee
```

Two hardcoded users exist. Example:
```
curl -v -X GET http://localhost:8080/employees/0 | jq

{
  "id": 0,
  "firstName": "Adam",
  "lastName": "Smith",
  "job": "Philosopher"
}

curl -v -X GET http://localhost:8080/employees/1 | jq 

{
  "id": 1,
  "firstName": "John",
  "lastName": "Locke",
  "job": "Philosopher"
}

curl -v -X GET http://localhost:8080/employees/2 | jq

{
  "error": "not found"
}

curl -v -X GET http://localhost:8080/employees/abc | jq

{
  "error": "bad request"
}
```

# To Do 

[ ] add swagger

[ ] add elasticsearch

[ ] add interfaces

[ ] add postgres

[ ] add fixtures

[ ] add tests

[ ] add goldenfiles

[ ] add kubernetes

[ ] add auth

[ ] add tracing

[ ] add redis

[ ] add nats jetstream

[ ] add prometheus

[ ] add graphana

[ ] add example requests

[ ] add terraform

[ ] add makefile

[ ] add zappr

[ ] add ci workflow