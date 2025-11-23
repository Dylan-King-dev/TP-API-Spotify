package main

import (
	"TPSpotify/router"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TokenData struct {
	Access_token string `json:"access_token"`
	Token_type   string `json:"token_type"`
	Expires_in   int    `json:"expires_in"`
}

type DamsoData struct {
	Items []struct {
		Name   string `json:"name"`
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
		ReleaseDate string `json:"release_date"`
		TotalTracks int    `json:"total_tracks"`
	} `json:"items"`
}

type LaylowData struct {
	Name  string `json:"name"`
	Album struct {
		Name   string `json:"name"`
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
		ReleaseDate string `json:"release_date"`
	} `json:"album"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	DurationMs    int `json:"duration_ms"`
	ExternarlURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

func GetToken() (string, error) {
	urlApi := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", "83ee08a63f2541bfb1a71c287a06ea3d")
	data.Set("client_secret", "a477a7c2704e4e21ae79330cf5055cec")

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodPost, urlApi, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Oupss, une erreur est survenue lors de la création de la requête :", err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Oupss, une erreur est survenue lors de l'envoi de la requête :", err)
		return "", err
	}

	defer res.Body.Close()

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		fmt.Println("Oupss, une erreur est survenue lors de la lecture de la réponse :", errBody.Error())
		return "", errBody
	}
	var token TokenData

	json.Unmarshal(body, &token)
	return token.Access_token, nil

}

func GetDamso() {
	token, err := GetToken()
	if err != nil {
		fmt.Println("Impossible de récupérer le token :", err)
		return
	}

	urlApi := "https://api.spotify.com/v1/artists/2UwqpfQtNuhBwviIC0f2ie/albums"

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	req, errReq := http.NewRequest(http.MethodGet, urlApi, nil)
	if errReq != nil {
		fmt.Println("Oupss, une erreur est survenue lors de la création de la requête :", errReq.Error())
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	res, errResp := httpClient.Do(req)
	if errResp != nil {
		fmt.Println("Oupss, une erreur est survenue lors de l'envoi de la requête :", errResp.Error())
		return
	}

	defer res.Body.Close()

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		fmt.Println("Oupss, une erreur est survenue lors de la lecture de la réponse :", errBody.Error())

	}

	var decodeData DamsoData

	json.Unmarshal(body, &decodeData)
	fmt.Println(decodeData.Items[0].Name, decodeData.Items[0].Images, decodeData.Items[0].ReleaseDate, decodeData.Items[0].TotalTracks)
}

func GetLaylow() {
	token, err := GetToken()
	if err != nil {
		fmt.Println("Impossible de récupérer le token :", err)
		return
	}

	urlApi := "https://api.spotify.com/v1/tracks/67Pf31pl0PfjBfUmvYNDCL?si=3ae39232f83d4963"

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, errReq := http.NewRequest(http.MethodGet, urlApi, nil)
	if errReq != nil {
		fmt.Println("Oupss, une erreur est survenue lors de la création de la requête :", errReq.Error())
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	res, errResp := httpClient.Do(req)
	if errResp != nil {
		fmt.Println("Oupss, une erreur est survenue lors de l'envoi de la requête :", errResp.Error())
		return
	}

	defer res.Body.Close()

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		fmt.Println("Oupss, une erreur est survenue lors de la lecture de la réponse :", errBody.Error())

	}
	var decodeData LaylowData

	json.Unmarshal(body, &decodeData)
	fmt.Println(decodeData.Name, decodeData.Album, decodeData.Artists, decodeData.DurationMs, decodeData.ExternarlURLs)
}

func main() {

	router.SetupRoutes()

	port := ":8080"
	fmt.Println("Server running at http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
