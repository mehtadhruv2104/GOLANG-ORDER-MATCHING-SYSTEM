# GOLANG-ORDER-MATCHING-SYSTEM
This application is matching engine that receives market and limit orders and tries 
to find the best match for each order.
I have implemented an in memory heap that always gives back the best match 
It is one of the most efficient solution.

How to set up the project

-> Install Mysql using Docker

docker run --name new-mysql -e MYSQL_ROOT_PASSWORD=your_password -p 3306:3306 -d mysql:latest

-> Create a new database

start mysql server

docker exec -it new-mysql mysql -u root -p

Run the following command in SQL server

CREATE DATABASE newdb;

-> Clone the Repo

git clone https://github.com/mehtadhruv2104/GOLANG-ORDER-MATCHING-SYSTEM

-> Update .env file with DBURL

DB_URL= yourDBURL

-> Run migrations to initilaize tables

cd migrations

goose mysql $(yourDBURL) up

-> Initialize go dependencies

go mod tidy

-> Run the server

go run main.go


-> Sample Curl Request

curl 'http://localhost:8080/api/orders' \
    --header 'Content-Type: application/json' \
    --data '{
    "symbol": "AMZ",
    "side": "sell",
    "type": "limit",
    "price": 8,
    "quantity": 5,
    "status": "open"
  }'

  curl 'http://localhost:8080/api/orders' \
    --header 'Content-Type: application/json' \
    --data '{
    "symbol": "AMZ",
    "side": "buy",
    "type": "market",
    "price": 15,
    "quantity": 50,
    "status": "open"
  }'