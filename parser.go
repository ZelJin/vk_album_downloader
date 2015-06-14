package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Get Vk user id from and saving location from the flags
	currentPath, _ := os.Getwd()
	userIDPtr := flag.String("uid", "6447964", "User id of your Vkontakte account.")
	pathPtr := flag.String("path", currentPath, "Albums download path")

	flag.Parse()
	fmt.Println("Welcome to VK image downloader!")
	fmt.Println("Parsed values:")
	fmt.Println("uid:", *userIDPtr)
	fmt.Println("path:", *pathPtr)

}
