go-shopify
==========

Golang tool for connecting to Shopify's API

## Installation

You need to have Git and Go already installed.
Run this in your terminal

```sh
go get github.com/arduino/go-shopify
```

## Usage

Import it in your Go code:

```go
import (
  "github.com/arduino/go-shopify/shopify"
)
```

## Client Creation

To initialize a client you need the shopify private app password

```go
shop := shopify.NewClient("your-shop-domain", "app-password")
```

See functions used in example/example.go 
