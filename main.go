package main

import (
    "github.com/mukezhz/learn-huma/bootstrap"

    "github.com/joho/godotenv"
)

func main() {
    _ = godotenv.Load()
    _ = bootstrap.RootApp.Execute()
}
