package main

import (
    "fmt"
    "bytes"
    "strconv"
)



//  SSRect structure describes a rectangle by the Left,
//  Right, Top, and Bottom coordinates. Additional methods
//  are provided that assist with the SSGrid logic directly.
type SSRect struct {
    L    int
    T    int
    R    int
    B    int
    midx int
    midy int
}

//  Midx returns the x-coordinate of the SSRect's midpoint
//  It also stores this information for later use to avoid
//  having to calculate this repeatedly.
func (rect *SSRect) Midx() int {
    if rect.midx == 0 {
        rect.midx = (rect.L + rect.R) / 2
    }
    return rect.midx
}

//  Midy returns the y-coordinate of the SSRect's midpoint
//  It also stores this information for later use to avoid
//  having to calculate this repeatedly.
func (rect *SSRect) Midy() int {
    if rect.midy == 0 {
        rect.midy = (rect.T + rect.B) / 2
    }
    return rect.midy
}

//  ContainsPoint determines if the point given by (x, y)
//  is located in or on this SSRect
func (rect *SSRect) ContainsPoint(x, y int) bool {
    return rect.L <= x && x <= rect.R && rect.T <= y && y <= rect.B
}

//  ContainsRect determines if the other SSRect is fully contained
//  within this one.
func (rect *SSRect) ContainsRect(other *SSRect) bool {
    return rect.L <= other.L && rect.R >= other.R && rect.T <= other.T && rect.B >= other.B
}

//  DisjointRect determines if the other SSRect does not collide
//  with this one.
func (rect *SSRect) DisjointFrom(other *SSRect) bool {
    return rect.T > other.B || other.T > rect.B || rect.R < other.L || other.R < rect.L
}




//  SSGrid Structure encapsulates the 2D Grid challenge.
//  The main features are the Update and Subsum methods.
//  For all running times provided, n is whichever width
//  or height is greater.
type SSGrid struct {
    root *ssquad
}

//  Reads a value in the grid in O(log(n)) time
func (grid *SSGrid) Read(x, y int) int {
    if !grid.root.bounds.ContainsPoint(x, y) {
        return 0
    }
    return grid.root.read(x, y)
}

//  Update will update the grid location (x, y) with num.
//  This operation will take O(log(n)) time.
func (grid *SSGrid) Update(x, y, num int) {
    if !grid.root.bounds.ContainsPoint(x, y) {
        fmt.Println("(", x, ",", y, ") is not in the grid.")
    } else {
        grid.root.update(x, y, num)
    }
}

//  Subsum will calculate the sum of all cells contained
//  within the subgrid described by rect.
//  This operation will take span O(log(n)) time.
func (grid *SSGrid) Subsum(rect *SSRect) int {
    return grid.root.subsum(rect)
}

//  Provides the string representing the subgrid described by rect.
func (grid *SSGrid) DisplayRect(rect *SSRect) string {
    var buffer bytes.Buffer
    for r := rect.L; r <= rect.R; r++ {
        for c := rect.T; c <= rect.B; c++ {
            num := grid.Read(r, c)
            if num < 100 {
                buffer.WriteString(" ")
                if num < 10 {
                    buffer.WriteString(" ")
                }
            }
            buffer.WriteString(strconv.Itoa(num))
        }
        buffer.WriteString("\n")
    }
    return buffer.String()
}

//  Calls DisplayRect for the entire grid.
func (grid *SSGrid) String() string {
    return grid.DisplayRect(grid.root.bounds)
}


//  NewSSGrid will create a new SSGrid with dimensions m x n
func NewSSGrid(m, n int) *SSGrid {
    return &SSGrid{root:&ssquad{bounds:&SSRect{L: 0, T: 0, R: m - 1, B: n - 1}}}
}


//  INTERNAL


//  SSQuad is an internal structure. It is a QuadTree which also contains
//  information about which subgrid it is responsible for as well as the sum
//  of the entire subgrid
type ssquad struct {
    ul      *ssquad
    ur      *ssquad
    bl      *ssquad
    br      *ssquad
    bounds  *SSRect
    data    int
}

