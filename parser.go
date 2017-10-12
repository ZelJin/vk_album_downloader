package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
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
	userIDPtr := flag.String("uid", "1", "User id of your Vkontakte account.")
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
		fmt.Print("Parsing album" + album.title)
		// Create a folder for that album
		albumFolder := path.Join(*pathPtr, album.title)
		os.Mkdir(albumFolder, 0777)
		photos := parsePhotos(*userIDPtr, album.id)
		for _, photo := range photos {
			// Save photo by url to the destination folder
			downloadPhoto(albumFolder, photo.id, photo.link)
		}

	}

}

// parseAlbums call VK API to get the list of user's albums.
// returns an array of album metadata
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

// parsePhotos calls VK API to get the list of photos in the album.
// returns an array of photo metadata
func parsePhotos(uid string, albumID string) []Photo {
	// Initialize an empty []Album object
	var photos []Photo

	// Perform a request to the API
	url := fmt.Sprintf("https://api.vk.com/method/photos.get.json?owner_id=%v&album_id=%v&v=5.34", uid, albumID)
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

// downloadPhoto downloads an online picture from a remote URL and stores it
// in a folder defined by albumPath variable.
func downloadPhoto(albumPath string, name string, url string) {
	//fmt.Println(url)
	response, error := http.Get(url)
	if error != nil {
		println(error)
	}
	imageData, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	ioutil.WriteFile(path.Join(albumPath, name+".jpg"), imageData, 0777)
	println("Downloaded photo " + name)
}
