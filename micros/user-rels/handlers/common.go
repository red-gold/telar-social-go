package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
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
	models "github.com/red-gold/ts-serverless/micros/user-rels/models"
	socialModels "github.com/red-gold/ts-serverless/micros/user-rels/models"
)

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

	httpReq, httpErr := http.NewRequest(method, *coreConfig.AppConfig.InternalGateway+prettyURL, bodyReader)
	if httpErr != nil {
		return nil, httpErr
	}

	digest := hmac.Sign(bytesReq, []byte(*coreConfig.AppConfig.PayloadSecret))
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

// increaseUserFollowCount Increase user follow count
func increaseUserFollowCount(userId uuid.UUID, inc int, userInfoInReq *UserInfoInReq) {

	actionURL := fmt.Sprintf("/profile/follow/inc/%d/%s", inc, userId.String())

	// Create user headers for http request
	userHeaders := getHeadersFromUserInfoReq(userInfoInReq)

	_, actionErr := functionCall(http.MethodPut, []byte(actionURL), actionURL, userHeaders)

	if actionErr != nil {
		errorMessage := fmt.Sprintf("Function call error: %s - %s", actionURL, actionErr.Error())
		log.Error(errorMessage)
	}
}

// increaseUserFollowerCount Increase user follower count
func increaseUserFollowerCount(userId uuid.UUID, inc int, userInfoInReq *UserInfoInReq) {

	actionURL := fmt.Sprintf("/profile/follower/inc/%d/%s", inc, userId.String())

	// Create user headers for http request
	userHeaders := getHeadersFromUserInfoReq(userInfoInReq)

	_, actionErr := functionCall(http.MethodPut, []byte(actionURL), actionURL, userHeaders)

	if actionErr != nil {
		errorMessage := fmt.Sprintf("Function call error: %s - %s", actionURL, actionErr.Error())
		log.Error(errorMessage)
	}
}

// sendFollowNotification Send follow notification
func sendFollowNotification(model *socialModels.FollowModel, userInfoInReq *UserInfoInReq) {

	// Create user headers for http request
	userHeaders := getHeadersFromUserInfoReq(userInfoInReq)

	URL := fmt.Sprintf("/@/%s", userInfoInReq.SocialName)
	notificationModel := &models.NotificationModel{
		OwnerUserId:          userInfoInReq.UserId,
		OwnerDisplayName:     userInfoInReq.DisplayName,
		OwnerAvatar:          userInfoInReq.Avatar,
		Title:                userInfoInReq.DisplayName,
		Description:          fmt.Sprintf("%s is following you.", userInfoInReq.DisplayName),
		URL:                  URL,
		NotifyRecieverUserId: model.RightUser.UserId,
		TargetId:             model.RightUser.UserId,
		IsSeen:               false,
		Type:                 "follow",
	}
	notificationBytes, marshalErr := json.Marshal(notificationModel)
	if marshalErr != nil {
		fmt.Printf("Cannot marshal notification! error: %s", marshalErr.Error())

	}

	notificationURL := "/notifications"
	_, notificationIndexErr := functionCall(http.MethodPost, notificationBytes, notificationURL, userHeaders)
	if notificationIndexErr != nil {
		fmt.Printf("\nCannot save notification on follow user! error: %s\n", notificationIndexErr.Error())
	}

}
