package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/adrg/xdg"
	"golang.org/x/oauth2"
)

// GetClient retrieves a token, saves the token, then returns the generated client.
func GetClient(ctx context.Context, config *oauth2.Config) (client *http.Client, err error) {
	tokFile, err := xdg.CacheFile("gphotos-fb/token.json")
	if err != nil {
		return
	}

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(ctx, config)
		if err != nil {
			return
		}
		err = saveToken(tokFile, tok)
		if err != nil {
			return
		}
	}
	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("Unable to read authorization code: %w", err)
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve token from web: %w", err)
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (tok *oauth2.Token, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = cerr
		}
	}()

	tok = &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) (err error) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}

	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = cerr
		}
	}()

	return json.NewEncoder(f).Encode(token)
}
