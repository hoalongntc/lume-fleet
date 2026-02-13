package lume

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client talks to the Lume HTTP API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:7777"
	}
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) ListVMs() ([]VM, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/lume/vms")
	if err != nil {
		return nil, fmt.Errorf("lume: list VMs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lume: list VMs: status %d: %s", resp.StatusCode, body)
	}

	var vms []VM
	if err := json.NewDecoder(resp.Body).Decode(&vms); err != nil {
		return nil, fmt.Errorf("lume: list VMs: decode: %w", err)
	}
	return vms, nil
}

func (c *Client) GetVM(name string) (*VM, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/lume/vms/" + name)
	if err != nil {
		return nil, fmt.Errorf("lume: get VM %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lume: get VM %q: status %d: %s", name, resp.StatusCode, body)
	}

	var vm VM
	if err := json.NewDecoder(resp.Body).Decode(&vm); err != nil {
		return nil, fmt.Errorf("lume: get VM %q: decode: %w", name, err)
	}
	return &vm, nil
}

func (c *Client) CreateVM(req CreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("lume: create VM: marshal: %w", err)
	}

	// Use a longer timeout for create since it's async and the server may take time to respond.
	createClient := &http.Client{Timeout: 5 * time.Minute}
	resp, err := createClient.Post(c.BaseURL+"/lume/vms", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("lume: create VM %q: %w", req.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lume: create VM %q: status %d: %s", req.Name, resp.StatusCode, respBody)
	}
	return nil
}

func (c *Client) RunVM(name string, req RunRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("lume: run VM: marshal: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/lume/vms/"+name+"/run", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("lume: run VM %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lume: run VM %q: status %d: %s", name, resp.StatusCode, respBody)
	}
	return nil
}

func (c *Client) StopVM(name string) error {
	resp, err := c.HTTPClient.Post(c.BaseURL+"/lume/vms/"+name+"/stop", "application/json", nil)
	if err != nil {
		return fmt.Errorf("lume: stop VM %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lume: stop VM %q: status %d: %s", name, resp.StatusCode, body)
	}
	return nil
}

// WaitForCreation polls until the VM is no longer provisioning.
func (c *Client) WaitForCreation(name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		vm, err := c.GetVM(name)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if vm.ProvisioningOperation == nil && vm.Status == "stopped" {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("lume: VM %q creation timed out after %v", name, timeout)
}
