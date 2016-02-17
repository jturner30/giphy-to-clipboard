package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// Download a file using a GET request
// Based almost entirely on
// https://github.com/thbar/golang-playground/blob/master/download-files.go
// TODO: Is this really the way to handle errors?
// TODO: That error message isn't particularly helpful
func downloadFromFile(url string) (string, error) {
	out, err := os.Create("/tmp/giphy.gif")
	if err != nil {
		return "", err
	}
	defer out.Close()

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if response.Status != "200 OK" {
		return "", errors.New("Could not find resource")
	}

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return "", err
	}
	return out.Name(), nil
}

// Runs an AppleScript in lieu of having a decent clipboard API to work with
// Applescript courtesy of
// https://apple.stackexchange.com/questions/66459/copying-files-to-the-clipboard-using-applescript
func copyToClipboard() error {
	c := exec.Command(
		"osascript",
		"-e",
		"tell application \"Finder\" to set the clipboard to ( POSIX file \"/tmp/giphy.gif\"  ) ",
	)
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

// Usage: giphy-to-clipboard url
// url should be a resource at i.giphy.com, though theoretically this program
// can copy an arbitrary resource into a file named giphy.gif. This is probably
// insecure and most likely should be avoided
func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		color.Blue("Usage: giphy-to-clipboard <giphy url>")
		return
	}
	url := args[0]
	_, err := downloadFromFile(url)
	if err != nil {
		color.Red(fmt.Sprintf("Error copying %s", url))
		return
	}
	fmt.Println("Copying", url, " to clipboard")
	err = copyToClipboard()
	if err != nil {
		errString := fmt.Sprintf("Copying %s, to clipboard failed. Encountered error: %s", url, err)
		color.Red(errString)
		return
	}
	color.Green("Your GIF has been copied!")
}
