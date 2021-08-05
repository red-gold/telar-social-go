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
	"github.com/gofrs/uuid"
	coreConfig "github.com/red-gold/telar-core/config"
	log "github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	models "github.com/red-gold/ts-serverless/micros/vang/models"
)

type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type UserInfoInReq struct {
	UserId      uuid.UUID `json:"userId"`
	Username    string    `json:"username"`
	Avatar      string    `json:"avatar"`
	DisplayName string    `json:"displayName"`
	SystemRole  string    `json:"systemRole"`
}

// getHeadersFromUserInfoReq
func getHeadersFromUserInfoReq(info *UserInfoInReq) map[string][]string {
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{info.UserId.String()}
	userHeaders["email"] = []string{info.Username}
	userHeaders["avatar"] = []string{info.Avatar}
	userHeaders["displayName"] = []string{info.DisplayName}
	userHeaders["role"] = []string{info.SystemRole}

	return userHeaders
}

// getUserInfoReq
func getUserInfoReqFromCurrentUser(currentUser types.UserContext) *UserInfoInReq {
	userInfoInReq := &UserInfoInReq{
		UserId:      currentUser.UserID,
		Username:    currentUser.Username,
		Avatar:      currentUser.Avatar,
		DisplayName: currentUser.DisplayName,
		SystemRole:  currentUser.SystemRole,
	}
	return userInfoInReq
}

// getUserInfoReq
func getUserInfoReq(c *fiber.Ctx) *UserInfoInReq {
	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		return &UserInfoInReq{}
	}
	return getUserInfoReqFromCurrentUser(currentUser)

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

// Dispatch action
func dispatchAction(action Action, roomId uuid.UUID, userInfoInReq *UserInfoInReq) {

	actionURL := fmt.Sprintf("/actions/dispatch/%s", roomId.String())

	actionBytes, marshalErr := json.Marshal(action)
	if marshalErr != nil {
		errorMessage := fmt.Sprintf("Marshal notification Error %s", marshalErr.Error())
		fmt.Println(errorMessage)
	}
	// Create user headers for http request
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{userInfoInReq.UserId.String()}
	userHeaders["email"] = []string{userInfoInReq.Username}
	userHeaders["avatar"] = []string{userInfoInReq.Avatar}
	userHeaders["displayName"] = []string{userInfoInReq.DisplayName}
	userHeaders["role"] = []string{userInfoInReq.SystemRole}

	_, actionErr := functionCall(http.MethodPost, actionBytes, actionURL, userHeaders)

	if actionErr != nil {
		errorMessage := fmt.Sprintf("Cannot send action request! error: %s", actionErr.Error())
		fmt.Println(errorMessage)
	}
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

// getProfileBySocialName Get user profile by social name
func getProfileBySocialName(socialName string) (*models.UserProfileModel, error) {
	profileURL := fmt.Sprintf("/profile/social/%s", socialName)
	foundProfileData, err := functionCall(http.MethodGet, []byte(""), profileURL, nil)
	if err != nil {
		if err == NotFoundHTTPStatusError {
			return nil, nil
		}
		log.Error("functionCall (%s) -  %s", profileURL, err.Error())
		return nil, fmt.Errorf("getProfileBySocialName/functionCall")
	}
	var foundProfile models.UserProfileModel
	err = json.Unmarshal(foundProfileData, &foundProfile)
	if err != nil {
		log.Error("Unmarshal foundProfile -  %s", err.Error())
		return nil, fmt.Errorf("getProfileBySocialName/unmarshal")
	}
	return &foundProfile, nil
}

// getProfilesByUserIds Get user profiles by user IDs
func getProfilesByUserIds(model models.GetProfilesModel, userInfoInReq *UserInfoInReq) ([]models.UserProfileModel, error) {
	profileURL := "/profile/dto/ids"
	body, marshalErr := json.Marshal(model)
	if marshalErr != nil {
		errorMessage := fmt.Sprintf("Marshal models.GetProfilesModel Error %s", marshalErr.Error())
		fmt.Println(errorMessage)
		return nil, marshalErr
	}

	// Create user headers for http request
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{userInfoInReq.UserId.String()}
	userHeaders["email"] = []string{userInfoInReq.Username}
	userHeaders["avatar"] = []string{userInfoInReq.Avatar}
	userHeaders["displayName"] = []string{userInfoInReq.DisplayName}
	userHeaders["role"] = []string{userInfoInReq.SystemRole}

	foundProfilesData, err := functionCall(http.MethodPost, body, profileURL, nil)
	if err != nil {
		if err == NotFoundHTTPStatusError {
			return nil, nil
		}
		log.Error("functionCall (%s) -  %s", profileURL, err.Error())
		return nil, fmt.Errorf("getProfilesByUserIds/functionCall")
	}
	var foundProfile []models.UserProfileModel
	err = json.Unmarshal(foundProfilesData, &foundProfile)
	if err != nil {
		log.Error("Unmarshal foundProfiles -  %s", err.Error())
		return nil, fmt.Errorf("getProfilesByUserIds/unmarshal")
	}
	return foundProfile, nil
}

// dispatchProfileByUserIds Dispatch profile by user Ids
func dispatchProfileByUserIds(model models.DispatchProfilesModel, userInfoInReq *UserInfoInReq) error {
	profileURL := "/profile/dispatch"
	body, marshalErr := json.Marshal(model)
	if marshalErr != nil {
		errorMessage := fmt.Sprintf("Marshal models.DispatchProfilesModel Error %s", marshalErr.Error())
		fmt.Println(errorMessage)
	}

	// Create user headers for http request
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{userInfoInReq.UserId.String()}
	userHeaders["email"] = []string{userInfoInReq.Username}
	userHeaders["avatar"] = []string{userInfoInReq.Avatar}
	userHeaders["displayName"] = []string{userInfoInReq.DisplayName}
	userHeaders["role"] = []string{userInfoInReq.SystemRole}

	_, err := functionCall(http.MethodPost, body, profileURL, userHeaders)
	if err != nil {
		if err == NotFoundHTTPStatusError {
			return nil
		}
		log.Error("functionCall (%s) -  %s", profileURL, err.Error())
		return fmt.Errorf("dispatchProfileByUserIds/functionCall")
	}
	return nil
}
