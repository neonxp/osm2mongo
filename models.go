package main

import (
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/osm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coords struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type Node struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OsmID     int64              `bson:"osm_id"`
	Visible   bool               `bson:"visible"`
	Version   int                `bson:"version,omitempty"`
	Timestamp time.Time          `bson:"timestamp"`
	Tags      map[string]string  `bson:"tags,omitempty"`
	Location  Coords             `bson:"location"`
}

type Way struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OsmID     int64              `bson:"osm_id"`
	Visible   bool               `bson:"visible"`
	Version   int                `bson:"version"`
	Timestamp time.Time          `bson:"timestamp"`
	Nodes     []int64            `bson:"nodes"`
	Tags      map[string]string  `bson:"tags"`
}

type Relation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OsmID     int64              `bson:"osm_id"`
	Visible   bool               `bson:"visible"`
	Version   int                `bson:"version"`
	Timestamp time.Time          `bson:"timestamp"`
	Members   []Member           `bson:"members"`
	Tags      map[string]string  `bson:"tags"`
}

type Member struct {
	Type osm.Type `bson:"type"`
	Ref  int64    `bson:"ref"`
	Role string   `bson:"role"`

	Version  int
	Location Coords `bson:"location"`

	// Orientation is the direction of the way around a ring of a multipolygon.
	// Only valid for multipolygon or boundary relations.
	Orientation orb.Orientation `bson:"orienation,omitempty"`
}
