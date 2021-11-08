# Warchest

The purpose of this app is to retrieve stats based on cryptocurrency holdings to produce a consolidated net profit.

## Design Requirements

- Read in JSON config file (exposed via Environment Variable "WARCHEST_CONFIG")
- Calculate net profit based on coins held, purchased price, and transaction fees

## Execution and Build requirements

This Readme assumes `make`, `docker` and `golang >= 1.16` are installed. Strictly speaking only `golang >= 1.16` is 
required, but that isn't very production-like is it?

### Building

`make`

## Viewing Your Warchest

The following environmental variables are available for execution:  

* CB_API_KEY=`<your api key>` 
* CB_API_SECRET=`<api keys dirty little secret>`
* WARCHEST_CONFIG=`<path to your warchest transaction config>` -- WIP

When the api key and api secret are set, warchest will query for all of the coins available in the wallet associated
with the api key, and then proceed to calculate the total net profit for the supported keys (currently only DOGE and 
SHIB).

> NOTE: `WARCHEST_CONFIG` is meant to skip the querying of available coins' transactions.
> This is still a WIP and doesn't do anything helpful for execution (only useful for dev).

There is also a Makefile target to help make execution easier:

`make run`

Your service will be available at http://localhost:8080/

## Demo mode

If `CB_API_KEY=demo` when executing the binary, the command line utility will return the calculations provided by
the config `src/config/testdata/CoinConfig.json`. 

If the server flag has been passed in, then the server will start in Demo mode providing a UI representation of the
demo config.

Additionally, there is a Makefile target:

`make demo`

# Testing

## Unittests
Go Gadget Go! Go test!

`make test`

## Integration and Beyond tests
Placeholder buckets that have a few tests enabling quicker test building.

### L2
L2 testing for this project includes the verification of API responses.

`make L2`

### L3 (WIP)
L3 testing for this project will include a higher level test using cypress or selenium

`make L3`

## Configuration -- W.I.P.

Configs are still a work in progress, ideally the coins and their transactions should be saved each
time they are pulled from the API. 

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

(The intention was to add this to your path as the name hasn't changed yet, but these are effectively
the transactions. Command line options for handling this better are in the works.)

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
