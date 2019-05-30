package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func createIndexes(db *mongo.Database) error {
	opts := options.CreateIndexes().SetMaxTime(1000)
	nodes := db.Collection("nodes")
	log.Println("creating indexes for nodes")
	created, err := nodes.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{"osm_id", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true).SetUnique(false),
		}, {
			Keys:    bsonx.Doc{{"osm_id", bsonx.Int32(-1)}, {"version", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true).SetUnique(false),
		}, {
			Keys:    bsonx.Doc{{"tags", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true),
		},
	}, opts)
	if err != nil {
		return err
	}
	log.Println(created)
	log.Println("creating geoindexes for nodes")
	if err := geoIndex(nodes, "location"); err != nil {
		return err
	}

	log.Println("creating indexes for ways")
	ways := db.Collection("ways")
	created, err = ways.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{"osm_id", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true).SetUnique(false),
		}, {
			Keys:    bsonx.Doc{{"osm_id", bsonx.Int32(-1)}, {"version", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true).SetUnique(false),
		}, {
			Keys:    bsonx.Doc{{"tags", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true),
		},
	}, opts)
	if err != nil {
		return err
	}
	log.Println(created)

	relations := db.Collection("relations")
	log.Println("creating geoindexes for relations")
	created, err = relations.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{"osm_id", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true).SetUnique(false),
		}, {
			Keys:    bsonx.Doc{{"osm_id", bsonx.Int32(-1)}, {"version", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true).SetUnique(false),
		}, {
			Keys:    bsonx.Doc{{"tags", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true),
		},
		{
			Keys:    bsonx.Doc{{"members.ref", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true),
		},
	}, opts)
	if err != nil {
		return err
	}
	log.Println(created)
	if err := geoIndex(relations, "members.coords"); err != nil {
		return err
	}
	log.Println("indexes created")
	return nil
}

func geoIndex(col *mongo.Collection, key string) error {
	_, err := col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{{
				Key: key, Value: bsonx.String("2dsphere"),
			}},
			Options: options.Index().SetSphereVersion(2).SetSparse(true).SetBackground(true),
		},
	)
	return err
}
