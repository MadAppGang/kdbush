// MIT License
//
// Copyright (c) 2016 MadAppGang

// Package kdbush provides a very fast static spatial index for 2D points based on a flat KD-tree.
//
// Very fast, but limited. Here are main limitations:
//
// 1. Points only, no rectangles
//
// 2. 2 dimensional
//
// 3. indexing 16-40 times faster then  rtreego(https://github.com/dhconnelly/rtreego) (TODO: benchmark)
//
// 4. Implements radius search  (rtreego and go.geo only have range search)
//
// There are three amazing other options for geospatial indexing:
//   tile38 - http://tile38.com
//   go.geo - https://github.com/paulmach/go.geo/tree/master/quadtree
//   rtreego - https://github.com/dhconnelly/rtreego
//
// All this modules are dynamic and complex.
//
//
// This implementation is based on:
//
// JS library: https://github.com/mourner/kdbush
//
// C++11 port: https://github.com/mourner/kdbush.hpp
//
// If you liked the project, start it please: https://github.com/MadAppGang/kdbush
//
package kdbush
