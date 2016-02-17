package main

import (
    "fmt"
    "math/rand"
    "strconv"
)

const (
    w = 10
    h = 10
    ra = 100
    nums = 100
)

func printArr(arr [][h]int) {
    for i := 0; i < w; i++ {
        for j := 0; j < h; j++ {
            num := arr[i][j]
            if num < 100 {
                fmt.Print(" ")
                if num < 10 {
                    fmt.Print(" ")
                }
            }
            fmt.Print(strconv.Itoa(num))
        }
        fmt.Println()
    }
}

func main() {
    dummy := [w][h]int{}

    grid := NewSSGrid(w, h)

    for i := 0; i < nums; i++ {
        x, y, num := rand.Intn(w), rand.Intn(h), rand.Intn(ra)
        dummy[x][y] = num
        grid.Update(x, y, num)
    }
    fmt.Println("Dummy array:")
    printArr(dummy[:])
    fmt.Println()
    fmt.Println("SSGrid:")
    fmt.Println(grid.DisplayRect(SSRect{L: 0, T: 0, R: w - 1, B: h - 1}))
    fmt.Println()
    rect := SSRect{L: 1, T: 2, R: 4, B: 3}
    fmt.Println("Sample rect from:", rect)
    fmt.Println(grid.DisplayRect(rect))
    fmt.Println(grid.Subsum(SSRect{L: 1, T: 2, R: 4, B: 3}))
}

