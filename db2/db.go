package db2

import (
	"context"
	"database/sql"
	"log"

	_ "embed"

	"github.com/lunardoesdev/singerbench/mydb"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var initsql string
var Queries *mydb.Queries
var Db sql.DB

func init() {
	ctx := context.Background()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("%w\n", err)
	}

	if _, err := db.ExecContext(ctx, initsql); err != nil {
		log.Fatalf("%w\n", err)
	}

	Queries = mydb.New(db)
}

func Addproxy(ctx context.Context, proxy string) error {
	prox := sql.NullString{
		String: proxy,
		Valid:  true,
	}
	_, err := Queries.GetProxyIdByLink(ctx, prox)
	if err != nil {
		//proxy doesnt exist
		if err := Queries.AddProxy(ctx, prox); err != nil {
			return err
		}
	}

	return nil
}
