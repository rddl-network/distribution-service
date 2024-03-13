package service

import (
	"context"
	"fmt"
	"net/http"
)

type IShamirClient interface {
	IssueTransaction(amount uint64, address string) (err error)
}

type ShamirClient struct {
	host string
}

func NewShamirClient(host string) *ShamirClient {
	return &ShamirClient{host: host}
}

func (sc *ShamirClient) IssueTransaction(amount uint64, address string) (err error) {
	url := fmt.Sprintf("%s/%s/%d", sc.host, address, amount)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, nil)
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
