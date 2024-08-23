package main

import (
	"cloud_distributed_storage/store/ceph"
	"fmt"
)

func main() {
	bucket := ceph.GetCephBucket("userfile")

	//data, _ := bucket.Get("/ceph/17c05aeaf81d0a676ea2951c03799fd65d23f978")
	//tmpFile, _ := os.Create("/tmp/test_file")
	//tmpFile.Write(data)
	//return

	//err := bucket.PutBucket(s3.PublicRead)
	//fmt.Printf("create bucket err: %v\n", err)

	res, _ := bucket.List("", "", "", 100)
	fmt.Printf("object keys: %v\n", res)

	//err = bucket.Put("/testupload/a.txt", []byte("just for test"), "octet-stream", s3.PublicRead)
	//fmt.Printf("upload err: %v\n", err)
	//
	//res, err = bucket.List("", "", "", 100)
	//fmt.Printf("object keys: %v\n", res)

}
