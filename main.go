package main

import (
    "fmt"
    "math/rand"
)

const (
    w = 1000
    h = 1000
    ra = 99
    nums = 1000000
)

func main() {

    grid := NewSSGrid(w, h)

    for i := 0; i < nums; i++ {
        x, y, num := rand.Intn(w), rand.Intn(h), rand.Intn(ra)
        grid.Update(x, y, num)
    }
    fmt.Println(grid.DisplayRect(&SSRect{L: 0, T: 0, R: w - 1, B: h - 1}))
    fmt.Println()
    rect := &SSRect{L: 1, T: 1, R: 8, B: 8}
    fmt.Println("Sample rect from:", *rect)
    fmt.Println(grid.DisplayRect(rect))
    fmt.Println(grid.Subsum(rect))
}

