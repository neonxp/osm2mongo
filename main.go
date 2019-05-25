package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func main() {
	dbconnection := flag.String("dbconnection", "mongodb://localhost:27017", "Mongo database name")
	dbname := flag.String("dbname", "osm", "Mongo database name")
	osmfile := flag.String("osmfile", "", "OSM file")
	flag.Parse()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*dbconnection))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	db := client.Database(*dbname)
	if err := read(db, *osmfile); err != nil {
		log.Fatal(err)
	}

}

func read(db *mongo.Database, file string) error {
	nodes := db.Collection("nodes")
	_, _ = nodes.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"osm_id", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	)
	_, _ = nodes.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"coords", bsonx.Int32(1)}},
			Options: options.Index().SetSphereVersion(2).SetSparse(true),
		},
	)

	ways := db.Collection("ways")
	_, _ = ways.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"osm_id", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	)
	_, _ = ways.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"nodes", bsonx.Int32(1)}},
			Options: options.Index().SetSparse(true),
		},
	)

	relations := db.Collection("relations")
	_, _ = nodes.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"osm_id", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	)
	_, _ = nodes.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"members.ref", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	)
	_, _ = nodes.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"members.coords", bsonx.Int32(1)}},
			Options: options.Index().SetSphereVersion(2).SetSparse(true),
		},
	)

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	opts := (new(options.ReplaceOptions)).SetUpsert(true)
	nc := 0
	wc := 0
	rc := 0

	scanner := osmpbf.New(context.Background(), f, 3)
	defer scanner.Close()

	for scanner.Scan() {
		o := scanner.Object()
		switch o := o.(type) {
		case *osm.Way:
			nodes := make([]int64, 0, len(o.Nodes))
			for _, v := range o.Nodes {
				nodes = append(nodes, int64(v.ID))
			}
			w := Way{
				OsmID:     int64(o.ID),
				Tags:      convertTags(o.Tags),
				Nodes:     nodes,
				Timestamp: o.Timestamp,
				Version:   o.Version,
				Visible:   o.Visible,
			}
			if _, err = ways.ReplaceOne(context.Background(), bson.M{"osm_id": int64(o.ID)}, w, opts); err != nil {
				return err
			}
			wc++
		case *osm.Node:
			n := Node{
				OsmID: int64(o.ID),
				Location: Coords{
					Type: "Point",
					Coordinates: []float64{
						o.Lon,
						o.Lat,
					}},
				Tags:      convertTags(o.Tags),
				Version:   o.Version,
				Timestamp: o.Timestamp,
				Visible:   o.Visible,
			}
			if _, err = nodes.ReplaceOne(context.Background(), bson.M{"osm_id": int64(o.ID)}, n, opts); err != nil {
				return err
			}
			nc++
		case *osm.Relation:
			members := make([]Member, len(o.Members))
			for _, v := range o.Members {
				members = append(members, Member{
					Type:        v.Type,
					Version:     v.Version,
					Orientation: v.Orientation,
					Ref:         v.Ref,
					Role:        v.Role,
					Location: Coords{
						Type: "Point",
						Coordinates: []float64{
							v.Lon,
							v.Lat,
						}},
				})
			}
			r := Relation{
				OsmID:     int64(o.ID),
				Tags:      convertTags(o.Tags),
				Version:   o.Version,
				Timestamp: o.Timestamp,
				Visible:   o.Visible,
				Members:   members,
			}
			if _, err = relations.ReplaceOne(context.Background(), bson.M{"osm_id": int64(o.ID)}, r, opts); err != nil {
				return err
			}
			rc++
		}
		fmt.Printf("\rNodes: %d Ways: %d Relations: %d", nc, wc, rc)
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return scanErr
	}
	return nil
}

func convertTags(tags osm.Tags) map[string]string {
	result := make(map[string]string, len(tags))
	for _, t := range tags {
		result[t.Key] = t.Value
	}
	return result
}
