package main

import (
	"filestore-server/service/apigw/route"
)

func main() {
	r := route.Router()
	r.Run(":8080")
}
