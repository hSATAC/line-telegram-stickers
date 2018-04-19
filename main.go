package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hSATAC/line-telegram-stickers/mtproto"
)

func usage() {
	fmt.Print("Move Line stickers to Telegram.\n\nUsage:\n\n")
	fmt.Print("    line-telegram-stickers <command> [arguments]\n\n")
	fmt.Print("The commands are:\n\n")
	fmt.Print("    auth  <phone_number>            auth connection by code\n")
	fmt.Print("    move	<line_pack_id>             move Line sticker pack to telegram\n")
	fmt.Println()
}

var commands = map[string]int{
	"auth": 1,
	"move": 1,
}

func main() {
	var err error

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	valid := false
	for k, v := range commands {
		if os.Args[1] == k {
			if len(os.Args) < v+2 {
				usage()
				os.Exit(1)
			}
			valid = true
			break
		}
	}

	if !valid {
		usage()
		os.Exit(1)
	}

	m, err := mtproto.NewMTProto(os.Getenv("HOME") + "/.line-telegram-stickers")
	if err != nil {
		fmt.Printf("Create failed: %s\n", err)
		os.Exit(2)
	}

	err = m.Connect()
	if err != nil {
		fmt.Printf("Connect failed: %s\n", err)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "auth":
		err = m.Auth(os.Args[2])
	case "move":
		packID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		// Try connection before all operations
		_, err = m.GetForeignPeer("@Stickers")
		if err != nil {
			if err.Error() == "RPC: mtproto.TL_rpc_error{error_code:401, error_message:\"AUTH_KEY_UNREGISTERED\"}" {
				fmt.Println("Please run `line-telegram-stickers auth` before move stickers.")
			} else {
				fmt.Printf("Connection to telegram error: %s\n", err)
			}
			os.Exit(2)
		}
		packName := fmt.Sprintf("line_telegram_stickers_%d", packID)
		// download pack
		dir := downloadAndFilterPack(packID)
		defer os.RemoveAll(dir)
		// upload pack
		uploadPack(m, packName, dir)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
