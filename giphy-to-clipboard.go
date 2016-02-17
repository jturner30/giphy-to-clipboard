package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func downloadFromFile(url string) (string, error) {
	out, err := os.Create("/tmp/giphy.gif")
	if err != nil {
		fmt.Println("Error while creating tmp file")
		return "", err
	}
	defer out.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading file from ", url)
		return "", err
	}

	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Println("Error while downloading file from ", url)
		return "", err
	}
	return out.Name(), nil
}

func copyToClipboard() error {
	c := exec.Command("osascript", "-e", "tell application \"Finder\" to set the clipboard to ( POSIX file \"/tmp/giphy.gif\"  ) ")
	err := c.Run()
	if err != nil {
		fmt.Println("Error copying file to clipboard. Got error", err)
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: giphy-to-clipboard <giphy url>")
		return
	}
	url := args[0]
	_, err := downloadFromFile(url)
	if err != nil {
		fmt.Println("Error copying url", url, "got error", err)
	}
	fmt.Println("Copying", url, " to clipboard")
	copyToClipboard()
}
