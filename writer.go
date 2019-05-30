package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func write(ctx context.Context, db *mongo.Database, nodesCh chan Node, waysCh chan Way, relationsCh chan Relation, initial bool, blockSize int, worker int) error {
	nodes := db.Collection("nodes")
	ways := db.Collection("ways")
	relations := db.Collection("relations")
	opts := (new(options.BulkWriteOptions)).SetOrdered(false)
	nodesBuffer := make([]mongo.WriteModel, 0, blockSize)
	waysBuffer := make([]mongo.WriteModel, 0, blockSize)
	relationsBuffer := make([]mongo.WriteModel, 0, blockSize)
	nc := 0
	wc := 0
	rc := 0
	for {
		select {
		case w := <-waysCh:
			if initial {
				um := mongo.NewInsertOneModel()
				um.SetDocument(w)
				waysBuffer = append(waysBuffer, um)
			} else {
				um := mongo.NewUpdateOneModel()
				um.SetUpsert(true)
				um.SetUpdate(w)
				um.SetFilter(bson.M{"osm_id": w.OsmID})
				waysBuffer = append(waysBuffer, um)
			}

		case n := <-nodesCh:
			if initial {
				um := mongo.NewInsertOneModel()
				um.SetDocument(n)
				nodesBuffer = append(nodesBuffer, um)
			} else {
				um := mongo.NewUpdateOneModel()
				um.SetUpsert(true)
				um.SetUpdate(n)
				um.SetFilter(bson.M{"osm_id": n.OsmID})
				nodesBuffer = append(nodesBuffer, um)
			}
		case r := <-relationsCh:
			if initial {
				um := mongo.NewInsertOneModel()
				um.SetDocument(r)
				relationsBuffer = append(relationsBuffer, um)
			} else {
				um := mongo.NewUpdateOneModel()
				um.SetUpsert(true)
				um.SetUpdate(r)
				um.SetFilter(bson.M{"osm_id": r.OsmID})
				relationsBuffer = append(relationsBuffer, um)
			}
		case <-ctx.Done():
			log.Printf("[%d] saving last info in buffers...", worker)
			if _, err := nodes.BulkWrite(context.Background(), nodesBuffer, opts); err != nil {
				return err
			}
			if _, err := ways.BulkWrite(context.Background(), waysBuffer, opts); err != nil {
				return err
			}
			if _, err := relations.BulkWrite(context.Background(), relationsBuffer, opts); err != nil {
				return err
			}
			log.Printf("[%d] Done", worker)
			return nil
		}
		if len(nodesBuffer) == blockSize {
			nc++
			log.Printf("[%d] nodes %d ways %d relations %d", worker, nc, wc, rc)
			if _, err := nodes.BulkWrite(context.Background(), nodesBuffer, opts); err != nil {
				return err
			}
			nodesBuffer = make([]mongo.WriteModel, 0)
		}
		if len(waysBuffer) == blockSize {
			wc++
			log.Printf("[%d] nodes %d ways %d relations %d", worker, nc, wc, rc)
			if _, err := ways.BulkWrite(context.Background(), waysBuffer, opts); err != nil {
				return err
			}
			waysBuffer = make([]mongo.WriteModel, 0)
		}
		if len(relationsBuffer) == blockSize {
			rc++
			log.Printf("[%d] nodes %d ways %d relations %d", worker, nc, wc, rc)
			if _, err := relations.BulkWrite(context.Background(), relationsBuffer, opts); err != nil {
				return err
			}
			relationsBuffer = make([]mongo.WriteModel, 0)
		}
	}
}
