package main

import (
    "fmt"
)

func gen_m_num (m int, num *[]int) {
    for i:=1; i <= m; i++ {
        *num = append(*num, i)
    }
}

func gen_n_combin (n int, pre []int, input []int) {
    if n == 1 {
        for _, v := range input {
            fmt.Println(append(pre, v))
        }
    } else if len(input) > n {
        for i := 1; i <= len(input) - (n-1); i++ {
            new_pre := append(pre, input[i-1])
            gen_n_combin(n-1, new_pre, input[i:])
        }
    } else if len(input) == n {
        fmt.Println(append(pre, input...))
    }
}

func main() {
    var m int
    var n int
    fmt.Scanln(&m, &n)
    list := make([]int, 0)
    gen_m_num(m, &list)
    fmt.Println("gen num:", list)
    gen_n_combin(n, make([]int, 0), list)
}

