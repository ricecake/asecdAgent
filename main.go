package main

import (
    "os"
    "fmt"
    "github.com/ricecake/asecdAgent/cmd"
)

func main() {
    if err := cmd.RootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }
}
