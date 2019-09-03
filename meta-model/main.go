package main

import (
	"fmt"
	"github.com/facedamon/meta-model/pkg"
	"github.com/facedamon/meta-model/routers"
	_ "github.com/facedamon/meta-model/sql"
	"net/http"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", pkg.HTTPPort),
		Handler:        router,
		ReadTimeout:    pkg.ReadTimeout,
		WriteTimeout:   pkg.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
