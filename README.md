# Warchest

The purpose of this app is to retrieve stats based on cryptocurrency holdings to produce a consolidated net profit.

## Requirements

- Read in JSON config file (exposed via Environment Variable "WARCHEST_CONFIG")
- Calculate net profit based on coins held, purchased price, and transaction fees

## Config

The config file needs to be in JSON format, and must at least contain key `purchased_coins` that is
an array of purchased coin objects:

**Empty Config**
```
{
    "purchased_coins": []
}
```

**`purchased_coin` object**

The purchased coin object must have the following fields:

```
{
    "coin_symbol": "ETH",
    "amount": 10.1,
    "purchased_price": 34.5,
    "transaction_fee": 6.56,
    "purchase_exchange_rate": 0.001
}
```

### Example Config
A full example of a config can be seen below:

``` 
{
  "purchased_coins": [
    {
      "coin_symbol": "ETH",
      "amount": 10.1,
      "purchased_price": 34.5,
      "transaction_fee": 6.56,
      "purchase_exchange_rate": 0.001
    },
    {
      "coin_symbol": "ALGO",
      "amount": 5.0,
      "purchased_price": 2.5,
      "transaction_fee": 0.35,
      "purchase_exchange_rate": 0.40
    }
  ]
}
```

## Make Targets

Current make targets:

#### `make test`
- Will run all the unittests and any integration tests

#### `make build`
- Will build the warchest binary

#### `make deploy`
- Will deploy warchest in a docker container (WIP - not yet available)

#### TBD