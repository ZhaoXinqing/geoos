package wkt

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spatial-go/geoos"
)

// UnmarshalString encode to geom
func UnmarshalString(s string) (geoos.Geometry, error) {
	p := Parser{NewLexer(strings.NewReader(s))}
	return p.Parse()
}

// MarshalString decode to string
func MarshalString(geom geoos.Geometry) string {
	buf := bytes.NewBuffer(nil)
	wkt(buf, geom)
	return buf.String()
}

func wkt(buf *bytes.Buffer, geom geoos.Geometry) {
	switch g := geom.(type) {
	case geoos.Point:
		_, _ = fmt.Fprintf(buf, "POINT(%g %g)", g.Lon(), g.Lat())
	case geoos.MultiPoint:
		if len(g) == 0 {
			buf.Write([]byte(`MULTIPOINT EMPTY`))
			return
		}

		buf.Write([]byte(`MULTIPOINT(`))
		for i, p := range g {
			if i != 0 {
				buf.WriteByte(',')
			}
			_, _ = fmt.Fprintf(buf, "(%g %g)", p.Lon(), p.Lat())
		}
		buf.WriteByte(')')
	case geoos.LineString:
		if len(g) == 0 {
			buf.Write([]byte(`LINESTRING EMPTY`))
			return
		}

		buf.Write([]byte(`LINESTRING`))
		writeLineString(buf, g)
	case geoos.MultiLineString:
		if len(g) == 0 {
			buf.Write([]byte(`MULTILINESTRING EMPTY`))
			return
		}

		buf.Write([]byte(`MULTILINESTRING(`))
		for i, ls := range g {
			if i != 0 {
				buf.WriteByte(',')
			}
			writeLineString(buf, ls)
		}
		buf.WriteByte(')')
	case geoos.Ring:
		wkt(buf, geoos.Polygon{g})
	case geoos.Polygon:
		if len(g) == 0 {
			buf.Write([]byte(`POLYGON EMPTY`))
			return
		}

		buf.Write([]byte(`POLYGON(`))
		for i, r := range g {
			if i != 0 {
				buf.WriteByte(',')
			}
			writeLineString(buf, geoos.LineString(r))
		}
		buf.WriteByte(')')
	case geoos.MultiPolygon:
		if len(g) == 0 {
			buf.Write([]byte(`MULTIPOLYGON EMPTY`))
			return
		}

		buf.Write([]byte(`MULTIPOLYGON(`))
		for i, p := range g {
			if i != 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte('(')
			for j, r := range p {
				if j != 0 {
					buf.WriteByte(',')
				}
				writeLineString(buf, geoos.LineString(r))
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case geoos.Collection:
		if len(g) == 0 {
			buf.Write([]byte(`GEOMETRYCOLLECTION EMPTY`))
			return
		}
		buf.Write([]byte(`GEOMETRYCOLLECTION(`))
		for i, c := range g {
			if i != 0 {
				buf.WriteByte(',')
			}
			wkt(buf, c)
		}
		buf.WriteByte(')')
	default:
		panic("unsupported type")
	}
}

func writeLineString(buf *bytes.Buffer, ls geoos.LineString) {
	buf.WriteByte('(')
	for i, p := range ls {
		if i != 0 {
			buf.WriteByte(',')
		}

		_, _ = fmt.Fprintf(buf, "%g %g", p.Lon(), p.Lat())
	}
	buf.WriteByte(')')
}
