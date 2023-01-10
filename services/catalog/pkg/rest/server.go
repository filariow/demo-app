package rest

import (
	"context"
	"encoding/json"
	"eshop-catalog/pkg/persistence"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const BasePath = "/products/"

type HttpServer struct {
	Mux *http.ServeMux
	r   persistence.Repository
}

func NewHttpServer(r persistence.Repository) *HttpServer {
	m := http.NewServeMux()

	s := &HttpServer{
		Mux: m,
		r:   r,
	}

	m.HandleFunc(BasePath, s.products)

	return s
}

func (s *HttpServer) products(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	if r.URL.Path == BasePath {
		s.listProducts(ctx, w)
	} else {
		p := strings.TrimLeft(r.URL.Path, BasePath)
		log.Printf("get product: %s, %s", r.URL.Path, p)
		s.getProduct(ctx, w, p)
	}

	return
}

func (s *HttpServer) encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (s *HttpServer) listProducts(ctx context.Context, w http.ResponseWriter) {
	pp, err := s.r.List(ctx)
	if err != nil {
		log.Printf("error retrieving products: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d, err := json.Marshal(pp)
	if err != nil {
		log.Printf("error unmarshaling products: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(d); err != nil {
		log.Printf("error writing response: %v", err)
		return
	}
}

func (s *HttpServer) getProduct(ctx context.Context, w http.ResponseWriter, id string) {
	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id: it must be a valid uuid"))
		return
	}

	p, err := s.r.Read(ctx, id)
	if err != nil {
		log.Printf("error retrieving product with id '%s': %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	d, err := json.Marshal(p)
	if err != nil {
		log.Printf("error unmarshaling product: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(d); err != nil {
		log.Printf("error writing response: %v", err)
		return
	}
}
