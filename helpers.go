package main

import "github.com/paulmach/osm"

func convertTags(tags osm.Tags) []Tag {
	result := make([]Tag, 0, len(tags))
	for _, t := range tags {
		result = append(result, Tag{
			Key:   t.Key,
			Value: t.Value,
		})
	}
	return result
}
