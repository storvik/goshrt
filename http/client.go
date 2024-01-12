package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/storvik/goshrt"
)

const errSendReq string = "could not send request, "
const errInvalidStatus string = "received invalid statuscode from endpoint, "
const errDecodingResp string = "error decoding response from endpoint, "
const errMarshal string = "error marshal request, "

type Client struct {
	Address string
	Key     string
}

func (c *Client) ShrtAdd(s *goshrt.Shrt) error {
	postBody, err := json.Marshal(s)
	if err != nil {
		return errors.New(errMarshal + err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, c.Address+"/api/shrt", bytes.NewReader(postBody))
	if err != nil {
		return errors.New("could not create new request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Key)

	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		return errors.New(errSendReq + err.Error())
	}
	defer res.Body.Close()

	// Read response
	if res.StatusCode != http.StatusCreated {
		return errors.New("received invalid statuscode from endpoint")
	}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&s)
	if err != nil {
		return errors.New("error decoding response from endpoint")
	}

	return nil
}

func (c *Client) ShrtGet(s *goshrt.Shrt) error {
	postBody, err := json.Marshal(s)
	if err != nil {
		return errors.New(errMarshal + err.Error())
	}

	var slug string
	if s.ID > 0 {
		slug = fmt.Sprintf("/api/shrt/%d", s.ID)
	} else {
		slug = fmt.Sprintf("/api/shrt/%s/%s", s.Domain, s.Slug)
	}

	req, err := http.NewRequest(http.MethodGet, c.Address+slug, bytes.NewReader(postBody))
	if err != nil {
		return errors.New("could not create new request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Key)

	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		return errors.New(errSendReq + err.Error())
	}
	defer res.Body.Close()

	// Read response
	if res.StatusCode != http.StatusOK {
		return errors.New(errInvalidStatus + strconv.Itoa(res.StatusCode))
	}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&s)
	if err != nil {
		return errors.New(errDecodingResp + err.Error())
	}

	return nil
}

func (c *Client) ShrtDelete(s *goshrt.Shrt) error {
	postBody, err := json.Marshal(s)
	if err != nil {
		return errors.New(errMarshal + err.Error())
	}

	slug := fmt.Sprintf("/api/shrt/%d", s.ID)

	req, err := http.NewRequest(http.MethodDelete, c.Address+slug, bytes.NewReader(postBody))
	if err != nil {
		return errors.New("could not create new request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Key)

	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		return errors.New("could not send request, " + err.Error())
	}
	defer res.Body.Close()

	// Read response
	if res.StatusCode != http.StatusOK {
		return errors.New(errInvalidStatus + strconv.Itoa(res.StatusCode))
	}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&s)
	if err != nil {
		return errors.New(errDecodingResp + err.Error())
	}

	return nil
}

func (c *Client) ShrtGetList(d string) ([]*goshrt.Shrt, error) {
	var slug string
	if d == "" {
		slug = "/api/shrts"
	} else {
		slug = fmt.Sprintf("/api/shrts/%s", d)
	}

	req, err := http.NewRequest(http.MethodGet, c.Address+slug, http.NoBody)
	if err != nil {
		return nil, errors.New("could not create new request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Key)

	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New(errSendReq + err.Error())
	}
	defer res.Body.Close()

	// Read response
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(errInvalidStatus + strconv.Itoa(res.StatusCode))
	}

	var shrts []*goshrt.Shrt
	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&shrts)
	if err != nil {
		return nil, errors.New(errDecodingResp + err.Error())
	}

	return shrts, nil
}
