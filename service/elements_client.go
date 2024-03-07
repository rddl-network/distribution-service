package service

import (
	elementsrpc "github.com/rddl-network/elements-rpc"
	"github.com/rddl-network/elements-rpc/types"
)

type IElementsClient interface {
	ListReceivedByAddress(url string, params []string) (receivedTx []types.ListReceivedByAddressResult, err error)
}

type ElementsClient struct{}

func NewElementsClient() *ElementsClient {
	return &ElementsClient{}
}

func (ec *ElementsClient) ListReceivedByAddress(url string, params []string) (receivedTx []types.ListReceivedByAddressResult, err error) {
	return elementsrpc.ListReceivedByAddress(url, params)
}
