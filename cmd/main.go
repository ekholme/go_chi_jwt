package main

import (
	gochijwt "github.com/ekholme/go_chi_jwt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	s := gochijwt.NewServer(r)

	s.Router.Use(middleware.Logger)

	s.Run()
}
