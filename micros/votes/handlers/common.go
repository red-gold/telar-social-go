package handlers

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexellis/hmac"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	coreConfig "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
)

type ResultAsync struct {
	Result []byte
	Error  error
}

type UserInfoInReq struct {
	UserId      uuid.UUID `json:"uid"`
	Username    string    `json:"email"`
	DisplayName string    `json:"displayName"`
	SocialName  string    `json:"socialName"`
	Avatar      string    `json:"avatar"`
	Banner      string    `json:"banner"`
	TagLine     string    `json:"tagLine"`
	SystemRole  string    `json:"role"`
	CreatedDate int64     `json:"createdDate"`
}

// getHeadersFromUserInfoReq
func getHeadersFromUserInfoReq(info *UserInfoInReq) map[string][]string {
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{info.UserId.String()}
	userHeaders["email"] = []string{info.Username}
	userHeaders["avatar"] = []string{info.Avatar}
	userHeaders["banner"] = []string{info.Banner}
	userHeaders["tagLine"] = []string{info.TagLine}
	userHeaders["displayName"] = []string{info.DisplayName}
	userHeaders["socialName"] = []string{info.SocialName}
	userHeaders["role"] = []string{info.SystemRole}

	return userHeaders
}

// getUserInfoReq
func getUserInfoReq(c *fiber.Ctx) *UserInfoInReq {
	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		return &UserInfoInReq{}
	}
	userInfoInReq := &UserInfoInReq{
		UserId:      currentUser.UserID,
		Username:    currentUser.Username,
		Avatar:      currentUser.Avatar,
		DisplayName: currentUser.DisplayName,
		SystemRole:  currentUser.SystemRole,
	}
	return userInfoInReq

}

// getHeaderInfoReq
func getHeaderInfoReq(c *fiber.Ctx) map[string][]string {
	return getHeadersFromUserInfoReq(getUserInfoReq(c))
}

// functionCall send request to another function/microservice using HMAC validation
func functionCall(method string, bytesReq []byte, url string, header map[string][]string) ([]byte, error) {
	prettyURL := utils.GetPrettyURLf(url)
	bodyReader := bytes.NewBuffer(bytesReq)

	fullURL := *coreConfig.AppConfig.InternalGateway + prettyURL
	httpReq, httpErr := http.NewRequest(method, fullURL, bodyReader)
	if httpErr != nil {
		return nil, httpErr
	}

	digest := hmac.Sign(bytesReq, []byte(*coreConfig.AppConfig.PayloadSecret))
	httpReq.Header.Set("Content-type", "application/json")
	log.Info("\ndigest: %s, header: %v \n", "sha1="+hex.EncodeToString(digest), types.HeaderHMACAuthenticate)
	httpReq.Header.Add(types.HeaderHMACAuthenticate, "sha1="+hex.EncodeToString(digest))

	if header != nil {
		for k, v := range header {
			httpReq.Header[k] = v
		}
	}

	c := http.Client{}
	res, reqErr := c.Do(httpReq)
	if reqErr != nil {
		return nil, fmt.Errorf("Error while sending request [%s] - %s", fullURL, reqErr.Error())
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	resData, readErr := ioutil.ReadAll(res.Body)
	if resData == nil || readErr != nil {
		return nil, fmt.Errorf("failed to read response from [%s].", fullURL)
	}

	if res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to call api, invalid status: %s", res.Status)
	}

	return resData, nil
}

// readPostAsync Read post async
func readPostAsync(postId uuid.UUID, infoReq *UserInfoInReq) <-chan ResultAsync {
	r := make(chan ResultAsync)
	go func() {
		defer close(r)
		postURL := fmt.Sprintf("/posts/%s", postId.String())

		post, err := functionCall(http.MethodGet, []byte(""), postURL, getHeadersFromUserInfoReq(infoReq))
		if err != nil {
			r <- ResultAsync{Error: err}
			return
		}
		r <- ResultAsync{Result: post}

	}()
	return r
}
