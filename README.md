# TouchID Go

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/aethiopicuschan/touchid-go.svg)](https://pkg.go.dev/github.com/aethiopicuschan/touchid-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/aethiopicuschan/touchid-go)](https://goreportcard.com/report/github.com/aethiopicuschan/touchid-go)
[![CI](https://github.com/aethiopicuschan/touchid-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/aethiopicuschan/touchid-go/actions/workflows/ci.yaml)

`touchid-go` is a Go library for Touch ID authentication on macOS. It provides a simple API and supports `context.Context` for managing timeouts and cancellations.

It is possible that it works on iOS as well, but it has not been confirmed.

## Installation

```sh
go get -u github.com/aethiopicuschan/touchid-go
```

## Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aethiopicuschan/touchid-go"
)

func main() {
	// Wait for 5 seconds for the user to authenticate
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Authenticate with Touch ID
	ok, err := touchid.Authenticate(ctx, "Sample Text")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Authentication timed out")
		} else {
			fmt.Println("Error:", err.Error())
		}
		os.Exit(1)
	}

	// Check the result
	if ok {
		fmt.Println("Authentication succeeded")
	} else {
		fmt.Println("Authentication failed")
		os.Exit(1)
	}
}
```
