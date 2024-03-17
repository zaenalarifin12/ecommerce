package client

import (
	"fmt"
)

func ProductByUUID(id string) string {
	url := fmt.Sprintf("/products/%s", id)
	return url
}

var APIProductClient = NewClient("http://localhost:5001")
