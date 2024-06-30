package tencent

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sms4go/infra"
	"strconv"
	"strings"
	"time"
)

const (
	Algorithm         = "TC3-HMAC-SHA256"
	HTTPRequestMethod = "POST"
	ContentType       = "application/json; charset=utf-8"
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

func hmacSha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

// GenerateSignature 生成腾讯云发送短信接口签名
func generateSignature(config *Config, templateId string, messages, phones []string, timestamp int64) (string, error) {
	// 格式化日期为UTC时区
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	canonicalUri := "/"
	canonicalQueryString := ""
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-tc-action:%s\n",
		"application/json; charset=utf-8", config.requestUrl, strings.ToLower(config.action))
	signedHeaders := "content-type;host;x-tc-action"

	// 创建请求参数的JSON表示
	params := generateRequestBody(phones, config.SDKAppId, config.Signature, templateId, messages)
	payload, err := json.Marshal(params)

	if err != nil {
		return "", err
	}
	hashedRequestPayload := sha256Hex(string(payload))
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", HTTPRequestMethod, canonicalUri, canonicalQueryString, canonicalHeaders, signedHeaders, hashedRequestPayload)
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, config.service)
	hashedCanonicalRequest := sha256Hex(canonicalRequest)
	stringToSign := fmt.Sprintf("%s\n%d\n%s\n%s", Algorithm, timestamp, credentialScope, hashedCanonicalRequest)

	// 生成签名密钥
	secretDate := hmacSha256(date, "TC3"+config.AccessKeySecret)
	secretService := hmacSha256(config.service, secretDate)
	secretSigning := hmacSha256("tc3_request", secretService)

	signature := hex.EncodeToString([]byte(hmacSha256(stringToSign, secretSigning)))

	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", Algorithm, config.AccessKeyId, credentialScope, signedHeaders, signature), nil
}

// GenerateHeadersMap 生成腾讯云短信请求头map
func generateHeadersMap(authorization, action, version, territory, requestUrl string, timestamp int64) map[string]string {
	headers := map[string]string{
		"Authorization":  authorization,
		"Content-Type":   ContentType,
		"Host":           requestUrl,
		"X-TC-Action":    action,
		"X-TC-Timestamp": strconv.FormatInt(timestamp, 10),
		"X-TC-Version":   version,
		"X-TC-Region":    territory,
	}
	return headers
}

// GenerateRequestBody 生成腾讯云短信请求body
func generateRequestBody(phones []string, sdkAppId, signatureName, templateId string, templateParamSet []string) *infra.TencentPayload {
	payload := &infra.TencentPayload{
		PhoneNumberSet:   phones,
		SignName:         signatureName,
		SmsSdkAppId:      sdkAppId,
		TemplateId:       templateId,
		TemplateParamSet: templateParamSet,
	}
	return payload
}
