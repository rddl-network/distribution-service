package service_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	shamir "github.com/rddl-network/shamir-coordinator-service/types"
	"github.com/stretchr/testify/assert"
)

func TestAdvisoryDistribution(t *testing.T) {
	app, db, mocks := setupService(t)
	defer db.Close()

	mocks.shamirClientMock.EXPECT().SendTokens(gomock.Any(), "VJLHinV6iAVSw7Mwx1yRe3jLT86pbjoygJRPNroXhcKxAmtX2EZZM4wAhW993umuquWG7wujcPXw98f9", "9615.38461540", gomock.Any()).Times(1).Return(shamir.SendTokensResponse{}, nil)

	app.DistributeToAdvisories()
}

func TestReadWriteLastBlockHeight(t *testing.T) {
	app, db, _ := setupService(t)
	defer db.Close()

	_, err := app.ReadLastBlockHeight()
	assert.NoError(t, err)

	var testValue int64 = 20
	err = app.WriteLastBlockHeight(testValue)
	assert.NoError(t, err)

	height, err := app.ReadLastBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, testValue, height)
}

type distRequest struct {
	lastWrittenBlock   int64
	currentBlockHeight int64
	run                bool
}

func TestRunDistribution(t *testing.T) {
	app, db, _ := setupService(t)
	defer db.Close()
	distRequests := []distRequest{
		{0, 25280, true},
		{0, 50560, true},
		{0, 25200, false},
		{25280, 25280, false},
		{25280, 50560, true},
		{25280, 75840, true},
		{25280, 50500, false},
		{25280, 25200, false},
	}
	for i, request := range distRequests {
		fmt.Printf("Index: %d, Request: %+v\n", i, request)
		result := app.RunDistribution(request.currentBlockHeight, request.lastWrittenBlock)
		assert.Equal(t, request.run, result)
	}
}
