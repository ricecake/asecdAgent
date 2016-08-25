// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	//"github.com/spf13/viper"
	"github.com/gorilla/websocket"
	"gopkg.in/olebedev/go-duktape.v2"
	"log"
	//"github.com/spf13/pflag"
	//"github.com/ugorji/go/codec"
	//"crypto/sha1"
)

var addr string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		//addr := viper.GetString("addr")
		fmt.Println("server called")
		db, err := bolt.Open("asecd.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		db.Close()
		ctx := duktape.New()
		ctx.EvalString(`2 + 3`)
		result := ctx.GetNumber(-1)
		ctx.Pop()
		fmt.Println("result is:", result)
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
		log.Printf("Address: %s", addr)
		log.Printf("connecting to %s", u.String())

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		done := make(chan struct{})

		go func() {
			defer c.Close()
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recv: %s", message)
				err2 := c.WriteMessage(websocket.TextMessage, message)
				if err2 != nil {
					log.Println("write:", err2)
					return
				}
			}
		}()

                err2 := c.WriteMessage(websocket.TextMessage, []byte("Test!"))
                if err2 != nil {
                        log.Println("write:", err)
                        return
                }

		for {
			select {
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().StringVarP(&addr, "addr", "a", "localhost:8080", "service address")

}
