# Docs

### Courier can initiate connection @

`/courier/realtime/:id`

### Customer can initiate connection @

`/customer/realtime/:id`

see `handler/ws.go` and `handler/handler.go` for more details with regards to realtime connections.

### Get Delivery Quote

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

### Rate Delivery (Customer)

`/v1/customer-rate-trip`

**method:** POST

**data params:**

```
{
    message: "Good Service",
    deliveryId: 1,
    customerId: 1,
    rating: 5
}
```

**response:**

```
{
    message: "success"
    data: true
}
```

### Rate Delivery (Courier)

`/v1/courier-rate-trip`

**method:** POST

**data params:**

```
{
    message: "Good Service",
    deliveryId: 1,
    courierId: 1
    rating: 5
}
```

**response:**

```
{
    message: "success"
    data: true
}
```
