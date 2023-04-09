package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/storvik/goshrt"
)

func ppShrt(s *goshrt.Shrt) {
	fmt.Printf("ID\t\t%d\n", s.ID)
	fmt.Printf("Domain\t\t%s\n", s.Domain)
	fmt.Printf("Slug\t\t%s\n", s.Slug)
	fmt.Printf("Destination\t%s\n", s.Dest)
	fmt.Printf("Expiry\t\t%s\n", s.Expiry.Format("2006.02.01"))
}

func (a *AppConfig) shrtAdd(s *goshrt.Shrt) error {
	postBody, _ := json.Marshal(s)
	req, err := http.NewRequest("POST", a.Server.Address+"/api/shrt", bytes.NewReader(postBody))
	if err != nil {
		return errors.New("could not create new request")
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer secret")

	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return errors.New("could not send request, " + err.Error())
	}
	defer res.Body.Close()

	// Read response
	if res.StatusCode != http.StatusCreated {
		return errors.New("received invalid statuscode from endpoint")
	}
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }

	shrt := new(goshrt.Shrt)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&shrt)
	if err != nil {
		return errors.New("error decoding response from endpoint")
	}

	fmt.Printf("Successfully added shrt!\n")
	ppShrt(shrt)
	return nil
}

func (a *AppConfig) shrtGet(s *goshrt.Shrt) error {
	return nil
}
