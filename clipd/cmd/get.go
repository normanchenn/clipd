/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

const (
	SOCKETPATH = "/tmp/clipd.sock"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "gets the clipboard history",
	Long:  "get the clipboard history with the following flags: --at, --last, --from, --to",
	Run: func(cmd *cobra.Command, args []string) {
		params := make(map[string]int)
		flags := []string{"at", "last", "from", "to"}

		for _, val := range flags {
			addParam(cmd, val, params)
		}

		data := map[string]interface{}{
			"action": "get",
			"params": params,
		}
		sendCmd(data)

	},
}

func init() {
	getCmd.Flags().IntP("at", "a", -1, "get the clipboard history at a specific index")
	getCmd.Flags().IntP("last", "l", -1, "get the last n items clipboard history")
	getCmd.Flags().IntP("from", "f", -1, "get the clipboard history from a specific index - must have --to as well")
	getCmd.Flags().IntP("to", "t", -1, "get the clipboard history to a specific index - must be --from as well")
}

func addParam(cmd *cobra.Command, flag string, params map[string]int) {
	value, _ := cmd.Flags().GetInt(flag)
	if value != -1 {
		params[flag] = value
	}
}

func sendCmd(data map[string]interface{}) {
	conn, err := net.Dial("unix", SOCKETPATH)
	defer conn.Close()
	if err != nil {
		fmt.Println("Error connecting to socket: ", err)
		return
	}

	encoder := json.NewEncoder(conn)
	err = encoder.Encode(data)
	if err != nil {
		fmt.Println("Error encoding json: ", err)
		return
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println("Error reading response: ", err)
		return
	}
	fmt.Println(string(response[:n]))
}
