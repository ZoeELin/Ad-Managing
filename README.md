# Ad-Managing
This is the assignment of 2024 Dcard backend intern.

## Requirements
- [x] Admin API (POST /api/v1/ad)
- [x] Place API (GET /api/v1/ad)
- [ ] Test for API that can handle over 10,000 Requests Per Secoud.
- [ ] Ensure that the total active advertisements in the system (i.e., StartAt < NOW < EndAt) are less than 1000.
- [ ] Limit the number of advertisements created per day to not exceed 3000.

## Structure
```
.
├── 2024 Backend Intern Assignment.pdf
├── Dockerfile
├── README.md
├── controllers
│   ├── ads.go
│   └── post_ads.go
├── database
│   └── db.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
└── models
    └── models.go
```
- The main program is located in `main.go`.


## Dependencies
- Gin is a web framework for Go language that focuses on performance and minimalism. 
    - Gin is designed to be fast and lightweight.
    - It provides a HTTP router that allows you to define routes for handling different HTTP requests (GET, POST, PUT, DELETE, etc.) easily.
- Gorm is an ORM library for the Go programming language that provide CRUD operations on your database. 


## Quick Start

```bash
# Clone the repo
git clone https://github.com/ZoeELin/Ad-Managing.git
cd Ad-Managing

# Start the service
docker-compose up
```
Go to http://localhost:5000/api/v1/ad to POST or GET the data.


## API Detail

### Generate ADs
Sent a POST request to `/api/v1/ad` with the advertisement details in JSON format. 
JSON structure:
- Title
- StartAt
- EndAt
- Conditions
    - Age
    - Gender (M, F)
    - Country (TW, JP, US)
    - Platform (android, ios, web)

Example:
```bash
curl -X POST -H "Content-Type: application/json" \
"http://localhost:5000/api/v1/ad" \
--data '{
"title": "AD test",
"startAt": "2023-04-01T03:00:00.000Z",
"endAt": "2024-04-30T16:00:00.000Z",
"conditions": [
{
"ageStart": 28,
"ageEnd": 45,
"country": ["TW", "JP", US],
"platform": ["ios"]
}
]
}'
```

### List ADs 
Send a GET request to `/api/v1/ad` with the conditions as query parameters to list active advertisements that match specitic conditions.

Parameters:

- offset, limit: used for pagination
- age
- gender
- country
- platform

Example:
```bash
curl -X GET -H "Content-Type: application/json" \
"http://localhost:5000/api/v1/ad?offset=15&limit=5&age=24&gender=M&platform=ios"
```

## Testing


