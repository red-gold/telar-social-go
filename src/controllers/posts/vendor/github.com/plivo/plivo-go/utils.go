package plivo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"sort"
	"strings"
)

func Numbers(numbers ...string) string {
	return strings.Join(numbers, "<")
}

func headersWithSep(headers map[string]string, keyValSep, itemSep string, escape bool) string {
	v := url.Values{}
	for key, value := range headers {
		v.Set(key, value)
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		var prefix string
		if escape {
			prefix = url.QueryEscape(k) + keyValSep
		} else {
			prefix = k + keyValSep
		}

		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteString(itemSep)
			}
			buf.WriteString(prefix)
			if escape {
				buf.WriteString(url.QueryEscape(v))
			} else {
				buf.WriteString(v)
			}
		}
	}
	return buf.String()
}

// Some code from encode.go from the Go Standard Library
func Headers(headers map[string]string) string {
	return headersWithSep(headers, "=", ",", true)
}

// The old signature validation is deprecated. Will be marked deprecated in the next release.

func ComputeSignature(authToken, uri string, params map[string]string) string {
	originalString := fmt.Sprintf("%s%s", uri, headersWithSep(params, "", "", false))
	logrus.Infof("originalString: %s\n", originalString)
	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(originalString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func ValidateSignature(authToken, uri string, params map[string]string, signature string) bool {
	return ComputeSignature(authToken, uri, params) == signature
}

// Adding V2 signature validation

func ComputeSignatureV2(authToken, uri string, nonce string) string {
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	var originalString string = parsedUrl.Scheme + "://" + parsedUrl.Host + parsedUrl.Path + nonce
	mac := hmac.New(sha256.New, []byte(authToken))
	mac.Write([]byte(originalString))
	var messageMAC string = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return messageMAC
}

func ValidateSignatureV2(uri string, nonce string, signature string, authToken string) bool {
	return ComputeSignatureV2(authToken, uri, nonce) == signature
}
