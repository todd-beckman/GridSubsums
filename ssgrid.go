package main

import (
    "fmt"
    "bytes"
    "strconv"
)

func Unused() {
    fmt.Print("Just for fmt import during debugging")
}

type SSRect struct {
    L int
    T int
    R int
    B int
    midx int
    midy int
}

func (rect *SSRect) Midx() int {
    if rect.midx == 0 {
        rect.midx = rect.L + rect.R / 2
    }
    return rect.midx
}

func (rect *SSRect) Midy() int {
    if rect.midy == 0 {
        rect.midy = rect.T + rect.B / 2
    }
    return rect.midy
}

func (rect *SSRect) Contains(other SSRect) bool {
    return rect.L >= other.L && rect.T >= other.T && rect.R >= other.R && rect.B >= other.B
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

func (quad *ssquad) Read(x, y int) int {
    if quad.Bounds.L == quad.Bounds.R && quad.Bounds.T == quad.Bounds.B {
       return quad.Data
    }
    if quad.Bounds.Midx() <= x {
        if quad.Bounds.Midy() <= y {
            if quad.UL == nil {
                return 0
            } else {
                return quad.UL.Read(x, y)
            }
        } else {
            if quad.BL == nil {
                return 0
            } else {
                return quad.BL.Read(x, y)
            }
        }
    } else {
        if quad.Bounds.Midy() <= y {
            if quad.UR == nil {
                return 0
            } else {
                return quad.UR.Read(x, y)
            }
        } else {
            if quad.BR == nil {
                return 0
            } else {
                return quad.BR.Read(x, y)
            }
        }
    }
}

func (quad *ssquad) UpdateSum() {
    if quad.UL != nil && (quad.Bounds.L != quad.Bounds.R || quad.Bounds.T != quad.Bounds.B) {
        quad.Data = quad.UL.Data
        if quad.UR != nil {
            quad.Data += quad.UR.Data
        }
        if quad.BL != nil {
            quad.Data += quad.BL.Data
            if quad.BR != nil {
                quad.Data += quad.BR.Data
            }
        }
    }
}

func (quad *ssquad) Update(x, y, num int) {

    //  base case: 1 cell
    if quad.Bounds.L == quad.Bounds.R && quad.Bounds.T == quad.Bounds.B {
        if quad.Bounds.L != x || quad.Bounds.T != y {
            fmt.Println("Error", x, ",", y, "is not", quad.Bounds.L, ",", quad.Bounds.T)
            return
        }
        fmt.Println("Posting", num, "at", x, ",", y)
        quad.Data = num
        return
    }
    if quad.Bounds.Midx() <= x {
        if quad.Bounds.Midy() <= y {
            //  case: top left
            if quad.UL == nil {
                quad.UL = &ssquad{Bounds:SSRect{L:quad.Bounds.L, T:quad.Bounds.T, R:quad.Bounds.Midx(), B:quad.Bounds.Midy()}}
            }
            quad.UL.Update(x, y, num)
        } else {
            //  case: bottom left
            if quad.BL == nil {
                top := quad.Bounds.Midy() + 1
                if (top > quad.Bounds.B) {
                    top = quad.Bounds.B
                }
                quad.BL = &ssquad{Bounds:SSRect{L:quad.Bounds.L, T:top, R:quad.Bounds.Midx(), B:quad.Bounds.B}}
            }
            quad.BL.Update(x, y, num)
        }
    } else {
        if quad.Bounds.Midy() <= y {
            //  case: top right
            if quad.UR == nil {
                left := quad.Bounds.Midx() + 1
                if (left > quad.Bounds.R) {
                    left = quad.Bounds.R
                }
                quad.UR = &ssquad{Bounds:SSRect{L:left, T:quad.Bounds.T, R:quad.Bounds.R, B:quad.Bounds.Midy()}}
            }
            quad.UR.Update(x, y, num)
        } else {
            //  case: bottom right
            if quad.BR == nil {
                left := quad.Bounds.Midx() + 1
                if (left > quad.Bounds.R) {
                    left = quad.Bounds.R
                }
                top := quad.Bounds.Midy() + 1
                if (top > quad.Bounds.B) {
                    top = quad.Bounds.B
                }
                quad.BR = &ssquad{Bounds:SSRect{L:left, T:top, R:quad.Bounds.R, B:quad.Bounds.B}}
            }
            quad.BR.Update(x, y, num)
        }
    }
    quad.UpdateSum()
}

func (quad *ssquad) Subsum(rect SSRect) int{
    //  Case: not part of the sum
    if !rect.Contains(quad.Bounds) {
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
    return grid.Root.Read(x, y)
}

func (grid *SSGrid) Update(x, y, num int) {
    grid.Root.Update(x, y, num)
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