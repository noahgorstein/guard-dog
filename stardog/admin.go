package stardog

import (
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
	admin := ConnectionDetails{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	}
	return &admin
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

type getUsersResponse struct {
	Users []string `json:"users"`
}

type User struct {
	Name string
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

// func DeleteUserPermission(connectionDetails ConnectionDetails, user User) {
// 	url := connectionDetails.Endpoint + "/admin/permissions/user/" + user.Name + "/delete"
// 	request, err := http.NewRequest(
// 		http.MethodGet,
// 		url,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Printf("Could not make request %v", err)
// 	}
// 	request.Header.Add("Accept", "application/json")
// 	request.SetBasicAuth(connectionDetails.Username, connectionDetails.Password)

// 	response, err := http.DefaultClient.Post(request)
// 	if err != nil {
// 		log.Printf("Could not make a request: %v", err)
// 	}

// 	body, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		log.Printf("Could not read response body - %v", err)
// 	}
// }
