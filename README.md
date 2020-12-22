# Athena V3

This is the heart of a theoretical delivery service. <br/>
Features include geo-indexing, order-dispatch, proximity-searching, ETA, trip estimates, etc <br/>
We leverage google maps for features such as distance-matrix to enable us sort couriers in ascending order. <br/>
Find more documentation [here](https://github.com/electra-systems/core-api/blob/master/API_DOCS.md)

# Inbuilt Features

- [x] geo-indexing
- [x] geo-radius search
- [x] ETA
- [x] order creation
- [x] order dispatch
- [x] order acceptance
- [x] order order rejection
- [x] customer login/signup
- [x] customer ratings
- [x] courier

# Requirements

- Postgres
- Redis
- Uber H3

## Project Setup

1. clone the repo and make a copy of .env.sample as .env & update the env vars.

```bash
git clone https://github.com/gwuah/athena-v3.git
cp .env.sample .env
```

2. Run the app using either :

```bash
go run main.go
```

```bash
go build main.go
./main
```
