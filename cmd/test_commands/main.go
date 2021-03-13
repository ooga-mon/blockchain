/*
	This was designed to be a quick way to continuously test interacting with the http endpoints from a node.
	It assumes a node is up and running to connect to.
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/octothorped/blockchain/database"
	"github.com/octothorped/blockchain/node"
	"github.com/octothorped/blockchain/wallet"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	url := fmt.Sprintf("http://%s:%d/", node.DefaultIP, node.DefaultPort)
	wal := wallet.NewWallet()
	randomPerson := wallet.NewWallet()

	for scanner.Scan() {
		words := strings.SplitN(scanner.Text(), " ", 2)
		switch strings.ToLower(words[0]) {
		case "send":
			if len(words) < 2 {
				fmt.Println("No word following send was found.")
				continue
			}
			signedTx, _ := wal.CreateTransaction(randomPerson.PublicAddress, words[1])
			_, err := node.WriteRequest(url+"node/transaction", signedTx)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case "get":
			resp, err := http.Get(url + "/node/blocks")
			if err != nil {
				fmt.Println(err)
				continue
			}

			respBlocks := []database.Block{}
			err = node.ReadResponse(resp, &respBlocks)
			if err != nil {
				fmt.Println(err)
				continue
			}

			prettyJsonBlocks, _ := json.MarshalIndent(respBlocks, "", "    ")

			fmt.Println(string(prettyJsonBlocks))

		default:
			fmt.Printf("Not a valid command. Input: %s\n", words[0])
		}
	}

}
