package tencent

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const (
	Algorithm         = "TC3-HMAC-SHA256"
	HTTPRequestMethod = "POST"
	ContentType       = "application/json; charset=utf-8"
)

const (
	SHA256 = "HmacSHA256"
	SHA1   = "HmacSHA1"
)

func hmacSHA256(key []byte, msg string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(msg))
	return mac.Sum(nil)
}

func sha256Hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

// GenerateSignature 生成腾讯云发送短信接口签名
func generateSignature(config *Config, templateId string, messages, phones []string, timestamp string) (string, error) {
	// 格式化日期为UTC时区
	date := time.Now().Format("2006-01-02")
	canonicalUri := "/"
	canonicalQueryString := ""
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\n", ContentType, config.requestUrl)
	signedHeaders := "content-type;host"

	// 创建请求参数的JSON表示
	params := map[string]interface{}{
		"PhoneNumberSet":   phones,
		"SmsSdkAppId":      config.SDKAppId,
		"SignName":         config.Signature,
		"TemplateId":       templateId,
		"TemplateParamSet": messages,
	}
	payload, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	hashedRequestPayload := sha256Hex(string(payload))
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", HTTPRequestMethod, canonicalUri, canonicalQueryString, canonicalHeaders, signedHeaders, hashedRequestPayload)
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, config.service)
	hashedCanonicalRequest := sha256Hex(canonicalRequest)
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s", Algorithm, timestamp, credentialScope, hashedCanonicalRequest)

	// 生成签名密钥
	secretDate := hmacSHA256([]byte("TC3"+config.AccessKeySecret), date)
	secretService := hmacSHA256(secretDate, config.service)
	secretSigning := hmacSHA256(secretService, "tc3_request")

	signature := hex.EncodeToString(hmacSHA256(secretSigning, stringToSign))

	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", Algorithm, config.AccessKeyId, credentialScope, signedHeaders, signature), nil
}

// GenerateHeadersMap 生成腾讯云短信请求头map
func generateHeadersMap(authorization, timestamp, action, version, territory, requestUrl string) map[string]string {
	headers := map[string]string{
		"Authorization":  authorization,
		"Content-Type":   ContentType,
		"Host":           requestUrl,
		"X-TC-Action":    action,
		"X-TC-Timestamp": timestamp,
		"X-TC-Version":   version,
		"X-TC-Region":    territory,
	}
	return headers
}

// GenerateRequestBody 生成腾讯云短信请求body
func generateRequestBody(phones []string, sdkAppId, signatureName, templateId string, templateParamSet []string) map[string]interface{} {
	requestBody := map[string]interface{}{
		"PhoneNumberSet":   phones,
		"SmsSdkAppId":      sdkAppId,
		"SignName":         signatureName,
		"TemplateId":       templateId,
		"TemplateParamSet": templateParamSet,
	}
	return requestBody
}
