package main

import "cloud_distributed_storage/Backend/service/apigw/route"

func main() {
	r := route.Router()
	r.Run(":8081")
}
