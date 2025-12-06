# Bitaksi TaxiHub
TaxiHub is an example of a simple microservices architecture that manages driver registration. 
The project consists of two main components:
- **driver-service**: Provides driver CRUD and nearby searches on MongoDB.
- **gateway**: Validates requests (JWT), applies rate limits, and routes them to the driver-service.
(The services are configured to run with Docker Compose.)

## Features
- Driver creation (POST /drivers)
- Driver update (PUT /drivers/{id})
- Driver listing (GET /drivers)
- Pagination support
- Locating nearby drivers
- JWT validation via gateway
- Rate limit middleware
- Swagger UI documentation
- Working with MongoDB Atlas connection
  
## Project Structure
  bitaksi_TaxiHub/
│
├── driver-service/
│ ├── dto/
│ ├── handler/
│ ├── model/
│ ├── repository/
│ ├── service/
│ ├── docs/ # swagger outputs
│ ├── main.go
│ └── Dockerfile
│
├── gateway/
│ ├── middleware/
│ ├── main.go
│ └── Dockerfile
│
├── tools/
│ ├── make_token.go # JWT generation tool
│ ├── go.mod
│ └── go.sum
│
├── docker-compose.yml
└── .env.example

## Setup
1) Clone the repo and go to the folder
2) Create the .env file. There is an example file in the project.

## Run
To start all services with Docker Compose:

    docker compose up --build

Running services:
- Gateway: http://localhost:8080
- Driver-service: http://localhost:8081

## Swagger (API Documentation)
  http://localhost:8081/swagger/index.html
  
You can try all endpoints from here.

## JWT Token Generation
JWT is required for all requests to the Gateway. To generate tokens:

    cd tools
  
    go run make_token.go

This command generates a token signed with JWT_SECRET in the .env file.
The token expires in 24 hours. You do not need to generate a new one until it expires.

## API Usage
Authorization Header

For all calls made through the gateway, the header must be as follows:
  
    Authorization: Bearer <TOKEN>

 ### Driver listing (Pagination)

    curl -i "http://localhost:8080/drivers?page=1&pageSize=5" -H "Authorization: Bearer <TOKEN>"

 ### Create driver

    curl -X POST "http://localhost:8080/drivers" -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d '{"firstName":"Ahmet","lastName":"Demir","plate":"34ABC123","taxiType":"luks","carBrand":"Renault","carModel":"Corolla","lat":41.0431,"lon":29.0099}'

 ### Update Driver

    curl -X PUT "http://localhost:8080/drivers/<id>" -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d '{"firstName":"Mehmet","taxiType":"sari","lat":41.02,"lon":29.00}'

 ### List nearby drivers

  curl "http://localhost:8080/drivers/nearby?lat=41.02&lon=29.01&radius=2" -H "Authorization: Bearer <TOKEN>"


## Notes

Unit tests have not yet been added.


