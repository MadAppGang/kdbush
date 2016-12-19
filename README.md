# KDBsuh

[Godoc is here](https://godoc.org/github.com/MadAppGang/kdbush)

Package kdbush provides a very fast static spatial index for 2D points based on a flat KD-tree.
Very fast, but limited:

-  Points only, no rectangles
-  2 dimensional
- indexing 16-40 times faster then  rtreego(https: github.com/dhconnelly/rtreego) (TODO: benchmark)
- Implements radius search  (rtreego and go.geo only have range search)


There are three amazing other options for geospatial indexing:

- [Tile38](tile38.com)
- [go.geo](https://github.com/paulmach/go.geo/tree/master/quadtree)
- [rtreego](https://github.com/dhconnelly/rtreego)



This implementation is based on:
 - JS library: https: github.com/mourner/kdbush
 - C++11 port: https: github.com/mourner/kdbush.hpp


##Create Index example
All Items should implement Point interface:
```go
type Point interface {
	Coordinates() (X, Y float64)
}
```

Package represents simple struct, that implements that protocol, you could use it:
```go
type SimplePoint struct {
	X, Y float64
}
func (sp *SimplePoint)Coordinates()(float64, float64) {
	return sp.X, sp.Y
}

```


To create index, just supply array of points and nodeSize.
NodeSize is size of the KD-tree node, 64 by default. Higher means faster indexing but slower search, and vise versa.

```go
points := []Point{
  &SimplePoint{X:10,Y:10}, //0
  &SimplePoint{X:15,Y:11}, //1
  &SimplePoint{X:1,Y:22},
  &SimplePoint{X:22,Y:22},
  &SimplePoint{X:34,Y:12},
  &SimplePoint{X:19,Y:19}, //5
  &SimplePoint{X:32,Y:34},
}
bush := NewBush(points, 10)
```

##Search in range


kdbush implements Range function to find points in bounding box.
Returns an array of indices that refer to the items in the original points input slice.

```go
points := []Point{
  &SimplePoint{X:10,Y:10}, //0
  &SimplePoint{X:15,Y:11}, //1
  &SimplePoint{X:1,Y:22},
  &SimplePoint{X:22,Y:22},
  &SimplePoint{X:34,Y:12},
  &SimplePoint{X:19,Y:19}, //5
  &SimplePoint{X:32,Y:34},
}
bush := NewBush(points, 10)
result :=  bush.Range(10, 10, 21, 21)
fmt.Println(result)
// Output: [0 1 5]

```
