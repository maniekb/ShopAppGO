package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"example/web-service-gin/config"
)

type GitHubOauthToken struct {
	Access_token string
}

type GitHubUserResult struct {
	Name  string
	Photo string
	Email string
}

type GitHubUserRes struct {
	email  string
	avatar_url string
	login string
}

func GetGitHubOauthToken(code string) (*GitHubOauthToken, error) {
	const rootURl = "https://github.com/login/oauth/access_token"
	
	config, _ := config.LoadConfig(".")
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", config.GitHubClientID)
	values.Add("client_secret", config.GitHubClientSecret)

	query := values.Encode()

	queryString := fmt.Sprintf("%s?%s", rootURl, bytes.NewBufferString(query))
	req, err := http.NewRequest("POST", queryString, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	parsedQuery, err := url.ParseQuery(string(resBody))
	if err != nil {
		return nil, err
	}

	tokenBody := &GitHubOauthToken{
		Access_token: parsedQuery["access_token"][0],
	}

	return tokenBody, nil
}

func GetGitHubUser(accessToken string) (*GitHubUserResult, error) {
	rootUrl := "https://api.github.com/user"

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			DisableCompression: true,
	},
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("blad")
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("blad")
		return nil, errors.New("could not retrieve user")
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("blad")
		return nil, err
	}

	fmt.Println(resBody)

	var GitHubUserRes map[string]interface{}

	if err := json.Unmarshal(resBody, &GitHubUserRes); err != nil {
		return nil, err
	}

	userBody := &GitHubUserResult{
		Email: GitHubUserRes["email"].(string),
		Name:  GitHubUserRes["login"].(string),
		Photo: GitHubUserRes["avatar_url"].(string),
	}

	return userBody, nil
}