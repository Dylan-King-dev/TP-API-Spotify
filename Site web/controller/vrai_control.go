package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
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
	DurationMs   int `json:"duration_ms"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

type LaylowViewData struct {
	LaylowData
	Duration string
}

func Ameliore(durationMs int) string {
	totalSeconds := durationMs / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func GetToken() (string, error) {
	clientID := "83ee08a63f2541bfb1a71c287a06ea3d"
	clientSecret := "a477a7c2704e4e21ae79330cf5055cec"

	urlApi := "https://accounts.spotify.com/api/token"

	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	data := url.Values{}
	data.Set("grant_type", "client_credentials") // only grant_type in body

	httpClient := &http.Client{
		Timeout: 15 * time.Second, // safe timeout
	}

	req, err := http.NewRequest(http.MethodPost, urlApi, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	// Check HTTP status
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("Spotify API error: %s", string(body))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var token TokenData
	if err := json.Unmarshal(body, &token); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %v", err)
	}

	return token.Access_token, nil
}

func GetDamso() (*DamsoData, error) {
	token, err := GetToken()
	if err != nil {
		return nil, fmt.Errorf("Impossible de récupérer le token: %v", err)
	}

	urlApi := "https://api.spotify.com/v1/artists/2UwqpfQtNuhBwviIC0f2ie/albums"

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	req, errReq := http.NewRequest(http.MethodGet, urlApi, nil)
	if errReq != nil {
		return nil, fmt.Errorf("Erreur envoi requête: %v", errReq)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	res, errResp := httpClient.Do(req)
	if errResp != nil {
		return nil, fmt.Errorf("Erreur envoi requête: %v", errResp)
	}

	defer res.Body.Close()

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		return nil, fmt.Errorf("Erreur lecture réponse: %v", errBody)
	}

	var decodeData DamsoData

	json.Unmarshal(body, &decodeData)

	return &decodeData, nil
}

func GetLaylow() (*LaylowViewData, error) {
	token, err := GetToken()
	if err != nil {
		return nil, fmt.Errorf("Impossible de récupérer le token: %v", err)
	}

	urlApi := "https://api.spotify.com/v1/tracks/67Pf31pl0PfjBfUmvYNDCL?si=3ae39232f83d4963"

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, errReq := http.NewRequest(http.MethodGet, urlApi, nil)
	if errReq != nil {
		return nil, fmt.Errorf("Erreur envoi requête: %v", errReq)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	res, errResp := httpClient.Do(req)
	if errResp != nil {
		return nil, fmt.Errorf("Erreur envoi requête: %v", errResp)
	}

	defer res.Body.Close()

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		return nil, fmt.Errorf("Erreur lecture réponse: %v", errBody)
	}

	var decodeData LaylowData

	json.Unmarshal(body, &decodeData)

	Duration := Ameliore(decodeData.DurationMs)

	viewData := LaylowViewData{
		LaylowData: decodeData,
		Duration:   Duration,
	}

	return &viewData, nil
}

func DaHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/Wappers.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func DamsoPage(w http.ResponseWriter, r *http.Request) {
	damsoData, err := GetDamso()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("template/Damso.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, damsoData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LaylowPage(w http.ResponseWriter, r *http.Request) {
	laylowViewData, err := GetLaylow()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("template/Laylow.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, laylowViewData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
