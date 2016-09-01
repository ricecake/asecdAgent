// Copyright © 2016 Sebastian Green-Husted <ricecake@tfm.nu>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/websocket"
	"gopkg.in/olebedev/go-duktape.v2"
	"log"
	//"github.com/ugorji/go/codec"
	//"crypto/sha1"
)

var addr string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the daemon in server mode",
	Long: `Runs the asecdClient as a script execution daemon.
Will connect to the remote job control server, and wait to be passed work to exectute`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("asecd.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
		log.Printf("connecting to %s", u.String())

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		done := make(chan struct{})
		messageChannel := make(chan []byte)

		go func() {
			defer c.Close()
			defer close(done)
			defer close(messageChannel)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				messageChannel <- message
				err2 := c.WriteMessage(websocket.TextMessage, message)
				if err2 != nil {
					log.Println("write:", err2)
					return
				}
			}
		}()

		for {
			select {
			case message := <- messageChannel:
				go func(){
					// Pass a communication channel to the coroutine that it can use to request resources
					// Of the write thread by passing it's own channel along for the return
					Comm := make(chan struct{})
					defer close(Comm)

					ctx := duktape.New()
					defer ctx.DestroyHeap()
					ctx.EvalString(`2 + 3`)
					result := ctx.GetNumber(-1)
					ctx.Pop()
					fmt.Println("result is:", result)
					log.Printf("Recv: %s", message)
				}()
			case <-interrupt:
				log.Println("interrupt")
				// To cleanly close a connection, a client should send a close
				// frame and wait for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
					case <-done:
					case <-time.After(time.Second):
				}
				c.Close()
				return
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&addr, "addr", "a", "localhost:8080", "service address")

}
