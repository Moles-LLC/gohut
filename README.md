# gohut

The Hut for the Gopher

## Installation

```sh
go get github.com/moles-llc/gohut
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/moles-llc/gohut"
)

func main() {
	client := gohut.NewClient()
	
	servers, _, err := client.GetPublicServerList("box", 2, 0)
	if err != nil {
		panic(err)
	}
	
	for _, server := range servers {
		fmt.Println(server.Name)
	}
}
```
