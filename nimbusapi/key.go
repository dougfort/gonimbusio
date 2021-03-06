package nimbusapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// StartConjoined starts a conjoined archive, returning a token that can be
// used to identify individual file uploads
func StartConjoined(requester Requester, collectionName string, key string) (
	string, error) {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	values := url.Values{}
	values.Set("action", "start")
	path := fmt.Sprintf("/conjoined/%s?%s", url.QueryEscape(key),
		values.Encode())

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return "", err
	}

	response, err := requester.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		err = HTTPError{response.StatusCode,
			fmt.Sprintf("POST %s %s failed %s", hostName, path, response.Body)}
		return "", err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var conjoinedMap map[string]string
	err = json.Unmarshal(responseBody, &conjoinedMap)
	if err != nil {
		return "", err
	}

	return conjoinedMap["conjoined_identifier"], nil
}

// AbortConjoined aborts a conjoined archive
func AbortConjoined(requester Requester, collectionName string, key string,
	conjoinedIdentifier string) error {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	values := url.Values{}
	values.Set("action", "abort")
	values.Set("conjoined_identifier", conjoinedIdentifier)
	path := fmt.Sprintf("/conjoined/%s?%s", url.QueryEscape(key),
		values.Encode())

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return err
	}

	response, err := requester.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = HTTPError{response.StatusCode,
			fmt.Sprintf("POST %s %s failed %s", hostName, path, response.Body)}
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var conjoinedMap map[string]bool
	err = json.Unmarshal(responseBody, &conjoinedMap)
	if err != nil {
		return err
	}

	if !conjoinedMap["success"] {
		return fmt.Errorf("conjoined action=abort returned false")
	}

	return nil
}

// FinishConjoined marks a conjoined archive as completed
func FinishConjoined(requester Requester, collectionName string, key string,
	conjoinedIdentifier string) error {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)
	values := url.Values{}
	values.Set("action", "finish")
	values.Set("conjoined_identifier", conjoinedIdentifier)
	path := fmt.Sprintf("/conjoined/%s?%s", url.QueryEscape(key),
		values.Encode())

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return err
	}

	response, err := requester.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = HTTPError{response.StatusCode,
			fmt.Sprintf("POST %s %s failed %s", hostName, path, response.Body)}
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var conjoinedMap map[string]bool
	err = json.Unmarshal(responseBody, &conjoinedMap)
	if err != nil {
		return err
	}

	if !conjoinedMap["success"] {
		return fmt.Errorf("conjoined action=finish returned false")
	}

	return nil
}

// ConjoinedParams contains optional arguments to Key.Archive
type ConjoinedParams struct {
	ConjoinedIdentifier string
	ConjoinedPart       int
}

// Archive uploads the data from requestBody
func Archive(requester Requester, collectionName string, key string,
	conjoinedParams *ConjoinedParams, contentLength int64,
	requestBody io.Reader) (string, error) {
	method := "POST"
	hostName := requester.CollectionHostName(collectionName)

	var path string
	if conjoinedParams != nil && conjoinedParams.ConjoinedIdentifier != "" {
		values := url.Values{}
		values.Set("conjoined_identifier", conjoinedParams.ConjoinedIdentifier)
		values.Set("conjoined_part",
			strconv.Itoa(conjoinedParams.ConjoinedPart))
		path = fmt.Sprintf("/data/%s?%s", url.QueryEscape(key), values.Encode())
	} else {
		path = fmt.Sprintf("/data/%s", url.QueryEscape(key))
	}

	request, err := requester.CreateRequest(method, hostName, path, requestBody)
	if err != nil {
		return "", err
	}
	request.ContentLength = contentLength

	response, err := requester.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		err = HTTPError{response.StatusCode,
			fmt.Sprintf("POST %s %s failed %s", hostName, path, response.Body)}
		return "", err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var versionMap map[string]string
	err = json.Unmarshal(responseBody, &versionMap)
	if err != nil {
		return "", err
	}

	return versionMap["version_identifier"], nil
}

// RetrieveParams contains optional arguments to Key.Retrieve
type RetrieveParams struct {
	VersionID       string
	SliceOffset     int
	SliceSize       int
	ModifiedSince   interface{}
	UnmodifiedSince interface{}
}

// Retrieve fetches data from nimbus.io
func Retrieve(requester Requester, collectionName string, key string,
	retrieveParams RetrieveParams) (io.ReadCloser, error) {

	if retrieveParams.ModifiedSince != nil {
		return nil, fmt.Errorf("not implemented: retrieveParams.ModifiedSince")
	}
	if retrieveParams.UnmodifiedSince != nil {
		return nil, fmt.Errorf("not implemented: retrieveParams.UnmodifiedSince")
	}

	method := "GET"
	hostName := requester.CollectionHostName(collectionName)

	var path string
	if retrieveParams.VersionID == "" {
		path = fmt.Sprintf("/data/%s", url.QueryEscape(key))
	} else {
		values := url.Values{}
		values.Add("version_identifier", retrieveParams.VersionID)
		path = fmt.Sprintf("/data/%s?%s", url.QueryEscape(key), values.Encode())
	}

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return nil, err
	}

	successfulStatusCode := http.StatusOK

	if retrieveParams.SliceOffset > 0 || retrieveParams.SliceSize > 0 {
		var rangeArg string

		if retrieveParams.SliceSize > 0 {
			rangeArg = fmt.Sprintf("bytes=%d-%d", retrieveParams.SliceOffset,
				retrieveParams.SliceOffset+retrieveParams.SliceSize-1)
		} else {
			rangeArg = fmt.Sprintf("bytes=%d-", retrieveParams.SliceOffset)
		}

		request.Header.Add("range", rangeArg)
		successfulStatusCode = http.StatusPartialContent
	}

	response, err := requester.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != successfulStatusCode {
		err = HTTPError{response.StatusCode,
			fmt.Sprintf("GET %s %s failed %s", hostName, path, response.Body)}
		return nil, err
	}

	if retrieveParams.SliceSize > 0 && response.ContentLength != -1 &&
		response.ContentLength != int64(retrieveParams.SliceSize) {
		return nil, fmt.Errorf("content length mismatch: expected %d found %d",
			retrieveParams.SliceSize, response.ContentLength)
	}

	return response.Body, nil
}

// DeleteKey deletes a Key, hiding all versions
func DeleteKey(requester Requester, collectionName string, key string) error {
	return DeleteVersion(requester, collectionName, key, "")
}

// DeleteVersion deletes a specific version of a Key
func DeleteVersion(requester Requester, collectionName string, key string,
	versionIdentifier string) error {
	method := "DELETE"
	hostName := requester.CollectionHostName(collectionName)

	var path string
	if versionIdentifier == "" {
		path = fmt.Sprintf("/data/%s", url.QueryEscape(key))
	} else {
		values := url.Values{}
		values.Add("version", versionIdentifier)
		path = fmt.Sprintf("/data/%s?%s", url.QueryEscape(key), values.Encode())
	}

	request, err := requester.CreateRequest(method, hostName, path, nil)
	if err != nil {
		return err
	}

	response, err := requester.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = HTTPError{response.StatusCode,
			fmt.Sprintf("DELETE %s %s failed %s", hostName, path, response.Body)}
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var resultMap map[string]bool
	err = json.Unmarshal(responseBody, &resultMap)
	if err != nil {
		return err
	}
	if !resultMap["success"] {
		err = fmt.Errorf("unexpected 'false' for 'success'")
		return err
	}

	return nil
}
