package main

import (
	"net/http"

	"github.com/Vivino/go-shezmu"
	shezttp "github.com/Vivino/go-shezmu/http"
	"github.com/julienschmidt/httprouter"
)

func main() {
	sv := shezmu.Summon()
	server := shezttp.NewServer(sv, ":2255")
	server.Get("/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Write([]byte("It works!"))
	})
	go server.Start()

	sv.StartDaemons()
	sv.HandleSignals()
}
