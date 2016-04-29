Open Source API Client
======================

[godoc](https://godoc.org/github.com/OpenSourceOrg/api/client)

This package is used to connect and query the Open Source API. This uses
the internal data modeling of the API to preform requests over the line.

Example
-------

```go
package main

import (
	"github.com/opensourceorg/api/client"
	"log"
)

func ohshit(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	license, err := client.Get("Apache-2.0")
	ohshit(err)
	log.Printf("%s\n", license.Name)
}
```
