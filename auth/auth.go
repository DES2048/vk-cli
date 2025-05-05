package auth

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api/oauth"
)

func input(prompt string, val *string) {
	fmt.Print(prompt)
	fmt.Scanln(val)
}

// WriteTokenToFile function
func WriteTokenToFile(token *oauth.UserToken, filename string) error {
	return os.WriteFile(filename, []byte(token.AccessToken), 0o666)
}

// Auth function
func Auth() (*oauth.UserToken, error) {
	fmt.Println("Авторизация вк")
	var vkClientIDStr string
	// var vkClientSecret string
	var vkURLStr string

	input("App Id:", &vkClientIDStr)

	vkClientID, err := strconv.Atoi(vkClientIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse App Id")
	}

	oauthURL := oauth.ImplicitFlowUser(oauth.UserParams{
		ClientID: vkClientID,
		Scope:    oauth.ScopeUserOffline + oauth.ScopeUserVideo + oauth.ScopeGroupPhotos + oauth.ScopeUserGroups,
	})

	fmt.Printf("Url для авторизации: %s\n", oauthURL.String())
	input("Url с токеном:", &vkURLStr)

	vkURL, err := url.Parse(vkURLStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %s", err)
	}

	t, err := oauth.NewUserTokenFromURL(vkURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from url: %s", err)
	}

	return t, nil
}

// CheckAuth ...
func CheckAuth(tokenFile string) (string, error) {
	file, err := os.Open(tokenFile)
	if err != nil {
		if os.IsNotExist(err) {
			token, err := Auth()
			if err != nil {
				return "", err
			}

			err = WriteTokenToFile(token, tokenFile)
			if err != nil {
				return "", fmt.Errorf("failed to write token to file: %s", err)
			}
		}

		return "", fmt.Errorf("failed to open token file: %s", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read token from file: %s", err)
	}

	return strings.Trim(string(data), "\n"), nil
}
