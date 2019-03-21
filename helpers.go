package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SbAuthReply struct {
	Status bool   `json:"Status"`
	Msg    string `json:"Msg"`
	// User   User   `json:"User"`
}

// type User struct {
// 	Name        string `json:"Name"`
// 	Pw          string `json:"Pw"`
// 	IsSuperuser bool   `json:"IsSuperuser"`
// }

// NOTE:
// the commented out above part is currently not needed

func isAuth(token string) bool {
	// Validates a JWToken with the Sb Auth Service

	client := &http.Client{}
	url := fmt.Sprintf("%s%s", SbAuthServiceEndpointBaseUrl, "isauthenticated")
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("this panics 2")
	}
	req.Header.Add("token", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("this panics 3")
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("this panics 3.5")
		fmt.Println(err)
	}

	reply := SbAuthReply{}
	json.Unmarshal(body, &reply)

	return reply.Status
}
