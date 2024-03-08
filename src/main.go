package main

import (
	"context"
	"crebito/config"
	"crebito/controller"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	cfg := config.New()
	ctx := context.Background()

	driver, err := neo4j.NewDriverWithContext(fmt.Sprintf("bolt://%s:7687", cfg.DBHostname), neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		log.Panic(err)
	}
	defer driver.Close(ctx)

	for i := 10; i > 0; i-- {
		err := driver.VerifyConnectivity(context.Background())
		if err == nil {
			break
		}

		log.Printf("Error in conectivity=%s\n", err.Error())
		time.Sleep(time.Duration(i) * time.Second)
	}

	exs := driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
		AccessMode:   neo4j.AccessModeRead,
	})
	defer exs.Close(ctx)

	ts := driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
		AccessMode:   neo4j.AccessModeRead,
	})
	defer ts.Close(ctx)

	m := chi.NewMux()

	m.Get("/clientes/{id}/extrato", func(w http.ResponseWriter, r *http.Request) {
		controller.HandleExtract(w, r, exs)
	})

	m.Post("/clientes/{id}/transacoes", func(w http.ResponseWriter, r *http.Request) {
		controller.HandleTransaction(w, r, ts)
	})

	if err := http.ListenAndServe("0.0.0.0:"+cfg.Port, m); err != nil {
		log.Panicf("Error starting server : error=%v", err)
	}
}
