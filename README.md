# OpenStreetMaps to Mongo

Simple loader from osm dump file to mongodb. Based on https://github.com/paulmach/osm package.

## Build

`go build -o osm2mgo`

## Usage

`./osm2mgo flags`

### Flags:

* `-osmfile string` Path to OSM file (PBF format only) (default "./RU.osm.pbf")
* `-dbconnection string` Mongo database name (default "mongodb://localhost:27017")
* `-dbname string` Mongo database name (default "map")
* `-initial` Is initial import?
* `-indexes` Create indexes
* `-layers string` Layers to import (default "nodes,ways,relations")
* `-concurrency int` Workers count (default 32)
* `-block int` Block size to bulk write (default 1000)

## Example

```
# ./osm2mgo -osmfile ~/Downloads/RU.pbf
Nodes: 1294069 Ways: 0 Relations: 0
```