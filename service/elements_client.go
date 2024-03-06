package service

type IElementsClient interface {
}

type ElementsClient struct{}

func NewElementsClient() *ElementsClient {
	return &ElementsClient{}
}
