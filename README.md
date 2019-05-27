# OpenStreetMaps to Mongo

Simple loader from osm dump file to mongodb. Based on https://github.com/paulmach/osm package.

## Build

`go build -o osm2go`

## Usage

`./osm2go -osmfile PATH_TO_OSM_FILE [-dbconnection mongodb://localhost:27017] [-dbname osm] [-initial=true] [-concurrency=16] [-block=1000]`

* `osmfile` required, path to *.osm or *.osm.pbf file
* `dbconnection` optional, mongodb connection string (default: `mongodb://localhost:27017`) 
* `dbname` optional, mongodb database name (default: `osm`) 
* `initial` optional, use insert instead upsert. Faster, but not check if item exists (default: `false`)
* `concurrency` optional, parallel read processes (default, `16`)
* `block` optional, block size to bulk write (default: `1000`)

## Example

```
# ./osm2mgo -osmfile ~/Downloads/RU.pbf
Nodes: 1294069 Ways: 0 Relations: 0
```