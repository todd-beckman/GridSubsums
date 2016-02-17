package main

import (
    "fmt"
    "bytes"
    "strconv"
)

type SSRect struct {
    L    int
    T    int
    R    int
    B    int
    midx int
    midy int
}

func (rect *SSRect) Midx() int {
    if rect.midx == 0 {
        rect.midx = (rect.L + rect.R) / 2
    }
    return rect.midx
}

func (rect *SSRect) Midy() int {
    if rect.midy == 0 {
        rect.midy = (rect.T + rect.B) / 2
    }
    return rect.midy
}

func (rect *SSRect) ContainsRect(other SSRect) bool {
    return rect.L >= other.L && rect.T >= other.T && rect.R >= other.R && rect.B >= other.B
}

func (rect *SSRect) ContainsPoint(x, y int) bool {
    return rect.L <= x && x <= rect.R && rect.T <= y && y <= rect.B
}


type SSGrid struct {
    Root *ssquad
}

type ssquad struct {
    UL      *ssquad
    UR      *ssquad
    BL      *ssquad
    BR      *ssquad
    Bounds  SSRect //bounds
    Data    int
}

func (quad *ssquad) Read(x, y int) int{
    if quad.Bounds.L == x && quad.Bounds.T == y && quad.Bounds.L == quad.Bounds.R && quad.Bounds.T == quad.Bounds.B {
       return quad.Data
    }
    if quad.UL != nil && quad.UL.Bounds.ContainsPoint(x, y) {
        return quad.UL.Read(x, y)
    }
    if quad.UR != nil && quad.UR.Bounds.ContainsPoint(x, y) {
        return quad.UR.Read(x, y)
    }
    if quad.BL != nil && quad.BL.Bounds.ContainsPoint(x, y) {
        return quad.BL.Read(x, y)
    }
    if quad.BR != nil && quad.BR.Bounds.ContainsPoint(x, y) {
        return quad.BR.Read(x, y)
    }
    return 0
}

func left(quad *ssquad) int {
    left := quad.Bounds.Midx() + 1
    if (left > quad.Bounds.R) {
        return quad.Bounds.R
    }
    return left
}
func top(quad *ssquad) int {
    top := quad.Bounds.Midy() + 1
    if (top > quad.Bounds.B) {
        return quad.Bounds.B
    }
    return top
}

func (quad *ssquad) updateSum() {
    sum := 0
    if quad.UL != nil {
        sum += quad.UL.Data
    }
    if quad.UR != nil {
        sum += quad.UR.Data
    }
    if quad.BL != nil {
        sum += quad.BL.Data
    }
    if quad.BR != nil {
        sum += quad.BR.Data
    }
    quad.Data = sum
}

func (quad *ssquad) Update(x, y, num int) {
    if quad.Bounds.L == quad.Bounds.R && quad.Bounds.T == quad.Bounds.B {
        if quad.Bounds.L != x || quad.Bounds.T != y {
            fmt.Println("Error", x, ",", y, "is not", quad.Bounds.L, ",", quad.Bounds.T)
            return
        }
        quad.Data = num
        return
    }
    if quad.UL == nil {
        quad.UL = &ssquad{Bounds:SSRect{L:quad.Bounds.L, T:quad.Bounds.T, R:quad.Bounds.Midx(), B:quad.Bounds.Midy()}}
    }
    if quad.UL.Bounds.ContainsPoint(x, y) {
        quad.UL.Update(x, y, num)
        quad.updateSum()
        return
    }
    if quad.UR == nil {
        quad.UR = &ssquad{Bounds:SSRect{L:quad.Bounds.L, T:top(quad), R:quad.Bounds.Midx(), B:quad.Bounds.B}}
    }
    if quad.UR.Bounds.ContainsPoint(x, y) {
        quad.UR.Update(x, y, num)
        quad.updateSum()
        return
    }
    if quad.BL == nil {
        quad.BL = &ssquad{Bounds:SSRect{L:left(quad), T:quad.Bounds.T, R:quad.Bounds.R, B:quad.Bounds.Midy()}}
    }
    if quad.BL.Bounds.ContainsPoint(x, y) {
        quad.BL.Update(x, y, num)
        quad.updateSum()
        return
    }
    if quad.BR == nil {
        quad.BR = &ssquad{Bounds:SSRect{L:left(quad), T:top(quad), R:quad.Bounds.R, B:quad.Bounds.B}}
    }
    if quad.BR.Bounds.ContainsPoint(x, y) {
        quad.BR.Update(x, y, num)
        quad.updateSum()
        return
    }
}

func (quad *ssquad) Subsum(rect SSRect) int{
    //  Case: not part of the sum
    if !rect.ContainsRect(quad.Bounds) {
        return 0
    }
    //  Case: single cell
    if quad.UL == nil {
        return quad.Data
    }
    sums := make(chan int, 4)
    calcSum := func (q *ssquad) {
        if q == nil {
            sums <- 0
        } else {
            sums <- q.Subsum(rect)
        }
    }
    go calcSum(quad.UL)
    go calcSum(quad.UR)
    go calcSum(quad.BL)
    go calcSum(quad.BR)
    sum := 0
    for i := 0; i < 4; i++ {
        sum += <- sums
    }
    return sum
}


func NewSSGrid(m, n int) *SSGrid {
    return &SSGrid{Root:&ssquad{Bounds:SSRect{L: 0, T: 0, R: m - 1, B: n - 1}}}
}

func (grid *SSGrid) Read(x, y int) int {
    if !grid.Root.Bounds.ContainsPoint(x, y) {
        return 0
    }
    return grid.Root.Read(x, y)
}

func (grid *SSGrid) Update(x, y, num int) {
    if !grid.Root.Bounds.ContainsPoint(x, y) {
        fmt.Println("(", x, ",", y, ") is not in the grid.")
    } else {
        grid.Root.Update(x, y, num)
    }
}

func (grid *SSGrid) Subsum(rect SSRect) int {
    return grid.Root.Subsum(rect)
}

func (grid *SSGrid) DisplayRect(rect SSRect) string {
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

func (grid *SSGrid) String() string {
    return grid.DisplayRect(grid.Root.Bounds)
}