# Electra API

This is a core api for a delivery service that provides, geo-indexing, order-dispatching, proximity-searching, ETA, etc <br/>
We leverage mapbox for features such as distance-matrix to enable us sort hits in ascending order. <br/>
It's built to work as a core service that order services will integrate with. <br/>
Find more documentation [here](https://github.com/electra-systems/core-api/blob/master/API_DOCS.md)

# Status

- [x] geo-indexing
- [x] searching
- [x] ETA
- [x] dispatch 
- [x] order creation
- [x] customer login/signup
- [x] customer ratings
- [x] courier ratings


## Project Setup

1. clone the repo and make a copy of .env.sample as .env & update the env vars.

```bash
git clone https://github.com/elecra-systems/core-api.git
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
