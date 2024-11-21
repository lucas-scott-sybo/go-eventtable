# Events Table Example

A simple example of saving aggrigate events while still maintaing state as a row.


## Quick Start

download atlas if you have not already 

https://atlasgo.io/


Downlaod sqlc if you haven't already

https://docs.sqlc.dev/

Run

```shell
docker compose up
```

to run the postgres db

in another shell run

```shell
make generate
make migrate
make run
```

create a new user

```shell
curl --location 'http://127.0.0.1:5000/users' \
--header 'Content-Type: application/json' \
--data '{
    "name": "lucasthedistroyeddrasdf",
    "password": "JingleJangle841!"
}'
```

update a user

```shell
curl --location --request PUT 'http://127.0.0.1:5000/users' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "name": "bryan",
    "password": "something new"
}'
```

get all users

```shell
curl --location 'http://127.0.0.1:5000/users'
```

get all events

```shell
curl --location 'http://127.0.0.1:5000/events'
```
