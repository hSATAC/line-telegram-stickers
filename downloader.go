package main

import (
	"archive/zip"
	"fmt"
	"github.com/disintegration/imaging"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	file, err := ioutil.TempFile(os.TempDir(), "line-telegram-stickers-pack")
	if err != nil {
		fmt.Println("Error while creating", file, "-", err)
		return
	}
	defer os.Remove(file.Name())
	fmt.Println("Temp File created!")

	downloadLinePack(3333, file)

	dir, err := ioutil.TempDir(os.TempDir(), "line-telegram-stickers-pack")
	if err != nil {
		fmt.Println("Error while creating", dir, "-", err)
		return
	}

	fmt.Println("Unzip to", dir)
	unzip(file.Name(), dir)

	//dir := "/var/folders/q7/bcn45hls729bj9dncgfyv5d80000gn/T/line-telegram-stickers-pack612988553"
	filterFileInDir(dir)
	resizeFileInDir(dir)
}

func resizeFileInDir(dir string) {
	filepath.Walk(dir, func(path string, _ os.FileInfo, _ error) error {
		if path == dir {
			return nil
		}

		if match, _ := regexp.MatchString("[0-9]+@2x.png", path); match {
			src, err := imaging.Open(path)
			if err != nil {
				log.Fatalf("Open failed: %v", err)
			}
			b := src.Bounds()
			imgWidth := b.Max.X
			imgHeight := b.Max.Y

			if imgWidth == 512 && imgHeight <= 512 {
				return nil
			}

			if imgWidth <= 512 && imgHeight == 512 {
				return nil
			}

			if imgWidth > imgHeight {
				src = imaging.Resize(src, 512, 0, imaging.Lanczos)
			} else {
				src = imaging.Resize(src, 0, 512, imaging.Lanczos)
			}

			err = imaging.Save(src, path)
			if err != nil {
				log.Fatalf("Save failed: %v", err)
			}
			fmt.Println("Resized", path)
		}
		return nil
	})
}

func filterFileInDir(dir string) {
	filepath.Walk(dir, func(path string, _ os.FileInfo, _ error) error {
		if path == dir {
			return nil
		}

		if match, _ := regexp.MatchString("[0-9]+@2x.png", path); !match {
			fmt.Println("Delete file", path)
			os.Remove(path)
		}
		return nil
	})
}

func downloadLinePack(id int, file *os.File) (string, error) {
	url := fmt.Sprintf("http://dl.stickershop.line.naver.jp/products/0/0/1/%d/iphone/stickers@2x.zip", id)
	return downloadFromUrl(url, file)
}

func downloadFromUrl(url string, file *os.File) (string, error) {
	fmt.Println("Downloading", url, "to", file.Name())

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}
	defer response.Body.Close()

	n, err := io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}

	fmt.Println(n, "bytes downloaded.")

	return file.Name(), nil
}

func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
