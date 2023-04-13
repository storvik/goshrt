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

type Client struct {
	Address string
	Key     string
}

func (c *Client) ShrtAdd(s *goshrt.Shrt) error {
	postBody, _ := json.Marshal(s)
	req, err := http.NewRequest(http.MethodPost, c.Address+"/api/shrt", bytes.NewReader(postBody))
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
	if res.StatusCode != http.StatusCreated {
		return errors.New("received invalid statuscode from endpoint")
	}
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&s)
	if err != nil {
		return errors.New("error decoding response from endpoint")
	}

	return nil
}

func (c *Client) ShrtGet(s *goshrt.Shrt) error {
	postBody, _ := json.Marshal(s)
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
		return errors.New("could not send request, " + err.Error())
	}
	defer res.Body.Close()

	// Read response
	if res.StatusCode != http.StatusOK {
		return errors.New("received invalid statuscode from endpoint, " + strconv.Itoa(res.StatusCode))
	}
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&s)
	if err != nil {
		return errors.New("error decoding response from endpoint, " + err.Error())
	}

	return nil
}

func (c *Client) ShrtDelete(s *goshrt.Shrt) error {
	postBody, _ := json.Marshal(s)
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
		return errors.New("received invalid statuscode from endpoint, " + strconv.Itoa(res.StatusCode))
	}
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&s)
	if err != nil {
		return errors.New("error decoding response from endpoint, " + err.Error())
	}

	return nil
}

func (c *Client) ShrtGetList(d string) ([]*goshrt.Shrt, error) {
	var slug string
	if d == "" {
		slug = fmt.Sprintf("/api/shrts")
	} else {
		slug = fmt.Sprintf("/api/shrts/%s", d)
	}
	req, err := http.NewRequest(http.MethodGet, c.Address+slug, nil)
	if err != nil {
		return nil, errors.New("could not create new request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Key)

	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("could not send request, " + err.Error())
	}
	defer res.Body.Close()

	// Read response
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("received invalid statuscode from endpoint, " + strconv.Itoa(res.StatusCode))
	}
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }

	var shrts []*goshrt.Shrt
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&shrts)
	if err != nil {
		return nil, errors.New("error decoding response from endpoint, " + err.Error())
	}

	return shrts, nil
}
