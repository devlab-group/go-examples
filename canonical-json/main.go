package main

import (
	"encoding/json"
	"github.com/ugorji/go/codec"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	jh codec.JsonHandle
	h  = &jh
)

// Turn on an option to convert data into canonical format
func init() {
	jh.EncodeOptions.Canonical = true
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("canonical-json: request body read error %v\n", err)
		http.Error(w, "error request body processing", http.StatusInternalServerError)
	}

	var decoded interface{}
	json.Unmarshal(body, &decoded)

	// Encode decoded body into canonical format and send back to client
	var encoder *codec.Encoder = codec.NewEncoder(w, h)
	encoder.Encode(decoded)
}

func main() {
	http.HandleFunc("/canonical-json", jsonHandler)
	http.ListenAndServe(":8080", nil)
}
