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

- pending - Delivery has been accepted but does not yet have a electron assigned
- pending_pickup - Electron is assigned and is en route to pick up the items
- nearing_pickup - Electron is closing in on pickup point
- at_pickup - Electron is at the pickup location
- delivery_ongoing - Electron has picked up order is moving towards the dropoff
- nearing_dropoff - Electron is closing in on the dropoff point
- at_dropoff - Electron is at the dropoff location
- delivered - Electron has completed delivery
- canceled - Delivery has been canceled.
