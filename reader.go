package main

import (
	"context"
	"log"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

func read(ctx context.Context, file string, nodesCh chan Node, waysCh chan Way, relationsCh chan Relation, concurrency int, layers []string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := osmpbf.New(context.Background(), f, concurrency)
	defer scanner.Close()

	layersToImport := map[string]bool{
		"ways":      false,
		"nodes":     false,
		"relations": false,
	}

	for _, l := range layers {
		layersToImport[l] = true
	}

	for scanner.Scan() {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		o := scanner.Object()
		switch o := o.(type) {
		case *osm.Way:
			if !layersToImport["ways"] {
				continue
			}
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
			waysCh <- w
		case *osm.Node:
			if !layersToImport["nodes"] {
				continue
			}
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
			nodesCh <- n
		case *osm.Relation:
			if !layersToImport["relations"] {
				continue
			}
			members := make([]Member, 0, len(o.Members))
			for _, v := range o.Members {
				var location *Coords
				if v.Lat != 0.0 && v.Lon != 0.0 {
					location = &Coords{
						Type: "Point",
						Coordinates: []float64{
							v.Lon,
							v.Lat,
						}}
				}
				members = append(members, Member{
					Type:        v.Type,
					Version:     v.Version,
					Orientation: v.Orientation,
					Ref:         v.Ref,
					Role:        v.Role,
					Location:    location,
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
			relationsCh <- r
		}
	}
	log.Println("Read done")
	scanErr := scanner.Err()
	if scanErr != nil {
		return scanErr
	}
	return nil
}
