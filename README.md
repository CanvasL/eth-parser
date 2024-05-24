# eth-parser

A Golang-based parser for the Ethereum blockchain, offering APIs to subscribe to addresses and retrieve transactions in real-time.

## Start

```sh
go run main.go
```

## API

- Get the current block number

```sh
curl "http://localhost:8080/current_block"
```

- Subscribe to an Ethereum address:

```sh
curl -X POST "http://localhost:8080/subscribe?address=<ethereum_address>"
```

- Get transactions for an subscribed Ethereum address:

```sh
curl "http://localhost:8080/transactions?address=<ethereum_address>"
```
