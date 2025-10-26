package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	go func() {
		r := http.NewServeMux()

		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)

		log.Println(http.ListenAndServe("localhost:6060", r))
	}()

	db, err := NewPostgres(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.PingContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Postgres")

	repo := NewCommand(db, DollarWildcard)

	res, err := repo.Get(context.Background(), WithID(1))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

}

func NewPostgres(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	return db, err
}
