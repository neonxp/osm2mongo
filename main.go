package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/neonxp/rutina"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	dbconnection := flag.String("dbconnection", "mongodb://localhost:27017", "Mongo database name")
	dbname := flag.String("dbname", "map", "Mongo database name")
	osmfile := flag.String("osmfile", "", "OSM file")
	initial := flag.Bool("initial", false, "Is initial import")
	indexes := flag.Bool("indexes", false, "Create indexes")
	layersString := flag.String("layers", "nodes,ways,relations", "Layers to import")
	blockSize := flag.Int("block", 1000, "Block size to bulk write")
	concurrency := flag.Int("concurrency", 32, "Concurrency read and write")
	flag.Parse()
	layers := strings.Split(*layersString, ",")
	r := rutina.New()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*dbconnection))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	db := client.Database(*dbname)

	if *indexes {
		if err := createIndexes(db); err != nil {
			log.Fatal(err)
		}
		log.Println("Indexes created")
	}

	log.Printf("Started import file %s to db %s", *osmfile, *dbname)
	nodesCh := make(chan Node, 1)
	waysCh := make(chan Way, 1)
	relationsCh := make(chan Relation, 1)

	for i := 0; i < *concurrency; i++ {
		worker := i
		r.Go(func(ctx context.Context) error {
			return write(ctx, db, nodesCh, waysCh, relationsCh, *initial, *blockSize, worker)
		})
	}

	r.Go(func(ctx context.Context) error {
		return read(ctx, *osmfile, nodesCh, waysCh, relationsCh, *concurrency, layers)
	})
	if err := r.Wait(); err != nil {
		log.Fatal(err)
	}
}
