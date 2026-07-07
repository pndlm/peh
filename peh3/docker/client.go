package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

// now including our own simple docker client since the official stuff
// was a little bit difficult to deal with

type Client struct {
	HttpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", "/var/run/docker.sock")
				},
			},
		},
	}
}

type ApiError struct {
	Message string `json:"message"`
}

func getError(resp *http.Response) error {
	var apiErr ApiError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil || apiErr.Message == "" {
		return fmt.Errorf("Unknown error status %d", resp.StatusCode)
	}
	return errors.New(apiErr.Message)
}

type Container struct {
	ID     string            `json:"Id"`
	State  string            `json:"State"`
	Labels map[string]string `json:"Labels"`
}

type ContainerListOptions struct {
	All bool
}

func (dc *Client) ContainerList(opt ContainerListOptions) ([]Container, error) {
	query := url.Values{}
	query.Set("all", strconv.FormatBool(opt.All))
	req, _ := http.NewRequest("GET", "http://docker/containers/json?"+query.Encode(), nil)
	resp, err := dc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, getError(resp)
	}
	var containers []Container
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, err
	}
	return containers, nil
}

func (dc *Client) ContainerStop(id string) error {
	req, _ := http.NewRequest("POST", "http://docker/containers/"+id+"/stop", nil)
	resp, err := dc.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	if resp.StatusCode >= 300 {
		return getError(resp)
	}
	return nil
}

func (dc *Client) ContainerRemove(id string) error {
	req, _ := http.NewRequest("DELETE", "http://docker/containers/"+id, nil)
	resp, err := dc.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return getError(resp)
	}
	return nil
}
