package main

import "fmt"
import "log"
import "os"
import "github.com/spf13/pflag"
import "gopkg.in/olebedev/go-duktape.v2"
import "github.com/boltdb/bolt"
import "github.com/ricecake/asecdAgent/cmd"
//import "github.com/gorilla/websocket"
//import "github.com/ugorji/go/codec"
//import "crypto/sha1"

func main() {
    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
  ctx := duktape.New()
  ctx.EvalString(`2 + 3`)
  result := ctx.GetNumber(-1)
  ctx.Pop()
  fmt.Println("result is:", result)

    if err := cmd.RootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }
}
