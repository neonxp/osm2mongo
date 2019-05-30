# OpenStreetMaps to Mongo

Simple loader from osm dump file to mongodb. Based on https://github.com/paulmach/osm package.

## Build

`go build -o osm2go`

## Usage

`./osm2go -osmfile=PATH_TO_OSM_FILE`

All flags:

* `-osmfile` (required) OSM file
* `-initial` (default:false) Is initial import (uses insert, not upsert)
* `-indexes` (default:false) Create indexes (needs only first time)
* `-dbconnection` (default:"mongodb://localhost:27017") Mongo database name
* `-dbname` (default:"map") Mongo database name
* `-layers` (default:"nodes,ways,relations") Layers to import
* `-block` (default:1000) Block size to bulk write
* `-concurrency` (default:32) Concurrency read and write

## Example

```
# ./osm2mgo -osmfile ~/Downloads/RU.pbf
Nodes: 1294069 Ways: 0 Relations: 0
```