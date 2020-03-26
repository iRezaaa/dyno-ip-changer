package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"bytes"
)

func getDomainList(client *http.Client, apiKey string) (*DnsRequestResponse, error) {
	req, err := http.NewRequest("GET", "https://api.dynu.com/v2/dns", nil)

	if err != nil {
		return nil, &CustomError{
			Unknown,
			"Error while create new request!!",
		}
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("API-Key", apiKey)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("")
		return nil, &CustomError{
			SendRequest,
			"Error while send request to Dynu server, check your internet connection!",
		}
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)

		var responseModel DnsRequestResponse
		err := json.Unmarshal(respBody, &responseModel)

		if err != nil {
			return nil, &CustomError{
				Unmarshal,
				"Error while unmarshal response!",
			}
		}

		return &responseModel, nil
	} else if resp.StatusCode == 401 {
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)

		var responseModel ErrorResponse
		err := json.Unmarshal(respBody, &responseModel)

		if err != nil {
			return nil, &CustomError{
				Unmarshal,
				"Error while unmarshal error response!",
			}
		}

		return nil, &CustomError{
			Auth,
			fmt.Sprintf("Authentication and/or authorized has failed. server message {statusCode : statusCode : %d,type : %s, message : %s}", responseModel.StatusCode, responseModel.Type, responseModel.Message),
		}
	} else if resp.StatusCode == 500 {
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)

		var responseModel ErrorResponse
		err := json.Unmarshal(respBody, &responseModel)

		if err != nil {
			return nil, &CustomError{
				Unmarshal,
				"Error while unmarshal error response!",
			}
		}

		return nil, &CustomError{
			Response,
			fmt.Sprintf("Server send error message as response {statusCode : statusCode : %d,type : %s, message : %s}", responseModel.StatusCode, responseModel.Type, responseModel.Message),
		}
	} else {
		return nil, &CustomError{
			Unknown,
			fmt.Sprintf("Unknown error response!"),
		}
	}

	return nil, nil
}
func updateDomainIP(client *http.Client, apiKey string, currentData *Domain, newIP string) error {
	requestBody := &UpdateDomainIPRequest{
		Name:              currentData.Name,
		Group:             currentData.Group,
		Ipv4Address:       newIP,
		TTL:               currentData.TTL,
		Ipv4:              currentData.Ipv4,
		Ipv6:              currentData.Ipv6,
		Ipv4WildcardAlias: currentData.Ipv4WildcardAlias,
		Ipv6WildcardAlias: false,
		AllowZoneTransfer: false,
		DnsSec:            false,
	}
	jsonRequest, err := json.Marshal(requestBody)

	if err != nil {
		return &CustomError{
			Marshal,
			"Marshal json request",
		}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.dynu.com/v2/dns/%d", currentData.ID), bytes.NewBuffer(jsonRequest))

	if err != nil {
		return &CustomError{
			Unknown,
			"Error while creating request!",
		}
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("API-Key", apiKey)

	resp, err := client.Do(req)

	if err != nil {
		return &CustomError{
			SendRequest,
			"Error while send request to server!",
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}else if resp.StatusCode == 401 {
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)

		var responseModel ErrorResponse
		err := json.Unmarshal(respBody, &responseModel)

		if err != nil {
			return &CustomError{
				Unmarshal,
				"Error while unmarshal error response!",
			}
		}

		return &CustomError{
			Auth,
			fmt.Sprintf("Authentication and/or authorized has failed. server message {statusCode : statusCode : %d,type : %s, message : %s}", responseModel.StatusCode, responseModel.Type, responseModel.Message),
		}
	} else if resp.StatusCode == 500 || resp.StatusCode == 501 || resp.StatusCode == 502 {
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)

		var responseModel ErrorResponse
		err := json.Unmarshal(respBody, &responseModel)

		if err != nil {
			return &CustomError{
				Unmarshal,
				"Error while unmarshal error response!",
			}
		}

		return &CustomError{
			Response,
			fmt.Sprintf("Server send error message as response {statusCode : statusCode : %d,type : %s, message : %s}", responseModel.StatusCode, responseModel.Type, responseModel.Message),
		}
	} else {
		return &CustomError{
			Unknown,
			fmt.Sprintf("Unknown error response! status code is %d}",resp.StatusCode),
		}
	}

	return nil
}
