package stardog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ConnectionDetails struct {
	Endpoint string
	Username string
	Password string
}

func NewConnectionDetails(endpoint string, username string, password string) *ConnectionDetails {
	connectionDetails := ConnectionDetails{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	}
	return &connectionDetails
}

type Permission struct {
	Action       string   `json:"action"`
	ResourceType string   `json:"resource_type"`
	Resource     []string `json:"resource"`
}

func NewPermission(action string, resource_type string, resource []string) *Permission {
	permission := Permission{
		Action:       action,
		ResourceType: resource_type,
		Resource:     resource,
	}
	return &permission
}

func Alive(connectionDetails ConnectionDetails) bool {
	url := connectionDetails.Endpoint + "/admin/alive"
	request, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		log.Printf("Could not make request %v", err)
	}
	request.Header.Add("Accept", "application/json")
	request.SetBasicAuth(connectionDetails.Username, connectionDetails.Password)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request: %v", err)
	}

	return response.StatusCode == 200
}

type User struct {
	Name string
}

type getUsersResponse struct {
	Users []string `json:"users"`
}

func GetUsers(connectionDetails ConnectionDetails) []User {
	url := connectionDetails.Endpoint + "/admin/users"
	request, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		log.Printf("Could not make request %v", err)
	}
	request.Header.Add("Accept", "application/json")
	request.SetBasicAuth(connectionDetails.Username, connectionDetails.Password)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request: %v", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body - %v", err)
	}
	var result getUsersResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	var users []User
	for _, rec := range result.Users {
		newUser := User{
			Name: rec,
		}
		users = append(users, newUser)
	}
	return users
}

type GetUserDetailsResponse struct {
	Enabled     bool     `json:"enabled"`
	Superuser   bool     `json:"superuser"`
	Roles       []string `json:"roles"`
	Permissions []struct {
		Action       string   `json:"action"`
		ResourceType string   `json:"resource_type"`
		Resource     []string `json:"resource"`
		Explicit     bool     `json:"explicit"`
	} `json:"permissions"`
}

func GetUserDetails(connectionDetails ConnectionDetails, user User) GetUserDetailsResponse {
	url := connectionDetails.Endpoint + "/admin/users/" + user.Name
	request, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		log.Printf("Could not make request %v", err)
	}
	request.Header.Add("Accept", "application/json")
	request.SetBasicAuth(connectionDetails.Username, connectionDetails.Password)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request: %v", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body - %v", err)
	}
	var result GetUserDetailsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	return result
}

func DeleteUserPermission(connectionDetails ConnectionDetails, user User, permission Permission) bool {
	url := connectionDetails.Endpoint + "/admin/permissions/user/" + user.Name + "/delete"
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(permission)

	req, err := http.NewRequest(http.MethodPost, url, payloadBuf)
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.SetBasicAuth(connectionDetails.Username, connectionDetails.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return true
	} else {
		return false
	}
}

func AddUserPermission(connectionDetails ConnectionDetails, user User, permission Permission) bool {
	url := connectionDetails.Endpoint + "/admin/permissions/user/" + user.Name
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(permission)

	req, err := http.NewRequest(http.MethodPut, url, payloadBuf)
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.SetBasicAuth(connectionDetails.Username, connectionDetails.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return true
	} else {
		return false
	}
}

type Credentials struct {
	Username string   `json:"username"`
	Password []string `json:"password"`
}

func AddUser(connectionDetails ConnectionDetails, credentials Credentials) bool {
	url := connectionDetails.Endpoint + "/admin/users"
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(credentials)

	req, err := http.NewRequest(http.MethodPost, url, payloadBuf)
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.SetBasicAuth(connectionDetails.Username, connectionDetails.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return true
	} else {
		return false
	}
}

func DeleteUser(connectionDetails ConnectionDetails, user User) bool {
	url := connectionDetails.Endpoint + "/admin/users/" + user.Name

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return false
	}
	req.SetBasicAuth(connectionDetails.Username, connectionDetails.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return true
	} else {
		return false
	}
}

type Password struct {
	Password string `json:"password"`
}

func NewPassword(password string) Password {
	return Password{Password: password}
}

func NewUser(username string) User {
	return User{Name: username}
}

func ChangeUserPassword(connectionDetails ConnectionDetails, user User, password Password) bool {
	url := connectionDetails.Endpoint + "/admin/users/" + user.Name + "/pwd"
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(password)

	req, err := http.NewRequest(http.MethodPut, url, payloadBuf)
	if err != nil {
		return false
	}
	req.SetBasicAuth(connectionDetails.Username, connectionDetails.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}
