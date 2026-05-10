# Ticket Order Simulation
## Requirement
- Go 1.26 version
- make
- docker

## Migrate & Run Application Backend
```sh
cd ticketing
docker compose up
make migration
make run
```

## Create Data Ticket
```curl
postman request POST 'localhost:8001/v1/ticket' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Basic dXNlcjpwYXNz' \
  --body '{
    "title": "Special Show Anak Lamakera",
    "quota": 10,
    "price": 250000
}' \
  --auth-basic-username 'user' \
  --auth-basic-password 'pass'
```

## Run Script War Ticket
```sh
cd script-war
go run main.go
```

## Open Dashboard Worker
([Dashboard Queue Worker](http://localhost:8080))

## Postman
([See Postman](https://github.com/willy182/ticket-order-simulation/blob/main/ticketing/docs/Ticketing.postman_collection.json))