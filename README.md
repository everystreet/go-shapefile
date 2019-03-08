# go-shapefile

`go-shapefile` is a Go parser for the "shapefile" GIS file format. The package does not currently support writing of files.

This package has been developed primarily to work with [Natural Earth](https://www.naturalearthdata.com/), so may only contain a subset of features relevant to those data files. The "shapefile" format is actually a collection of files, of which this package currently supports the "shape" (.shp), "attribute" (.dbf) and character encoding (.cpg) files.

## Shape file (.shp)

The .shp file contains the geometry data in the form of variable-length records. A single record represents a particular shape type, although the records in a single file must all represent the same type. `go-shapefile` supports the following types:

| Shape type  | Supported          |
| ----------- |:------------------:|
| Point       | :heavy_check_mark: |
| Polyline    | :heavy_check_mark: |
| Polygon     | :heavy_check_mark: |
| MultiPoint  | :x:                |
| PointZ      | :x:                |
| PolylineZ   | :x:                |
| PolygonZ    | :x:                |
| MultiPointZ | :x:                |
| PointM      | :x:                |
| PolylineM   | :x:                |
| PolygonM    | :x:                |
| MultiPointM | :x:                |
| MultiPatch  | :x:                |

Format specification: https://www.esri.com/library/whitepapers/pdfs/shapefile.pdf

## Attribute file (.dbf)

The .dbf file contains attributes for each shape in the .shp file. Attributes are stored in the form of records, which consist of a number of fields. Field names and values are not standardized - they are specified as part of the .dbf file to suit the particular use case.

This file uses a format called "dBase", of which there are several variations in varying degrees of usage. The most common are dBase IV and dBase V, but `go-shapefile` currently only supports IV. Below is an overview of the supported field types:

| Field type       | Supported          |
| ---------------- |:------------------:|
| Character/string | :heavy_check_mark: |
| Numeric          | :heavy_check_mark: |
| Date             | :x:                |
| Floating point   | :x:                |
| Logical          | :x:                |
| Memo             | :x:                |

Note that dBase V contains many more field types.

## Character endoding file (.cpg)

The .cpg file is optional and contains the character encoding used inside the .dbf file. By default, and in this file's absense, the character encoding is assumed to be ASCII, but this file can be used to support Unicode strings. `go-shapefile` supports the following encodings:

* ASCII
* UTF-8
