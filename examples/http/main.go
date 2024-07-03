package main

import (
	"context"
	"net/http"

	"github.com/Vivino/go-shezmu"
	shezttp "github.com/Vivino/go-shezmu/http"
	"github.com/julienschmidt/httprouter"
)

var tracker = func(ctx context.Context, name string) (context.Context, func()) {
	return context.Background(), func() {}
}

func main() {
	sv := shezmu.Summon()
	server := shezttp.NewServer(sv, ":2255")
	server.Get("/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Write([]byte("It works!"))
	})
	go server.Start()

	sv.StartDaemons(tracker)
	sv.HandleSignals()
}
