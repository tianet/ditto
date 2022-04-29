package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"

	"github.com/tianet/ditto/pkg/application/schema"
)

type RestServer struct {
	Port      int
	Router    *mux.Router
	Endpoints []*Endpoint
}
type Endpoint struct {
	Fields *[]schema.Field
}

func NewRestServer(port int) (RestServer, error) {
	server := RestServer{
		Port:      port,
		Router:    mux.NewRouter(),
		Endpoints: []*Endpoint{},
	}
	fmt.Println(port)
	server.Router.HandleFunc("/", DefaultHandler)
	return server, nil
}

func (r RestServer) ListenAndServe() {
	printEndpoints(r.Router)
	http.ListenAndServe(fmt.Sprintf(":%d", r.Port), r.Router)
}

func (r RestServer) AddEndpoint(endpoint string, fields *[]schema.Field) {
	fmt.Printf("Adding endpoint %s\n", endpoint)
	r.Router.HandleFunc(endpoint, Endpoint{Fields: fields}.Handler)
}

func (s Endpoint) Handler(w http.ResponseWriter, req *http.Request) {
	message, err := schema.GenerateMessage(s.Fields)

	output, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(output))
}

func GetEndpoint(path string) string {
	return "/" + strings.Split(filepath.Base(path), ".")[0]
}

func DefaultHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Ditto is alive")
}

func printEndpoints(r *mux.Router) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		fmt.Printf("%v %s\n", methods, path)
		return nil
	})
}
