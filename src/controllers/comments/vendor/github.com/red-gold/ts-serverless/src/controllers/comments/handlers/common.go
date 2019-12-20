package handlers

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexellis/hmac"
	coreConfig "github.com/red-gold/telar-core/config"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
)

// functionCall send request to another function/microservice using HMAC validation
func functionCall(method string, bytesReq []byte, url string, header map[string][]string) ([]byte, error) {
	prettyURL := utils.GetPrettyURLf(url)
	bodyReader := bytes.NewBuffer(bytesReq)

	httpReq, httpErr := http.NewRequest(method, *coreConfig.AppConfig.Gateway+prettyURL, bodyReader)
	if httpErr != nil {
		return nil, httpErr
	}

	payloadSecret, psErr := utils.ReadSecret("payload-secret")

	if psErr != nil {
		return nil, fmt.Errorf("couldn't get payload-secret: %s", psErr.Error())
	}

	digest := hmac.Sign(bytesReq, []byte(payloadSecret))
	httpReq.Header.Set("Content-type", "application/json")
	fmt.Printf("\ndigest: %s, header: %v \n", "sha1="+hex.EncodeToString(digest), server.X_Cloud_Signature)
	httpReq.Header.Add(server.X_Cloud_Signature, "sha1="+hex.EncodeToString(digest))

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
		return nil, fmt.Errorf("failed to call admin check api, invalid status: %s", res.Status)
	}

	return resData, nil
}
