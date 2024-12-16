package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/planetmint/planetmint-go/util"
	"github.com/rddl-network/distribution-service/config"
)

// Distributes 10% of received funds to all validators
func (ds *DistributionService) DistributeToValidators() {
	cfg := config.GetConfig()

	distributionAmt, err := ds.getDistributionAmount()
	if err != nil {
		ds.logger.Error("msg", "Error while calculating distribution amount: "+err.Error())
		return
	}

	if distributionAmt == 0 {
		ds.logger.Error("msg", "No tokens to distribute.")
		return
	}

	liquidAddresses, err := ds.getBeneficiaries()
	if err != nil {
		ds.logger.Error("msg", "Error while fetching beneficiary addresses: "+err.Error())
		return
	}

	// CalculateShares
	share, _ := ds.calculateShares(distributionAmt, uint64(len(liquidAddresses)))

	// SendToAddresses
	ds.logger.Info("msg", "sending tokens", "addresses", strings.Join(liquidAddresses, ","), "amount", distributionAmt, "share", share)
	err = ds.sendToAddresses(liquidAddresses, share, cfg.Asset)
	if err != nil {
		ds.logger.Error("msg", "Error while sending to validators: "+err.Error())
		return
	}
}

func (ds *DistributionService) getDistributionAmount() (distributionAmt uint64, err error) {
	received, err := ds.checkReceivedBalance()
	if err != nil {
		return
	}

	ds.logger.Debug("msg", "Reading last occurrence")
	occurrence, err := ds.GetLastOccurrence()
	if err != nil {
		return
	}

	ds.logger.Debug("msg", "Storing current occurrence")
	err = ds.StoreOccurrence(time.Now().Unix(), received)
	if err != nil {
		return
	}

	if occurrence == nil {
		return CalculateValidatorDistributionAmount(0, received), nil
	}

	return CalculateValidatorDistributionAmount(occurrence.Amount, received), nil
}

// Checks for received asset on a given address
func (ds *DistributionService) checkReceivedBalance() (received uint64, err error) {
	cfg := config.GetConfig()
	ds.logger.Info("msg", "checking received balance", "address", cfg.FundAddress, " asset", cfg.Asset)

	confirmationString := strconv.Itoa(cfg.Confirmations)
	txDetails, err := ds.eClient.ListReceivedByAddress(cfg.GetElementsURL(),
		[]string{confirmationString, "false", "true", `"` + cfg.FundAddress + `"`, `"` + cfg.Asset + `"`},
	)
	if err != nil {
		return 0, err
	}

	for _, txDetail := range txDetails {
		received += util.RDDLToken2Uint(txDetail.Amount)
	}

	return
}

func (ds *DistributionService) getBeneficiaries() (addresses []string, err error) {
	plmntAddresses, err := ds.getActiveValidatorAddresses()
	if err != nil {
		return nil, err
	}

	ds.logger.Info("msg", "fetching liquid receive addresses", "planetmintAddresses", strings.Join(plmntAddresses, ","))
	return ds.getReceiveAddresses(plmntAddresses)
}

// getReceiveAddresses fetches receive addresses from the rddl-2-plmnt service
func (ds *DistributionService) getReceiveAddresses(addresses []string) (receiveAddresses []string, err error) {
	for _, address := range addresses {
		receiveAddress, err := ds.r2pClient.GetReceiveAddress(context.Background(), address)
		if err != nil {
			return nil, err
		}
		receiveAddresses = append(receiveAddresses, receiveAddress.LiquidAddress)
	}
	return
}

// Gets all active validator addresses
func (ds *DistributionService) getActiveValidatorAddresses() (addresses []string, err error) {
	valAddresses, err := ds.pmClient.GetValidatorAddresses()
	if err != nil {
		return nil, err
	}

	for _, address := range valAddresses {
		delegationAddresses, err := ds.pmClient.GetValidatorDelegationAddresses(address)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, delegationAddresses...)
	}

	return
}

// Calculates share per given address
func (ds *DistributionService) calculateShares(total uint64, numValidators uint64) (share uint64, remainder uint64) {
	if numValidators == 0 {
		return 0, total
	}

	share = total / numValidators
	remainder = total % numValidators
	return
}

func CalculateValidatorDistributionAmount(prev uint64, curr uint64) (distributionAmt uint64) {
	if prev == 0 {
		return curr / 100 * 10
	}

	return (curr - prev) / 100 * 10
}
