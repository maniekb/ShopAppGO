package utils

import (
    "encoding/json"
    "errors"
    "net/http"

    "golang.org/x/oauth2"
    facebookOAuth "golang.org/x/oauth2/facebook"
	"example/web-service-gin/config"
)

type FacebookUserDetails struct {
    ID    string
    Name  string
    Email string
}

func GetFacebookOAuthConfig() *oauth2.Config {

	config, _ := config.LoadConfig(".")
    return &oauth2.Config{
        ClientID:     config.FacebookClientID,
        ClientSecret: config.FacebookClientSecret,
        RedirectURL:  config.FacebookRedirectUrl,
        Endpoint:     facebookOAuth.Endpoint,
        Scopes:       []string{"email"},
    }
}

func GetFacebookUser(token string) (FacebookUserDetails, error) {
    var userDetails FacebookUserDetails
    req, _ := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email&access_token="+token, nil)
    res, err := http.DefaultClient.Do(req)

    if err != nil {
        return FacebookUserDetails{}, errors.New("Error occurred while getting information from Facebook")
    }

    decoder := json.NewDecoder(res.Body)
    decoderErr := decoder.Decode(&userDetails)
    defer res.Body.Close()

    if decoderErr != nil {
        return FacebookUserDetails{}, errors.New("Error occurred while getting information from Facebook")
    }

    return userDetails, nil
}