package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

//Album Represents the album structure parsed from Vkontakte.
type Album struct {
	id    string
	title string
}

//Photo Represents the photo structure parsed from Vkontakte.
type Photo struct {
	id   string
	link string
}

func main() {
	// Get Vk user id from and save location from the flags
	currentPath, _ := os.Getwd()
	userIDPtr := flag.String("uid", "6447964", "User id of your Vkontakte account.")
	pathPtr := flag.String("path", currentPath, "Albums download path")

	flag.Parse()
	fmt.Println("Welcome to VK image downloader!")
	fmt.Println("Parsed values:")
	fmt.Println("uid:", *userIDPtr)
	fmt.Println("path:", *pathPtr)

	//parseAlbums(*userIDPtr)

	// Parse all albums
	albums := parseAlbums(*userIDPtr)
	for _, album := range albums {
		// Create a folder for that album
		//os.Mkdir(album.title, 0777)
		photos := parsePhotos(*userIDPtr, album.id)
		for _, photo := range photos {
			// Save photo by url to the destination folder
			//downloadPhoto(path string, name string, url string)
		}

	}

}

func parseAlbums(uid string) []Album {
	// Initialize an empty []Album object
	var albums []Album

	// Perform a request to the API
	url := fmt.Sprintf("https://api.vk.com/method/photos.getAlbums.json?user_id=%v&v=5.34", uid)
	response, _ := http.Get(url)

	// Reading the request contents
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	// Parsing JSON
	var jsonResponse map[string]interface{}
	_ = json.Unmarshal(body, &jsonResponse)

	// Getting to albums
	jsonAlbums := jsonResponse["response"].(map[string]interface{})["items"].([]interface{})
	for _, jsonAlbum := range jsonAlbums {
		parsedAlbum := jsonAlbum.(map[string]interface{})
		albums = append(albums, Album{strconv.Itoa(int(parsedAlbum["id"].(float64))), parsedAlbum["title"].(string)})
	}

	return albums
}

func parsePhotos(uid string, albumID string) []Photo {
	// Initialize an empty []Album object
	var photos []Photo

	// Perform a request to the API
	url := fmt.Sprintf("https://api.vk.com/method/photos.get.json?owner_id=%v&album_id=%v&v=5.34", uid, albumID)
	fmt.Println(url)
	response, _ := http.Get(url)

	// Reading the request contents
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	// Parsing JSON
	var jsonResponse map[string]interface{}
	_ = json.Unmarshal(body, &jsonResponse)

	// Getting to photos
	jsonPhotos := jsonResponse["response"].(map[string]interface{})["items"].([]interface{})

	for _, jsonPhoto := range jsonPhotos {
		parsedPhoto := jsonPhoto.(map[string]interface{})
		// Possible photo resolutions: 75, 130, 604, 807, 1280, 2560
		// We need to add the one with the best quality
		id := strconv.Itoa(int(parsedPhoto["id"].(float64)))
		if link2560, ok := parsedPhoto["photo_2560"]; ok {
			photos = append(photos, Photo{id, link2560.(string)})
		} else if link1280, ok := parsedPhoto["photo_1280"]; ok {
			photos = append(photos, Photo{id, link1280.(string)})
		} else if link807, ok := parsedPhoto["photo_807"]; ok {
			photos = append(photos, Photo{id, link807.(string)})
		} else if link604, ok := parsedPhoto["photo_604"]; ok {
			photos = append(photos, Photo{id, link604.(string)})
		}
	}
	return photos
}

func downloadPhoto(path string, name string, url string) {

}
