# Warchest

The purpose of this app is to retrieve stats based on cryptocurrency holdings to produce a consolidated net profit.

## Requirements

- Read in JSON config file (exposed via Environment Variable "WARCHEST_CONFIG")
- Calculate net profit based on coins held, purchased price, and transaction fees

## Building

Assuming `make` and `golang >= 1.16` are installed, then just run:

`make`

## Contained Execution 

To spin up a docker container that includes a running version of the service, run:

`make run`

Your service will be available at http://localhost:8080/

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

# WIP

## Configuration

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

(The intention was to add this to your path as the name hasn't changd yet, but these are effectively
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
