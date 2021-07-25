package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/alexellis/hmac"
	uuid "github.com/gofrs/uuid"
	coreConfig "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	models "github.com/red-gold/ts-serverless/micros/posts/models"
)

const contentMaxLength = 20

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// StringRand
func StringRand(length int) string {
	return StringWithCharset(length, charset)
}

// generatPostURLKey
func generatPostURLKey(socialName, body, postId string) string {
	contetn := body
	if contentMaxLength <= len(body) {
		contetn = body[:contentMaxLength]
	}

	return strings.ToLower(fmt.Sprintf("%s_%s-post-%s-%s", socialName, strings.ReplaceAll(contetn, " ", "-"), strings.Split(postId, "-")[0], StringRand(5)))
}

// functionCall send request to another function/microservice using cookie validation
func functionCall(method string, bytesReq []byte, url string, header map[string][]string) ([]byte, error) {
	prettyURL := utils.GetPrettyURLf(url)
	bodyReader := bytes.NewBuffer(bytesReq)

	httpReq, httpErr := http.NewRequest(method, *coreConfig.AppConfig.InternalGateway+prettyURL, bodyReader)
	if httpErr != nil {
		return nil, httpErr
	}
	payloadSecret := *coreConfig.AppConfig.PayloadSecret

	digest := hmac.Sign(bytesReq, []byte(payloadSecret))
	httpReq.Header.Set("Content-type", "application/json")
	fmt.Printf("\ndigest: %s, header: %v \n", "sha1="+hex.EncodeToString(digest), types.HeaderHMACAuthenticate)
	httpReq.Header.Add(types.HeaderHMACAuthenticate, "sha1="+hex.EncodeToString(digest))
	if header != nil {
		for k, v := range header {
			httpReq.Header[k] = v
		}
	}
	c := http.Client{}
	res, reqErr := c.Do(httpReq)
	fmt.Printf("\nRes: %v\n", res)
	if reqErr != nil {
		return nil, fmt.Errorf("Error while sending admin check request!: %s", reqErr.Error())
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	resData, readErr := ioutil.ReadAll(res.Body)
	if resData == nil || readErr != nil {
		return nil, fmt.Errorf("failed to read response from admin check request.")
	}

	if res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, NotFoundHTTPStatusError
		}
		return nil, fmt.Errorf("failed to call %s api, invalid status: %s", prettyURL, res.Status)
	}

	return resData, nil
}

// getUserProfileByID Get user profile by user ID
func getUserProfileByID(userID uuid.UUID) (*models.UserProfileModel, error) {
	profileURL := fmt.Sprintf("/profile/dto/id/%s", userID.String())
	foundProfileData, err := functionCall(http.MethodGet, []byte(""), profileURL, nil)
	if err != nil {
		if err == NotFoundHTTPStatusError {
			return nil, nil
		}
		log.Error("functionCall (%s) -  %s", profileURL, err.Error())
		return nil, fmt.Errorf("getUserProfileByID/functionCall")
	}
	var foundProfile models.UserProfileModel
	err = json.Unmarshal(foundProfileData, &foundProfile)
	if err != nil {
		log.Error("Unmarshal foundProfile -  %s", err.Error())
		return nil, fmt.Errorf("getUserProfileByID/unmarshal")
	}
	return &foundProfile, nil
}
