package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

func main() {
	log.Println("Connecting to google storage...")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	bucketName := "gcp-demo-54321"

	bkt := client.Bucket(bucketName)
	attrs, err := bkt.Attrs(ctx)
	if err != nil {
		// TODO: Handle error.
		panic(err)
	}
	fmt.Printf("bucket %s, created at %s, is located in %s with storage class %s\n",
		attrs.Name, attrs.Created, attrs.Location, attrs.StorageClass)

	// Read a list of files

	query := &storage.Query{Prefix: ""}
	it := bkt.Objects(ctx, query)
	fileName := ""
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			panic(err)
		}
		log.Printf("Got %#v\n", obj)
		fileName = obj.Name
	}

	rc, err := bkt.Object(fileName).NewReader(ctx)
	if err != nil {
		panic(err)
	}

	defer rc.Close()
	// imgByte, err := ioutil.ReadAll(rc)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Printf("got %#v\n", jumbotron)

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(file, rc)
	if err != nil {
		panic(err)
	}
	file.Close()

}
