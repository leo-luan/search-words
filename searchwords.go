// sarch words in a dict
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// wg is used to prevent the main program to finish before the goroutines' result
var wg sync.WaitGroup

const (
	apiURL = "http://dicionario-aberto.net/search-json/"
)

// WordDefinion is used to receive the JSON unmarshalling
type WordDefinion struct {
	Entry struct {
		Sense []struct {
			Def string `json:"def"`
		} `json:"sense"`
	} `json:"entry"`
}

func searchDefinion(word string) {

	// create the request
	req, err := http.NewRequest("GET", apiURL+word, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// client is used to make the request
	client := &http.Client{}

	// make the request waiting for resp Response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	defer resp.Body.Close()

	var record WordDefinion

	// unmarshalling the response into the record struct
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	// is the Sense slice is empty, there is no definition for this word
	if len(record.Entry.Sense) == 0 {
		fmt.Println(word + ": Sem definição.")
	} else {
		def := record.Entry.Sense[0].Def
		def = strings.Split(def, ".<br/>")[0]
		def = strings.Split(def, ":")[0]
		fmt.Println("")
		fmt.Println(word+":", def)
	}
	fmt.Println("")
	wg.Done()
}

func main() {
	// use the OS args to search
	words := os.Args[1:]
	// add goroutines to the WaitGroup
	wg.Add(len(words))
	for _, word := range words {
		go searchDefinion(word)
	}
	// wait for all goroutines
	wg.Wait()
}
