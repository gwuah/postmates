# Postmates

This is the heart of a delivery service. Features include geo-indexing, order-dispatch, proximity-searching, ETA, trip estimates, etc <br/>
We use google maps for features such as distance-matrix and directions <br/>
Find more documentation [here](https://github.com/gwuah/postmates/blob/master/API_DOCS.md)

# Inbuilt Features

- [x] geo-indexing
- [x] geo-radius search
- [x] ETA
- [x] order creation
- [x] order dispatch
- [x] order acceptance
- [x] order rejection
- [x] customer login/signup
- [x] customer ratings

# Requirements

- Postgres
- Redis
- Uber H3

# Architecture

- We use websocket connections for realtime communications with courier and customers. The ws connections are store in-memory in a concurrency safe manner.
- Couriers are indexed using uber's h3 geo-indexing library and grouped in redis.
- When you perform a radius search(closest couriers), we use h3 to calculate all indices 2 levels at resolution 8, see image below. Then we query our courier index, powered by redis to find all the couriers in those locations. Then, we make a request using their lng/lats and the customer's lng/lat to google maps to get the distance and duration from the customer, then we sort that result, and then dispatch the order to these couriers in order of those closest the origin of the request. (See [image](https://github.com/gwuah/postmates/blob/master/img/radius.png)
- The couriers send location updates every 3 seconds. This allows us to know their locations in almost realtime.
- The dispatch logic gives a courier 5 seconds to accept an order, after which it is sent to the next closest/available courier. If none of the available couriers accept the request, the process starts all over again, till someone finally accepts it.

## Project Setup

1. Clone the repo and make a copy of .env.sample as .env & update the env vars.

```bash
git clone https://github.com/gwuah/postmates.git
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

# Demo

- `cd demo` and run `yarn` to install all required dependencies.
- run `node electrons.js` to initiate 3 couriers instances that are constantly sending location updates every 3 seconds
- run `node customer__delivery_request.js` to instantiate a customer that will create a delivery request.
- pay attention to the logs.
