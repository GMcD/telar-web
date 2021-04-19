package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexellis/hmac"
	"github.com/gofrs/uuid"
	coreConfig "github.com/red-gold/telar-core/config"
	server "github.com/red-gold/telar-core/server"
	utils "github.com/red-gold/telar-core/utils"
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

// Dispatch action
func dispatchAction(action Action, userInfoInReq *UserInfoInReq) {

	actionURL := fmt.Sprintf("/actions/dispatch/%s", userInfoInReq.UserId.String())

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
		if res.StatusCode == http.StatusNotFound {
			return nil, NotFoundHTTPStatusError
		}
		return nil, fmt.Errorf("failed to call %s api, invalid status: %s", prettyURL, res.Status)
	}

	return resData, nil
}
