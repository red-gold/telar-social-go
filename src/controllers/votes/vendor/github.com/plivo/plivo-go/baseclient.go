package plivo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const sdkVersion = "4.1.1"

type ClientOptions struct {
	HttpClient *http.Client
}

type BaseClient struct {
	httpClient *http.Client

	AuthId    string
	AuthToken string

	BaseUrl   *url.URL
	userAgent string

	RequestInterceptor  func(request *http.Request)
	ResponseInterceptor func(response *http.Response)
}

func (client *BaseClient) NewRequest(method string, params interface{}, baseRequestString string, formatString string,
	formatParams ...interface{}) (request *http.Request, err error) {

	if client == nil || client.httpClient == nil {
		err = errors.New("client and httpClient cannot be nil")
		return
	}

	for i, param := range formatParams {
		if param == nil || param == "" {
			err = errors.New(fmt.Sprintf("Request path parameter #%d is nil/empty but should not be so.", i))
			return
		}
	}

	requestUrl := *client.BaseUrl

	requestUrl.Path = fmt.Sprintf(baseRequestString, fmt.Sprintf(formatString, formatParams...))
	var buffer = new(bytes.Buffer)
	if method == "GET" {
		var values url.Values
		if values, err = query.Values(params); err != nil {
			return
		}

		requestUrl.RawQuery = values.Encode()
	} else {
		if reflect.ValueOf(params).Kind().String() != "map" {
			if err = json.NewEncoder(buffer).Encode(params); err != nil {
				return
			}
		} else if reflect.ValueOf(params).Kind().String() == "map" && !reflect.ValueOf(params).IsNil() {
			if err = json.NewEncoder(buffer).Encode(params); err != nil {
				return
			}
		}

	}

	request, err = http.NewRequest(method, requestUrl.String(), buffer)

	request.Header.Add("User-Agent", client.userAgent)
	request.Header.Add("Content-Type", "application/json")

	request.SetBasicAuth(client.AuthId, client.AuthToken)

	return
}

func (client *BaseClient) ExecuteRequest(request *http.Request, body interface{}) (err error) {
	if client == nil {
		return errors.New("client cannot be nil")
	}

	if client.httpClient == nil {
		return errors.New("httpClient cannot be nil")
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err == nil && data != nil && len(data) > 0 {
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			if body != nil {
				err = json.Unmarshal(data, body)
			}
		} else {
			if string(data) == "{}" && response.StatusCode == 404 {
				err = errors.New(string("Resource not found exception \n" + response.Status))
			} else {
				err = errors.New(string(data))
			}
		}
	}

	return
}
