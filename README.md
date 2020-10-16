# Electra Core API

Heartbeat of electra

## SETUP

1. git clone the repo and make a copy of .env.sample as .env & update the env vars.

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

# Docs

Heroku Endpoint
`https://electra-staging.herokuapp.com`

Curr Version
`v1`

## Get Delivery Quote
`/v1/get-delivery-cost`

**method:** POST

**data params:** 


```
{
    origin: {
        latitude: 5.677474538991623,
        longitude: -0.24460022375167725
    },
    destination: {
        latitude: 5.6796946725653745,
        longitude: -0.2447180449962616
    }
}
```

**response:**
```
{
    data: {
        estimate: {
            1: {
                productId: 1,
                price: 5
            }
        },
        distance: 2.3166666666666664,
        duration: 319
    },
    message: "success"
}
```
