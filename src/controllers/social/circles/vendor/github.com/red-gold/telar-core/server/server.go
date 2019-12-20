package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexellis/hmac"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	handler "github.com/openfaas-incubator/go-function-sdk"
	cf "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/utils"
	uuid "github.com/satori/go.uuid"
)

type Request struct {
	Body        []byte
	Header      http.Header
	QueryString string
	Method      string
	IpAddress   string
	UserID      uuid.UUID
	Username    string
	Avatar      string
	DisplayName string
	SystemRole  string
	CookieMap   string
	params      httprouter.Params
}

type RouteProtection int

const (
	RouteProtectionHMAC RouteProtection = iota
	RouteProtectionAdmin
	RouteProtectionCookie
	RouteProtectionPublic
)

const (
	X_Cloud_Signature = "X-Cloud-Signature"
)

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (req Request) GetParamByName(name string) string {
	for i := range req.params {
		if req.params[i].Key == name {
			return req.params[i].Value
		}
	}
	return ""
}

type Handle func(Request) (handler.Response, error)

type HandleWR func(http.ResponseWriter, *http.Request, Request) (handler.Response, error)

// ReqWR request handler with http.ResponseWriter and http.Request
func ReqWR(funcHandler HandleWR, protected RouteProtection) httprouter.Handle {
	config := &cf.AppConfig
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		for name, headers := range r.Header {
			for _, h := range headers {
				fmt.Printf("\n%v: %v\n", name, h)
			}
		}
		req := handleParseRequest(r)
		req.params = ps

		// Reading cookie in protected request
		if protected != RouteProtectionPublic {
			checkProtection(w, r, &req, config, protected)
		}
		result, resultErr := funcHandler(w, r, req)
		parseHeader(w, r, result, resultErr)
		w.Write(result.Body)
	}
}

// ReqFileWR request file handler with http.ResponseWriter and http.Request
func ReqFileWR(funcHandler HandleWR, protected RouteProtection) httprouter.Handle {
	config := &cf.AppConfig
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		req := handleParseFileRequest(r)
		req.params = ps

		// Reading cookie in protected request
		if protected != RouteProtectionPublic {
			checkProtection(w, r, &req, config, protected)
		}
		result, resultErr := funcHandler(w, r, req)
		parseHeader(w, r, result, resultErr)
		w.Write(result.Body)
	}
}

// Req request handler
func Req(funcHandler Handle, protected RouteProtection) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		req := handleParseRequest(r)
		requestHandler(w, r, ps, funcHandler, req, protected)
	}
}

//  requestHandler
func requestHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, funcHandler Handle, req Request, protected RouteProtection) {
	config := &cf.AppConfig

	req.params = ps
	// Reading cookie in protected request
	if protected != RouteProtectionPublic {
		checkProtection(w, r, &req, config, protected)
	}
	result, resultErr := funcHandler(req)
	parseHeader(w, r, result, resultErr)

	w.Write(result.Body)
}

