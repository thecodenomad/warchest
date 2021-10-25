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

**`transaction` object**

The the list of objects in the config should be the initial transactions when your crypto currencies were
purchased:

```
{
    "coin_symbol": "ETH",
    "amount": 10.1,
    "purchased_price": 34.5,
    "transaction_fee": 6.56
}
```

## Building

Assuming `make` and `golang >= 1.16` are installed, then just run `make`

## Executing 

Before you can run the binary, you will need to create a config file that includes all of the transactions
that have been made for all the coins in your wallet. See the below for formatting.

### Example Config
A full example of a config can be seen below:

``` 
{
  "coin_purchases": [
    {
      "coin_symbol": "ETH",
      "amount": 10.1,
      "purchased_price": 100.0,
      "transaction_fee": 6.56,
    },
    {
      "coin_symbol": "ALGO",
      "amount": 5.0,
      "purchased_price": 1.2,
      "transaction_fee": 0.35,
    }
  ]
}
```

Once the config is created, it can be specified at execution time

`WARCHEST_CONFIG=<your config filepath> ./warchest`

Ex.

```
$ ./warchest 
Server enabled: false
Updating crypto wallet
2021/10/25 12:59:35 There are 1 coin(s) in your wallet, calculating...
2021/10/25 12:59:35 There are 9 DOGE transactions in your wallet, calculating...
        Current rate for DOGE: 0.266650
        Initial Cost of DOGE: 0.266650
        Total Amount of DOGE: 1
        Current cost of DOGE: 0.266650
        Total profit for DOGE: 0.0
Current Wallet's Net Profit: 0.0
```
