package main

import (
	"context"
	"log"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

func read(ctx context.Context, file string, insertCh chan Object, concurrency int, layers []string) error {
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
			if !layersToImport["ways"] || !o.Visible {
				continue
			}
			nodes := make([]int64, 0, len(o.Nodes))
			for _, v := range o.Nodes {
				nodes = append(nodes, int64(v.ID))
			}

			w := Object{
				ID:        ID{ID: int64(o.ID), Type: WayType, Version: o.Version},
				Tags:      convertTags(o.Tags),
				Timestamp: o.Timestamp,
				Nodes:     nodes,
			}
			insertCh <- w
		case *osm.Node:
			if !layersToImport["nodes"] || !o.Visible {
				continue
			}
			w := Object{
				ID:        ID{ID: int64(o.ID), Type: NodeType, Version: o.Version},
				Tags:      convertTags(o.Tags),
				Timestamp: o.Timestamp,
				Location: Coords{
					Type: "Point",
					Coordinates: []float64{
						o.Lon,
						o.Lat,
					}},
			}
			insertCh <- w
		case *osm.Relation:
			if !layersToImport["relations"] || !o.Visible {
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
					Orientation: v.Orientation,
					Ref:         v.Ref,
					Role:        v.Role,
					Location:    location,
				})
			}
			w := Object{
				ID:        ID{ID: int64(o.ID), Type: RelationType, Version: o.Version},
				Tags:      convertTags(o.Tags),
				Timestamp: o.Timestamp,
				Members:   members,
			}
			insertCh <- w
		}
	}
	log.Println("Read done")
	scanErr := scanner.Err()
	if scanErr != nil {
		return scanErr
	}
	return nil
}
