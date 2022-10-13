package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Client(AuthKey string) (*http.Client, string, string) {
	var client *http.Client
	client = &http.Client{}
	baseUrl := os.Getenv("BASE_URL")

	return client, baseUrl, AuthKey
}

func GetKey(client *http.Client, baseUrl string, AuthKey string) (*string, error) {
	req, err := http.NewRequest("GET", baseUrl+"/get-key", nil)
	req.Header.Add("bot-auth", AuthKey)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	r := io.MultiReader(resp.Body)

	var v R

	err = json.NewDecoder(r).Decode(&v)
	if err != nil {
		log.Println(err)
	}

	return &v.Result.AuthKey, nil
}
