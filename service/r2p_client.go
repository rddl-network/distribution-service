package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IR2PClient interface {
	GetReceiveAddress(plmntAddress string) (receiveAddress string, err error)
}

type R2PClient struct {
	host string
}

type AddressBody struct {
	LiquidAddress         string `binding:"required" json:"liquid-address"`
	PlanetmintBeneficiary string `binding:"required" json:"planetmint-beneficiary"`
}

func NewR2PClient(host string) *R2PClient {
	return &R2PClient{host: host}
}

func (r2p *R2PClient) GetReceiveAddress(plmntAddress string) (receiveAddress string, err error) {
	url := fmt.Sprintf("%s/receiveaddress/%s", r2p.host, plmntAddress)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return "", nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var addressBody AddressBody
	err = json.Unmarshal(body, &addressBody)
	if err != nil {
		return "", err
	}

	return addressBody.LiquidAddress, nil
}
