package main

import (
	p "boilgopher/storage/postgresutil"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	if err := p.Migrate(); err != nil {
		log.Fatal(err)
	}
}
