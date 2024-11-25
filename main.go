package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	Id           string    `json:"id"`
	LongUrl      string    `json:"longUrl"`
	ShortUrl     string    `json:"shortUrl"`
	CreationTime time.Time `json:"creationTime"`
}

var UrlDB = make(map[string]URL)

func generateShortURL(OriginalUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalUrl))
	fmt.Println("Hasher", hasher)

	data := hasher.Sum(nil)
	fmt.Println("Hasher Data: ", data)
	hash := hex.EncodeToString(data)

	fmt.Println("EncodeToString", hash)
	fmt.Println("Final String", hash[:8])

	return hash[0:8]
}

func createUrl(originalUrl string) string {
	shortUrl := generateShortURL(originalUrl)
	id := shortUrl

	UrlDB[id] = URL{
		Id:           id,
		LongUrl:      originalUrl,
		ShortUrl:     shortUrl,
		CreationTime: time.Now(),
	}

	return shortUrl
}

func getUrl(id string) (URL, error) {
	url, ok := UrlDB[id]

	if !ok {
		return URL{}, errors.New("url not found")
	}

	return url, nil
}

func RootPageUrlHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Helloworld")
}

func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	// Creating request body structure
	var data struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	shortUrl := createUrl(data.URL)
	//fmt.Fprintf(w, shortUrl)

	response := struct {
		ShortURL string `json:"shortUrl"`
	}{ShortURL: shortUrl}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RedirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]

	url, err := getUrl(id)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.LongUrl, http.StatusFound)
}

func main() {
	// fmt.Println("Building URL Shortner")
	// OriginalUrl := "https://github.com/adij1206"
	// generateShortURL(OriginalUrl)

	// Register the handler function to handle all the request to root URL
	http.HandleFunc("/", RootPageUrlHandler)
	http.HandleFunc("/shortner", ShortUrlHandler)
	http.HandleFunc("/redirect/", RedirectUrlHandler)

	// Start the HTTP server on port 3000
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error While STarting Server", err)
	}
}
