package kdbush

import (
	"math"
)


// Interface, that should be implemented by indexing structure
// It's just simply returns points coordinates
// Called once, only when index created, so you could calc values on the fly for this interface
type Point interface {
	Coordinates() (X, Y float64)
}


// Minimal struct, that implements Point interface
type SimplePoint struct {
	X, Y float64
}

// SimplePoint's  implementation of Point interface
func (sp *SimplePoint)Coordinates()(float64, float64) {
	return sp.X, sp.Y
}


// A very fast static spatial index for 2D points based on a flat KD-tree.
// Points only, no rectangles
// static (no add, remove items)
// 2 dimensional
// indexing 16-40 times faster then  rtreego(https://github.com/dhconnelly/rtreego) (TODO: benchmark)
type KDBush struct {
	NodeSize int
	Points   []Point

	idxs   []int     //array of indexes
	coords []float64 //array of coordinates
}


// Create new index from points
// Structure don't copy points itself, copy only coordinates
// Returns pointer to new KDBush index object, all data in it already indexed
// Input:
// points - slice of objects, that implements Point interface
// nodeSize  - size of the KD-tree node, 64 by default. Higher means faster indexing but slower search, and vise versa.
func NewBush(points []Point, nodeSize int) *KDBush {
	b := KDBush{}
	b.buildIndex(points, nodeSize)
	return &b
}

// Finds all items within the given bounding box and returns an array of indices that refer to the items in the original points input slice.
func (bush *KDBush) Range(minX, minY, maxX, maxY float64) []int {
	stack := []int{0, len(bush.idxs) - 1, 0}
	result := []int{}
	var x, y float64

	for len(stack) > 0 {
		axis := stack[len(stack)-1];stack = stack[:len(stack)-1]
		right := stack[len(stack)-1];stack = stack[:len(stack)-1]
		left := stack[len(stack)-1];stack = stack[:len(stack)-1]

		if right-left <= bush.NodeSize {
			for i := left; i <= right; i++ {
				x = bush.coords[2*i]
				y = bush.coords[2*i+1]
				if x >= minX && x <= maxX && y >= minY && y <= maxY {
					result = append(result, bush.idxs[i])
				}
			}
			continue
		}

		m := floor(float64(left+right) / 2.0)

		x = bush.coords[2*m]
		y = bush.coords[2*m+1]

		if x >= minX && x <= maxX && y >= minY && y <= maxY {
			result = append(result, bush.idxs[m])
		}

		nextAxis := (axis + 1) % 2

		if (axis == 0 && minX <= x) || (axis != 0 && minY <= y) {
			stack = append(stack, left)
			stack = append(stack, m-1)
			stack = append(stack, nextAxis)
		}

		if (axis == 0 && maxX >= x) || (axis != 0 && maxY >= y) {
			stack = append(stack, m+1)
			stack = append(stack, right)
			stack = append(stack, nextAxis)
		}

	}
	return result
}

// Finds all items within a given radius from the query point and returns an array of indices.
func (bush *KDBush)Within(point Point, radius float64) []int {
	stack := []int{0, len(bush.idxs) - 1, 0}
	result := []int{}
	r2 := radius * radius
	qx, qy := point.Coordinates()

	for len(stack) > 0 {
		axis := stack[len(stack) - 1]; stack = stack[:len(stack) - 1]
		right := stack[len(stack) - 1]; stack = stack[:len(stack) - 1]
		left := stack[len(stack) - 1]; stack = stack[:len(stack) - 1]

		if right - left <= bush.NodeSize {
			for i := left; i <= right; i++ {
				dst := sqrtDist(bush.coords[2 * i], bush.coords[2 * i + 1], qx, qy)
				if dst <= r2 {
					result = append(result, bush.idxs[i])
				}
			}
			continue
		}

		m := floor(float64(left + right) / 2.0)
		x := bush.coords[2 * m]
		y := bush.coords[2 * m + 1]

		if sqrtDist(x, y, qx, qy) <= r2 {
			result = append(result, bush.idxs[m])
		}

		nextAxis := (axis + 1) % 2

		if (axis == 0 && (qx - radius <= x)) || (axis != 0 && (qy - radius <= y)) {
			stack = append(stack, left)
			stack = append(stack, m - 1)
			stack = append(stack, nextAxis)
		}

		if (axis == 0 && (qx + radius >= x)) || (axis != 0 && (qy + radius >= y)) {
			stack = append(stack, m + 1)
			stack = append(stack, right)
			stack = append(stack, nextAxis)
		}
	}

	return result
}

