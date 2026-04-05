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

func IterateProxies(pagesize int64) func(func(mydb.Proxy) bool) {
	return func(yield func(mydb.Proxy) bool) {
		var curPage int64 = 0
		var curElement int64 = 0

		ctx := context.Background()

		count, err := Queries.CountProxies(ctx)
		if err != nil {
			goto finish
		}

		for curElement < count {
			page, err := Queries.ListProxies(ctx, mydb.ListProxiesParams{
				Offset: curPage * pagesize,
				Limit:  pagesize,
			})
			if err != nil {
				goto finish
			}

			for _, v := range page {
				yield(v)
				curElement++
			}
			curPage++
		}

	finish:
	}
}

func IterateSubscriptions(pagesize int64) func(func(mydb.Subscription) bool) {
	return func(yield func(mydb.Subscription) bool) {
		var curPage int64 = 0
		var curElement int64 = 0

		ctx := context.Background()

		count, err := Queries.CountSubscriptions(ctx)
		if err != nil {
			goto finish2
		}

		for curElement < count {
			page, err := Queries.ListSubscriptoins(ctx, mydb.ListSubscriptoinsParams{
				Offset: curPage * pagesize,
				Limit:  pagesize,
			})
			if err != nil {
				goto finish2
			}

			for _, v := range page {
				yield(v)
				curElement++
			}
			curPage++
		}

	finish2:
	}
}

func IterateMeasurements(pagesize int64) func(func(mydb.Measurement) bool) {
	return func(yield func(mydb.Measurement) bool) {
		var curPage int64 = 0
		var curElement int64 = 0

		ctx := context.Background()

		count, err := Queries.CountMeasurements(ctx)
		if err != nil {
			goto finish3
		}

		for curElement < count {
			page, err := Queries.ListMeasurements(ctx, mydb.ListMeasurementsParams{
				Offset: curPage * pagesize,
				Limit:  pagesize,
			})
			if err != nil {
				goto finish3
			}

			for _, v := range page {
				yield(v)
				curElement++
			}
			curPage++
		}

	finish3:
	}
}
