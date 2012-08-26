package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	nimbusiohttp "nimbus.io/http"
	"strings"
)

const (
	testKey = "test key"
	testBody = "test body"
)

func main() {
	fmt.Println("start")
	var credentials *nimbusiohttp.Credentials
	var err error

	sp := flag.String("credentials", "", "path to credentials file")
	flag.Parse()
	if *sp == "" {
		credentials, err = nimbusiohttp.LoadCredentialsFromDefault()
	} else {
		credentials, err = nimbusiohttp.LoadCredentialsFromPath(*sp)
	}
	if err != nil {
		log.Fatalf("Error loading credentials %s\n", err)
	}

	requester, err := nimbusiohttp.NewRequester(credentials); if err != nil {
		log.Fatalf("NewRequester failed %s\n", err)
	}

	collectionList, err := nimbusiohttp.ListCollections(requester, credentials)
	if err != nil {
		log.Fatalf("Request failed %s\n", err)
	}
	fmt.Printf("starting collection list = %v\n", collectionList)

	collectionName := nimbusiohttp.ReservedCollectionName(credentials.Name, 
		fmt.Sprintf("test-%05d", len(collectionList)))
	collection, err := nimbusiohttp.CreateCollection(requester, credentials, 
		collectionName)
	if err != nil{
		log.Fatalf("CreateCollection failed %s\n", err)
	}
	fmt.Printf("created collection = %v\n", collection)

	archiveBody := strings.NewReader(testBody)
	versionIdentifier, err := nimbusiohttp.Archive(requester, credentials, 
		collectionName, testKey, archiveBody)
	if err != nil{
		log.Fatalf("Archive failed %s\n", err)
	}
	fmt.Printf("archived key '%s' to version %v\n", testKey, versionIdentifier)

	retrieveBody, err := nimbusiohttp.Retrieve(requester, credentials, 
		collectionName, testKey)
	if err != nil{
		log.Fatalf("Retrieve failed %s\n", err)
	}

	retrieveResult, err := ioutil.ReadAll(retrieveBody)
	retrieveBody.Close()
	if err != nil{
		log.Fatalf("read failed %s\n", err)
	}
	fmt.Printf("retrieved key '%s'; matches testBody = %v\n", testKey, 
		string(retrieveResult) == testBody)

	success, err := nimbusiohttp.DeleteCollection(requester, credentials, 
		collectionName)
	if err != nil{
		log.Fatalf("DeleteCollection failed %s\n", err)
	}
	fmt.Printf("deleted collection = %s %v\n", collectionName, success)

	fmt.Println("end")
}
