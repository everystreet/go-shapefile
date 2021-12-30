package shp_test

import (
	"encoding/hex"
	"testing"

	"github.com/everystreet/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestDecodePolyline(t *testing.T) {
	buf, err := hex.DecodeString(data)
	require.NoError(t, err)

	p, err := shp.DecodePolyline(buf, 0)
	require.NoError(t, err)

	require.Equal(t, box, p.BoundingBox())

	require.Equal(t, 3, len(p.Parts()))
	pointsEqual(t, part1, p.Parts()[0])
	pointsEqual(t, part2, p.Parts()[1])
	pointsEqual(t, part3, p.Parts()[2])
}

func pointsEqual(t *testing.T, expected, actual []shp.Point) {
	require.Equal(t, normalizePoints(expected), normalizePoints(actual))
}

func normalizePoints(points []shp.Point) []shp.Point {
	out := make([]shp.Point, len(points))
	for i, p := range points {
		out[i] = shp.MakePoint(p.X, p.Y)
	}
	return out
}

// 404 bytes of a polyline taken from Natural Earth.
// Consists of 3 parts with 8, 9 and 5 points, respectively.
const data string = "00000000008066c036936fb6b94932c000000000008066402ec5218a580530c00300000016000000000000000800000011000000000000000080664072d6329b2f1130c00000000000806640aae943ac228e30c0dc06830ea76b6640cd00718a25cd30c06677b1af335766408d99c529150331c099e0404d19536640093d9b559fa330c0ebd1846c17636640560bf797196f30c032d5fc773b6d66409a797db3096130c0000000000080664072d6329b2f1130c07b6b60ab044466409ab1683a3b8131c024b9fc87f44b66409e4143ff045731c01b12f758fa5666401bd82ac1e2a031c082c5e1ccaf516640ca1af5108d2632c033c9c859d83d664036936fb6b94932c06e179aeb342c6640271422e0102a32c0a9c1340c1f296640e10b93a982b931c0c3bb5cc4773566408c321b64926131c07b6b60ab044466409ab1683a3b8131c0f073dae0627966c02ec5218a580530c0653d0a175b7d66c06a4c0ddc748030c000000000008066c0aae943ac228e30c000000000008066c072d6329b2f1130c0f073dae0627966c02ec5218a580530c0"

var box = shp.MakeBoundingBox(-180, -18.28799, 180, -16.020882256741224)

var (
	part1 = shp.Part{
		shp.MakePoint(180, -16.067132663642447),
		shp.MakePoint(180, -16.555216566639196),
		shp.MakePoint(179.36414266196414, -16.801354076946883),
		shp.MakePoint(178.72505936299711, -17.01204167436804),
		shp.MakePoint(178.59683859511713, -16.639150000000004),
		shp.MakePoint(179.0966093629971, -16.433984277547403),
		shp.MakePoint(179.4135093629971, -16.379054277547404),
		shp.MakePoint(180, -16.067132663642447),
	}

	part2 = shp.Part{
		shp.MakePoint(178.12557, -17.50481),
		shp.MakePoint(178.3736, -17.33992),
		shp.MakePoint(178.71806, -17.62846),
		shp.MakePoint(178.55271, -18.15059),
		shp.MakePoint(177.93266000000003, -18.28799),
		shp.MakePoint(177.38146, -18.16432),
		shp.MakePoint(177.28504, -17.72465),
		shp.MakePoint(177.67087, -17.381140000000002),
		shp.MakePoint(178.12557, -17.50481),
	}

	part3 = shp.Part{
		shp.MakePoint(-179.79332010904864, -16.020882256741224),
		shp.MakePoint(-179.9173693847653, -16.501783135649397),
		shp.MakePoint(-180, -16.555216566639196),
		shp.MakePoint(-180, -16.067132663642447),
		shp.MakePoint(-179.79332010904864, -16.020882256741224),
	}
)
