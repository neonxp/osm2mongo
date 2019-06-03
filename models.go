package main

import (
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/osm"
)

type Coords struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type ItemType string

const (
	NodeType     ItemType = "node"
	WayType      ItemType = "way"
	RelationType ItemType = "relation"
)

type ID struct {
	ID      int64    `json:"id" bson:"id"`
	Type    ItemType `json:"type" bson:"type"`
	Version int      `json:"version" bson:"version"`
}

type Object struct {
	ID        ID        `json:"_id" bson:"_id"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Tags      []Tag     `json:"tags" bson:"tags"`
	Location  Coords    `json:"location,omitempty" bson:"location,omitempty"`
	Nodes     []int64   `json:"nodes,omitempty" bson:"nodes,omitempty"`
	Members   []Member  `json:"members,omitempty" bson:"members,omitempty"`
}

type Member struct {
	Type osm.Type `json:"type" bson:"type"`
	Ref  int64    `json:"ref" bson:"ref"`
	Role string   `json:"role" bson:"role"`

	Location *Coords `json:"location" bson:"location"`

	// Orientation is the direction of the way around a ring of a multipolygon.
	// Only valid for multipolygon or boundary relations.
	Orientation orb.Orientation `json:"orienation" bson:"orienation"`
}

type Tag struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}
