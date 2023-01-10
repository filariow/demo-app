package rest

import (
	"context"
	"encoding/json"
	"eshop-orders/pkg/models"
	"eshop-orders/pkg/persistence"
	"io"
	"log"
	"net/http"
	"strings"
)

const BasePath = "/orders/"

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

	m.HandleFunc(BasePath, s.orders)

	return s
}

func (s *HttpServer) orders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	if r.URL.Path == BasePath {
		switch r.Method {
		case http.MethodGet:
			s.listOrders(ctx, w)
		case http.MethodPost:
			s.createOrder(r, w)
		}
	} else {
		p := strings.Replace(r.URL.Path, BasePath, "", 1)
		s.getOrders(ctx, w, p)
	}
}

func (s *HttpServer) createOrder(r *http.Request, w http.ResponseWriter) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var o models.Order
	if err := json.Unmarshal(d, &o); err != nil {
		log.Printf("error unmarshaling data: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(o.OrderedProducts) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("at least one product to order must be provided"))
	}

	co, err := s.r.Create(r.Context(), o)
	if err != nil {
		log.Printf("error creating order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rd, err := json.Marshal(co)
	if err != nil {
		log.Printf("error unmarshaling orders: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(rd); err != nil {
		log.Printf("error writing response: %v", err)
		return
	}
}
func (s *HttpServer) listOrders(ctx context.Context, w http.ResponseWriter) {
	pp, err := s.r.List(ctx)
	if err != nil {
		log.Printf("error retrieving orders: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d, err := json.Marshal(pp)
	if err != nil {
		log.Printf("error unmarshaling orders: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(d); err != nil {
		log.Printf("error writing response: %v", err)
		return
	}
}

func (s *HttpServer) getOrders(ctx context.Context, w http.ResponseWriter, id string) {
	p, err := s.r.Read(ctx, id)
	if err != nil {
		log.Printf("error retrieving order with id '%s': %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	d, err := json.Marshal(p)
	if err != nil {
		log.Printf("error unmarshaling order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(d); err != nil {
		log.Printf("error writing response: %v", err)
		return
	}
}
