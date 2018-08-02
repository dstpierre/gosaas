package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dstpierre/gosaas/cache"
	"github.com/dstpierre/gosaas/controllers"
	"github.com/dstpierre/gosaas/data"
)

func main() {
	dn := flag.String("driver", "postgres", "name of the database driver to use, postgres or mongo are supported")
	ds := flag.String("datasource", "", "database connection string")
	q := flag.Bool("queue", false, "set as queue pub/sub subscriber and task executor")
	flag.Parse()

	if len(*dn) == 0 || len(*ds) == 0 {
		flag.Usage()
		return
	}

	api := controllers.NewAPI()

	// open the database connection
	db := &data.DB{}

	if err := db.Open(*dn, *ds); err != nil {
		log.Fatal("unable to connect to the database:", err)
	}

	api.DB = db

	// Set as Redis pub/sub subscriber for the queue executor if q is true
	cache.New(*q)

	if err := http.ListenAndServe(":8080", api); err != nil {
		log.Println(err)
	}
}