// checkProtection check protection
func checkProtection(w http.ResponseWriter, r *http.Request, req *Request, config *cf.Configuration, protected RouteProtection) {

	if protected == RouteProtectionHMAC {
		presented, hmacErr := checkHmacPresent(req)
		if hmacErr != nil {
			// TODO: return on invalid request
			writeError(w, "invalid HMAC digest!", hmacErr.Error(), http.StatusUnauthorized)
		}
		if !presented {
			writeError(w, "HMAC is not presented!", "", http.StatusUnauthorized)
		}

	} else if protected == RouteProtectionCookie || protected == RouteProtectionAdmin {
		presented, hmacErr := checkHmacPresent(req)
		if hmacErr != nil {
			// TODO: return on invalid request
			fmt.Printf("invalid HMAC digest: %s", hmacErr.Error())
		}
		if !presented {
			// TODO: return on invalid request

			// Read cookie
			cookieMap, cookieErr := readCookie(w, r, config)
			if cookieErr != nil {
				writeError(w, "Internal Error happened in reading cookie!",
					fmt.Sprintf("Unable to read cookies : %s", cookieErr.Error()), http.StatusInternalServerError)
			}

			// Parse cookie to claim
			claims, cookieErr := parseCookie(w, cookieMap, config)
			if cookieErr != nil {
				fmt.Printf("Error in reading cookie error: %s", cookieErr.Error())
			}

			// Parse claim to request
			parseErr := parseClaim(req, claims, protected)
			if parseErr != nil {
				writeError(w, "Can not parse claim", parseErr.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func parseHeader(w http.ResponseWriter, r *http.Request, result handler.Response, resultErr error) {
	if result.Header != nil {
		for k, v := range result.Header {
			w.Header()[k] = v
		}
	}

	if resultErr != nil {
		log.Print(resultErr)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if result.StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(result.StatusCode)
		}
	}
}

// handleParseRequest parse the request to openfaas handler.request
func handleParseRequest(r *http.Request) Request {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		bodyBytes, bodyErr := ioutil.ReadAll(r.Body)

		if bodyErr != nil {
			log.Printf("Error reading body from request.")
		}

		input = bodyBytes
	}

	return Request{
		Body:        input,
		Header:      r.Header,
		Method:      r.Method,
		QueryString: r.URL.RawQuery,
		IpAddress:   utils.GetIPAdress(r),
	}

}

// handleParseFileRequest parse the request to openfaas handler.request
func handleParseFileRequest(r *http.Request) Request {

	return Request{
		Header:      r.Header,
		Method:      r.Method,
		QueryString: r.URL.RawQuery,
		IpAddress:   utils.GetIPAdress(r),
	}

}

// readCookie read cookies in a map
func readCookie(w http.ResponseWriter, r *http.Request, config *cf.Configuration) (map[string]*http.Cookie, error) {

	cookieHeader, errCHeader := r.Cookie(*config.HeaderCookieName)
	if errCHeader != nil {
		writeError(w, "Cookie Header not found.", errCHeader.Error(), http.StatusUnauthorized)
		return nil, errCHeader

	}

	cookiePayload, errCPayload := r.Cookie(*config.PayloadCookieName)
	if errCPayload != nil {
		writeError(w, "Cookie Payload not found.", errCPayload.Error(), http.StatusUnauthorized)
		return nil, errCPayload
	}

	cookieSignature, errCSignature := r.Cookie(*config.SignatureCookieName)
	if errCSignature != nil {
		writeError(w, "Cookie Signature not found.", errCSignature.Error(), http.StatusUnauthorized)
		return nil, errCSignature
	}

	cookies := make(map[string]*http.Cookie)
	cookies["header"] = cookieHeader
	cookies["payload"] = cookiePayload
	cookies["sign"] = cookieSignature
	return cookies, nil
}

// parseCookie
func parseCookie(w http.ResponseWriter, cookieMap map[string]*http.Cookie, config *cf.Configuration) (jwt.MapClaims, error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		err = r.(error)
	// 	}
	// }()
	keydata, err := ioutil.ReadFile(*config.PublicKeyPath)
	if err != nil {
		writeError(w, "Error happened in reading cookie",
			fmt.Sprintf("unable to read path: %s, error: %s", *config.PublicKeyPath, err.Error()), http.StatusInternalServerError)
		return nil, err
	}

	publicKey, keyErr := jwt.ParseECPublicKeyFromPEM(keydata)
	if keyErr != nil {
		writeError(w, "Internal Error happened in reading cookie!",
			fmt.Sprintf("Unable to read public key : %s", keyErr.Error()), http.StatusInternalServerError)
		return nil, keyErr
	}

	cookie := fmt.Sprintf("%s.%s.%s", cookieMap["header"].Value, cookieMap["payload"].Value, cookieMap["sign"].Value)
	parsed, parseErr := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if parseErr != nil {
		log.Println(parseErr, cookie)
		writeError(w, "Unable to decode cookie, please clear your cookies and sign-in again", parseErr.Error(), http.StatusUnauthorized)
		return nil, parseErr
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Token claim is not valid!")
	}
}

func parseClaim(req *Request, claims jwt.MapClaims, rp RouteProtection) error {
	claimMap, ok := claims["claim"].(map[string]interface{})
	if ok {

		role, roleOk := claimMap["role"].(string)
		if roleOk {
			req.SystemRole = role
		}

		if rp == RouteProtectionAdmin && role != "admin" {
			return fmt.Errorf("adminAccessRole")
		}
		userId, userIdOk := claimMap["uid"].(string)
		fmt.Printf("UserID from claim %s", userId)
		if userIdOk {
			userUUID, userUuidErr := uuid.FromString(userId)
			if userUuidErr != nil {
				return userUuidErr
			}
			req.UserID = userUUID
		}
		username, usernameOk := claimMap["email"].(string)
		if usernameOk {
			req.Username = username
		}
		avatar, avatarOk := claimMap["avatar"].(string)
		if avatarOk {
			req.Avatar = avatar
		}
		displayName, displayNameOk := claimMap["displayName"].(string)
		if displayNameOk {
			req.DisplayName = displayName
		}

	}

	return nil
}

// validateRequest
func validateRequest(req *Request) (err error) {
	payloadSecret, err := utils.ReadSecret("payload-secret")

	if err != nil {
		return fmt.Errorf("couldn't get payload-secret: %s", err.Error())
	}

	xCloudSignature := req.Header.Get(X_Cloud_Signature)

	fmt.Printf("\nxCloudSignature: %s\n", xCloudSignature)
	err = hmac.Validate(req.Body, xCloudSignature, payloadSecret)

	if err != nil {
		return err
	}

	return nil
}

// checkHmacPresent check whether hmac header presented
func checkHmacPresent(req *Request) (bool, error) {

	xCloudSignature := req.Header.Get(X_Cloud_Signature)
	// Loop over header names
	for name, values := range req.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
	if xCloudSignature != "" {
		validErr := validateRequest(req)
		if validErr != nil {
			return true, validErr
		}
		userId := req.Header.Get("uid")
		userUUID, userUuidErr := uuid.FromString(userId)
		if userUuidErr != nil {
			return true, userUuidErr
		}
		req.UserID = userUUID
		req.Username = req.Header.Get("email")
		req.Avatar = req.Header.Get("avatar")
		req.DisplayName = req.Header.Get("displayName")
		req.SystemRole = req.Header.Get("role")
		return true, nil
	}
	return false, nil
}

func writeError(w http.ResponseWriter, err string, logErr string, status int) {
	log.Println(err)
	log.Println(logErr)
	w.Write([]byte(err))
	w.WriteHeader(status)
}
