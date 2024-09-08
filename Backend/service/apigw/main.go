package main

import "cloud_distributed_storage/Backend/service/apigw/route"

func main() {
	r := route.Router()
	err := r.Run(":8081")
	if err != nil {
		return
	}
}
