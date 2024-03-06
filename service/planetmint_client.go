package service

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IPlanetmintClient interface {
	GetValidatorAddresses() (addresses []string, err error)
	GetValidatorDelegationAddresses(validatorAddress string) (addresses []string, err error)
}

type PlanetmintClient struct {
	host string
}

func NewPlanetmintClient(host string) *PlanetmintClient {
	return &PlanetmintClient{host: host}
}

func (pmc *PlanetmintClient) GetValidatorAddresses() (addresses []string, err error) {
	client, err := pmc.getStakingQueryClient()
	if err != nil {
		return nil, err
	}

	validatorsRes, err := client.Validators(context.Background(), &stakingTypes.QueryValidatorsRequest{})
	if err != nil {
		return nil, err
	}

	for _, validator := range validatorsRes.GetValidators() {
		addresses = append(addresses, validator.OperatorAddress)
	}

	return
}

func (pmc *PlanetmintClient) GetValidatorDelegationAddresses(validatorAddress string) (addresses []string, err error) {
	client, err := pmc.getStakingQueryClient()
	if err != nil {
		return nil, err
	}

	delegationRes, err := client.ValidatorDelegations(context.Background(), &stakingTypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: validatorAddress,
	})
	if err != nil {
		return nil, err
	}

	for _, delegation := range delegationRes.GetDelegationResponses() {
		addresses = append(addresses, delegation.Delegation.DelegatorAddress)
	}

	return
}

func (pmc *PlanetmintClient) getStakingQueryClient() (stakingClient stakingTypes.QueryClient, err error) {
	grpcConn, err := pmc.getGRPCConnection()
	if err != nil {
		return nil, err
	}
	stakingClient = stakingTypes.NewQueryClient(grpcConn)
	return
}

func (pmc *PlanetmintClient) getGRPCConnection() (grpcConn *grpc.ClientConn, err error) {
	return grpc.Dial(
		pmc.host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
}
