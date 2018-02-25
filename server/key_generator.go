package main

import (
        "math/rand"
        "time"
        )

func init() {
    rand.Seed(time.Now().UnixNano())
}


var alpha_numeric_runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
func alphaNum(n int) string {
    return stringGenerator(n, alpha_numeric_runes)
}

func stringGenerator(n int, runes []rune) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = runes[rand.Intn(len(runes))]
    }
    return string(b)
}
