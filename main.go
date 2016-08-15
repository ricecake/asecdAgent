package main

import "fmt"
import "log"
import "gopkg.in/olebedev/go-duktape.v2"
import "github.com/boltdb/bolt"
import "github.com/gorilla/websocket"
import "github.com/ugorji/go/codec"

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
}