///// private method to sort the data

////////////////////////////////////////////////////////////////
/// Sorting stuff
////////////////////////////////////////////////////////////////

func (bush *KDBush) buildIndex(points []Point, nodeSize int) {
	bush.NodeSize = nodeSize
	bush.Points = points

	bush.idxs = make([]int, len(points))
	bush.coords = make([]float64, 2*len(points))

	for i, v := range points {
		bush.idxs[i] = i
		x, y := v.Coordinates()
		bush.coords[i*2] = x
		bush.coords[i*2+1] = y
	}

	sort(bush.idxs, bush.coords, bush.NodeSize, 0, len(bush.idxs)-1, 0)
}

func sort(idxs []int, coords []float64, nodeSize int, left, right, depth int) {
	if (right - left) <= nodeSize {
		return
	}

	m := floor(float64(left+right) / 2.0)

	sselect(idxs, coords, m, left, right, depth%2)

	sort(idxs, coords, nodeSize, left, m-1, depth+1)
	sort(idxs, coords, nodeSize, m+1, right, depth+1)

}

func sselect(idxs []int, coords []float64, k, left, right, inc int) {
	//whatever you want
	for right > left {
		if (right - left) > 600 {
			n := right - left + 1
			m := k - left + 1
			z := math.Log(float64(n))
			s := 0.5 * math.Exp(2.0*z/3.0)
			sds := 1.0
			if float64(m)-float64(n)/2.0 < 0 {
				sds = -1.0
			}
			n_s := float64(n) - s
			sd := 0.5 * math.Sqrt(z*s*n_s/float64(n)) * sds
			newLeft := iMax(left, floor(float64(k)-float64(m)*s/float64(n)+sd))
			newRight := iMax(right, floor(float64(k)+float64(n-m)*s/float64(n)+sd))
			sselect(idxs, coords, k, newLeft, newRight, inc)
		}

		t := coords[2*k+inc]
		i := left
		j := right

		swapItem(idxs, coords, left, k)
		if coords[2*right+inc] > t {
			swapItem(idxs, coords, left, right)
		}

		for i < j {
			swapItem(idxs, coords, i, j)
			i += 1
			j -= 1
			for coords[2*i+inc] < t {
				i += 1
			}
			for coords[2*j+inc] > t {
				j -= 1
			}
		}

		if coords[2*left+inc] == t {
			swapItem(idxs, coords, left, j)
		} else {
			j += 1
			swapItem(idxs, coords, j, right)
		}

		if j <= k {
			left = j + 1
		}
		if k <= j {
			right = j - 1
		}
	}
}

func swapItem(idxs []int, coords []float64, i, j int) {
	swapi(idxs, i, j)
	swapf(coords, 2*i, 2*j)
	swapf(coords, 2*i+1, 2*j+1)
}

func swapf(a []float64, i, j int) {
	t := a[i]
	a[i] = a[j]
	a[j] = t
}

func swapi(a []int, i, j int) {
	t := a[i]
	a[i] = a[j]
	a[j] = t
}

func iMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func floor(in float64) int {
	out := math.Floor(in)
	return int(out)
}


func sqrtDist(ax, ay, bx, by float64) float64{
	dx := ax - bx
	dy := ay - by
	return dx*dx + dy*dy
}