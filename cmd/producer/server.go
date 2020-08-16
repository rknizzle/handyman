package main

import (
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
		// send task message to the queue to be consumed
		// TODO: pass request body as the task message
		res, err := p.Produce("")
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