func (quad *ssquad) read(x, y int) int{
    if quad.bounds.L == x && quad.bounds.T == y && quad.bounds.L == quad.bounds.R && quad.bounds.T == quad.bounds.B {
       return quad.data
    }
    if quad.ul != nil && quad.ul.bounds.ContainsPoint(x, y) {
        return quad.ul.read(x, y)
    }
    if quad.ur != nil && quad.ur.bounds.ContainsPoint(x, y) {
        return quad.ur.read(x, y)
    }
    if quad.bl != nil && quad.bl.bounds.ContainsPoint(x, y) {
        return quad.bl.read(x, y)
    }
    if quad.br != nil && quad.br.bounds.ContainsPoint(x, y) {
        return quad.br.read(x, y)
    }
    return 0
}

func left(quad *ssquad) int {
    left := quad.bounds.Midx() + 1
    if (left > quad.bounds.R) {
        return quad.bounds.R
    }
    return left
}
func top(quad *ssquad) int {
    top := quad.bounds.Midy() + 1
    if (top > quad.bounds.B) {
        return quad.bounds.B
    }
    return top
}

func (quad *ssquad) updateSum() {
    sum := 0
    if quad.ul != nil {
        sum += quad.ul.data
    }
    if quad.ur != nil {
        sum += quad.ur.data
    }
    if quad.bl != nil {
        sum += quad.bl.data
    }
    if quad.br != nil {
        sum += quad.br.data
    }
    quad.data = sum
}

func (quad *ssquad) update(x, y, num int) {
    if quad.bounds.L == quad.bounds.R && quad.bounds.T == quad.bounds.B {
        if quad.bounds.L != x || quad.bounds.T != y {
            fmt.Println("Error", x, ",", y, "is not", quad.bounds.L, ",", quad.bounds.T)
            return
        }
        quad.data = num
        return
    }
    if quad.ul == nil {
        quad.ul = &ssquad{bounds:&SSRect{L:quad.bounds.L, T:quad.bounds.T, R:quad.bounds.Midx(), B:quad.bounds.Midy()}}
    }
    if quad.ul.bounds.ContainsPoint(x, y) {
        quad.ul.update(x, y, num)
        quad.updateSum()
        return
    }
    if quad.ur == nil {
        quad.ur = &ssquad{bounds:&SSRect{L:quad.bounds.L, T:top(quad), R:quad.bounds.Midx(), B:quad.bounds.B}}
    }
    if quad.ur.bounds.ContainsPoint(x, y) {
        quad.ur.update(x, y, num)
        quad.updateSum()
        return
    }
    if quad.bl == nil {
        quad.bl = &ssquad{bounds:&SSRect{L:left(quad), T:quad.bounds.T, R:quad.bounds.R, B:quad.bounds.Midy()}}
    }
    if quad.bl.bounds.ContainsPoint(x, y) {
        quad.bl.update(x, y, num)
        quad.updateSum()
        return
    }
    if quad.br == nil {
        quad.br = &ssquad{bounds:&SSRect{L:left(quad), T:top(quad), R:quad.bounds.R, B:quad.bounds.B}}
    }
    if quad.br.bounds.ContainsPoint(x, y) {
        quad.br.update(x, y, num)
        quad.updateSum()
        return
    }
}

func (quad *ssquad) subsum(rect *SSRect) int {
    //  case: no intersection = irrelevant
    if rect.DisjointFrom(quad.bounds) {
        return 0
    }
    //  case: fully encapsulated by input
    if rect.ContainsRect(quad.bounds) {
        return quad.data
    }
    //  case: partial encapsulation
    sums := make(chan int, 4)
    sum := 0
    calcSum := func (q *ssquad) {
        if q == nil {
            sums <- 0
        } else {
            sums <- q.subsum(rect)
        }
    }
    go calcSum(quad.ul)
    go calcSum(quad.ur)
    go calcSum(quad.bl)
    go calcSum(quad.br)
    for i := 0; i < 4; i++ {
        sum += <- sums
    }
    return sum
}