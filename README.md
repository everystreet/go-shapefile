# go-shapefile

`go-shapefile` is a Go parser for the "shapefile" GIS file format. The package can be used to read shapefiles, and optionally convert them into GeoJSON "features" (supported by [`go-geojson`](https://github.com/everystreet/go-geojson)) - allowing a shapefile to be written directly to databases such as PostGIS, MongoDB, or Couchbase, etc.

**The package does not currently support writing of files.**

## Usage

How you choose to use this package will depend on your use case. `go-shapefile` supports reading the shapefiles in the following forms:

* .zip file containing mandatory .shp and .dbf files, with optional .cpg file
* Unzipped .shp and .dbf files, with optional character encoding
* .shp and .dbf files separately, with optional character encoding

### Basic example

Reading a zipped shapefile is achieved by using the `ZipScanner`. The example below shows a basic example of this, where error handling has been omitted for brevity.

```go
file, err := os.Open("path/to/ne_110m_admin_0_sovereignty.zip")
stat, err := r.Stat()

// Create new ZipScanner
// The filename can be replaced with an empty string if you don't want to check filenames inside the zip file
scanner := shapefile.NewZipScanner(file, stat.Size(), "ne_110m_admin_0_sovereignty.zip")

// Optionally get file info: shape type, number of records, bounding box, etc.
info, err := scanner.Info()
fmt.Println(info)

// Start the scanner
err = scanner.Scan()

// Call Record() to get each record in turn, until either the end of the file, or an error occurs
for {
    record := scanner.Record()
    if record == nil {
        break
    }

    // Each record contains a shape (from .shp file) and attributes (from .dbf file)
    fmt.Println(record)
}

// Err() returns the first error encountered during calls to Record()
err = scanner.Err()
```

### GeoJSON example

Using the example above, we can optionally convert shapefile records to GeoJSON features. `go-shapefile` achieves this by using [`go-geojson`](https://github.com/everystreet/go-geojson), meaning that you can use the standard `json.Marshal` to produce a JSON object that can be understood by any software that can work with the GeoJSON standard.

```go
record := scanner.Record()
feature := record.GeoJSONFeature()

jsonData, err := json.Marshal(feature)
fmt.Println(string(jsonData))
```

## Features

This package has been primarily developed to work with [Natural Earth](https://www.naturalearthdata.com/), so may only contain the subset of shapefile features relevant to those data files. The "shapefile" format is actually a collection of files, of which this package currently supports the "shape" (.shp), "attribute" (.dbf) and character encoding (.cpg) files.

### Shape file (.shp)

The .shp file contains the geometry data in the form of variable-length records. A single record represents a particular shape type, although the records in a single file must all represent the same type. `go-shapefile` supports the following types:

| Shape type  |     Supported      |
| ----------- | :----------------: |
| Point       | :heavy_check_mark: |
| Polyline    | :heavy_check_mark: |
| Polygon     | :heavy_check_mark: |
| MultiPoint  |        :x:         |
| PointZ      |        :x:         |
| PolylineZ   |        :x:         |
| PolygonZ    |        :x:         |
| MultiPointZ |        :x:         |
| PointM      |        :x:         |
| PolylineM   |        :x:         |
| PolygonM    |        :x:         |
| MultiPointM |        :x:         |
| MultiPatch  |        :x:         |

Format specification: [https://www.esri.com/library/whitepapers/pdfs/shapefile.pdf](./docs/shapefile.pdf).

### Attribute file (.dbf)

The .dbf file contains attributes for each shape in the .shp file. Attributes are stored in the form of records, which consist of a number of fields. Field names and values are not standardized - they are specified as part of the .dbf file to suit the particular use case.

This file uses a format called "dBase", of which there are several variations in varying degrees of usage. The most common are dBase IV and dBase V, but `go-shapefile` currently only supports IV. Below is an overview of the supported field types:

| Field type       |     Supported      |
| ---------------- | :----------------: |
| Character/string | :heavy_check_mark: |
| Numeric          | :heavy_check_mark: |
| Date             | :heavy_check_mark: |
| Floating point   | :heavy_check_mark: |
| Logical          |        :x:         |
| Memo             |        :x:         |

Note that dBase V contains many more field types.

### Character endoding file (.cpg)

The .cpg file is optional and contains the character encoding used inside the .dbf file. By default, and in this file's absense, the character encoding is assumed to be ASCII, but this file can be used to support Unicode strings. `go-shapefile` supports the encoding labels defined by https://encoding.spec.whatwg.org/#names-and-labels.
