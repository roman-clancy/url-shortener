package main

import (
    "fmt"
    "url_shortener/internal/config"
)

func main() {
    cfg := config.MustLoad()
    fmt.Println(cfg)
}
