# Employee Service

Test Go Service

## Usage

```
go run ./cmd/employee
```

Two hardcoded users exist. Example requests:
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