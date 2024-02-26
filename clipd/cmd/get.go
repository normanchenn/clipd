/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "gets the clipboard history",
	Long:  "get the clipboard history with the following flags: --at, --last, --from, --to",
	Run: func(cmd *cobra.Command, args []string) {
		data := "get"
		at, _ := cmd.Flags().GetInt("at")
		last, _ := cmd.Flags().GetInt("last")
		from, _ := cmd.Flags().GetInt("from")
		to, _ := cmd.Flags().GetInt("to")

		if at == -1 && last == -1 && from == -1 && to == -1 {
			data += fmt.Sprintf(" at=%d", 0)
		} else if at != -1 {
			data += fmt.Sprintf(" at=%d", at)
		} else if last != -1 {
			data += fmt.Sprintf(" last=%d", last)
		} else if from != -1 && to != -1 {
			data += fmt.Sprintf(" from=%d to=%d", from, to)
		} else {
			fmt.Println("flags not valid")
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

func sendCmd(command string) {
	conn, err := net.Dial("unix", "/tmp/clipd.sock")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	command += "\n"
	data := []byte(command)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(response[:n]))
}
