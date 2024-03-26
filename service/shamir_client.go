package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type IShamirClient interface {
	IssueTransaction(amount string, address string) (err error)
}

type ShamirClient struct {
	host string
}

func NewShamirClient(host string) *ShamirClient {
	return &ShamirClient{host: host}
}

type SendTokensRequest struct {
	Recipient string `json:"recipient"`
	Amount    string `json:"amount"`
}

func (sc *ShamirClient) IssueTransaction(amount string, address string) (err error) {
	url := sc.host + "/send"

	body := &SendTokensRequest{
		Recipient: address,
		Amount:    amount,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return
}
