package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	coreConfig "github.com/red-gold/telar-core/config"
)

func Call(url string, data []byte) error {
	c := http.Client{}
	appConfig := coreConfig.AppConfig
	reader := bytes.NewBuffer(data)
	fullURL := fmt.Sprintf("%s/%s", *appConfig.Gateway, url)
	req, _ := http.NewRequest(http.MethodPost, fullURL, reader)
	res, err := c.Do(req)
	if err != nil {
		log.Println(fullURL, err)
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	return nil
}
