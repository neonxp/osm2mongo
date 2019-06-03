package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func write(ctx context.Context, db *mongo.Database, insertCh chan Object, initial bool, blockSize int, worker int) error {
	nodes := db.Collection("items")
	opts := (new(options.BulkWriteOptions)).SetOrdered(false)
	buf := make([]mongo.WriteModel, 0, blockSize)
	ic := 0
	for {
		select {
		case w := <-insertCh:
			if initial {
				um := mongo.NewInsertOneModel()
				um.SetDocument(w)
				buf = append(buf, um)
			} else {
				um := mongo.NewUpdateOneModel()
				um.SetUpsert(true)
				um.SetUpdate(w)
				um.SetFilter(bson.M{"osm_id": w.ID})
				buf = append(buf, um)
			}
		case <-ctx.Done():
			if len(buf) > 0 {
				log.Printf("Worker: %d\tSaving last info in buffers...", worker)
				if _, err := nodes.BulkWrite(context.Background(), buf, opts); err != nil {
					return err
				}
			}
			log.Printf("Worker: %d\tDone", worker)
			return nil
		}
		if len(buf) == blockSize {
			ic++
			log.Printf("Worker: %d\tWriting block %d (%d objects)", worker, ic, ic*blockSize)
			if _, err := nodes.BulkWrite(context.Background(), buf, opts); err != nil {
				return err
			}
			buf = make([]mongo.WriteModel, 0)
		}
	}
}
