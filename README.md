# api

athena core api

## SETUP

1. git clone the repo and make a copy of .env.sample as .env & update the env vars.

```bash
git clone https://github.com/gwuah/api.git
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

# Delivery State Types

- pending - Delivery has been accepted but does not yet have a courier assigned
- pending_pickup - Courier is assigned and is en route to pick up the items
- nearing_pickup - Courier is closing in on pickup point
- at_pickup - Courier is at the pickup location
- delivery_ongoing - Courier has picked up order is moving towards the dropoff
- nearing_dropoff - Courier is closing in on the dropoff point
- at_dropoff - Courier is at the dropoff location
- delivered - Courier has completed delivery
- canceled - Delivery has been canceled.
