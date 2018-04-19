package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hSATAC/line-telegram-stickers/mtproto"
)

func uploadPack(m *mtproto.MTProto, packName string, dir string) {
	stickerPeer, err := m.GetForeignPeer("@Stickers")
	if err != nil {
		fmt.Printf("Could not find @Stickers: %s\n", err)
	}

	// cancal everything anyway.
	m.SendMsgToForeignPeer(stickerPeer, "/cancel")

	// add pack
	m.SendMsgToForeignPeer(stickerPeer, "/newpack")

	// pack name
	m.SendMsgToForeignPeer(stickerPeer, packName)

	// start to send pack
	filepath.Walk(dir, func(path string, _ os.FileInfo, _ error) error {
		if path == dir {
			return nil
		}
		// upload image
		fmt.Println("Uploading: ", path)
		err = m.SendMediaDocumentToForeignPeer(stickerPeer, path)
		if err != nil {
			fmt.Println("Upload failed:", err)
			return err
		}
		// choose emoji
		m.SendMsgToForeignPeer(stickerPeer, "🐈")
		return nil
	})

	// publish
	m.SendMsgToForeignPeer(stickerPeer, "/publish")

	// choose a short name
	m.SendMsgToForeignPeer(stickerPeer, packName)

	// done.
	fmt.Println("=======================================================================")
	fmt.Printf("Your pack should be published at https://t.me/addstickers/%s\n\n", packName)
	fmt.Println("Please check your telegram conversation with @Stickers for more detail.")
	fmt.Println("=======================================================================")

}
