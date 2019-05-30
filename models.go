package main

import (
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/osm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coords struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type Node struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OsmID     int64              `json:"osm_id" bson:"osm_id"`
	Visible   bool               `json:"visible" bson:"visible"`
	Version   int                `json:"version,omitempty" bson:"version,omitempty"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	Tags      []Tag              `json:"tags,omitempty" bson:"tags,omitempty"`
	Location  Coords             `json:"location" bson:"location"`
}

type Way struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OsmID     int64              `json:"osm_id" bson:"osm_id"`
	Visible   bool               `json:"visible" bson:"visible"`
	Version   int                `json:"version" bson:"version"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	Nodes     []int64            `json:"nodes" bson:"nodes"`
	Tags      []Tag              `json:"tags" bson:"tags"`
}

type Relation struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OsmID     int64              `json:"osm_id" bson:"osm_id"`
	Visible   bool               `json:"visible" bson:"visible"`
	Version   int                `json:"version" bson:"version"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	Members   []Member           `json:"members" bson:"members"`
	Tags      []Tag              `json:"tags" bson:"tags"`
}

type Member struct {
	Type osm.Type `json:"type" bson:"type"`
	Ref  int64    `json:"ref" bson:"ref"`
	Role string   `json:"role" bson:"role"`

	Version  int
	Location *Coords `json:"location,omitempty" bson:"location,omitempty"`

	// Orientation is the direction of the way around a ring of a multipolygon.
	// Only valid for multipolygon or boundary relations.
	Orientation orb.Orientation `json:"orienation,omitempty" bson:"orienation,omitempty"`
}

type Tag struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}
