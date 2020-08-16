package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rknizzle/handyman/pkg/producer"
	"net/http"
)

var p *producer.Producer

func main() {
	var err error
	p, err = producer.NewProducer()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/new", produce)

	port := 8080
	fmt.Printf("Starting server on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func produce(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		jsonBody := buf.String()

		res, err := p.Produce(jsonBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// write the task ID to the response
		js, err := json.Marshal(
			struct {
				TaskUUID string `json:"task_uuid"`
			}{
				TaskUUID: res.GetState().TaskUUID,
			})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		// 404 on unsupported HTTP methods
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}
