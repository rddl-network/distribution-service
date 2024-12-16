package service

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/rddl-network/distribution-service/config"
)

const (
	lastBlockHeightFileName = "lastBlockHeight.dat"
)

func (ds *DistributionService) DistributeToAdvisories() {
	// gather data
	currentBlockHeight, err := ds.pmClient.GetBlockHeight()
	if err != nil {
		ds.logger.Error("error querying block height:", err.Error())
		return
	}
	lastWrittenBlockHeight, err := ds.ReadLastBlockHeight()
	if err != nil {
		ds.logger.Error("error reading block height from file", err.Error())
		return
	}

	runDistribution := ds.RunDistribution(currentBlockHeight, lastWrittenBlockHeight)
	if !runDistribution {
		ds.logger.Debug("don't run distribution")
		return
	}

	err = ds.DistributeToAdvisoriesOnce()
	if err != nil {
		return
	}

	err = ds.WriteLastBlockHeight(currentBlockHeight)
	if err != nil {
		ds.logger.Error("error writing to block height file", err.Error())
		return
	}
}

func (ds *DistributionService) DistributeToAdvisoriesOnce() (err error) {
	cfg := config.GetConfig()
	distributions := config.GetWeeklyAdvisoryDistribution()
	for _, distribution := range distributions {
		address := distribution.Address
		if cfg.TestnetMode {
			address = cfg.TestnetAddress
		}
		err = ds.sendToAddress(address, distribution.Amount, cfg.Asset)
		if err != nil {
			err = fmt.Errorf("sending to address failed: %w", err)
			return
		}
	}
	return
}

func (ds *DistributionService) RunDistribution(currentBlockHeight int64, lastWrittenBlockHeight int64) (run bool) {
	var blocksPerDay int64 = 3600
	blocksPerWeek := blocksPerDay * 7
	var distributionOffset int64 = 75
	var distributionSettlementOffset int64 = 5

	barrierToPass := lastWrittenBlockHeight + blocksPerWeek + distributionOffset + distributionSettlementOffset

	if currentBlockHeight >= barrierToPass {
		run = true
	}

	return
}

func (ds *DistributionService) ReadLastBlockHeight() (blockHeight int64, err error) {
	cfg := config.GetConfig()
	filePath := cfg.DataPath + lastBlockHeightFileName

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("error opening file: %w", err)
		return
	}
	defer file.Close()

	// Read the int64 value
	err = binary.Read(file, binary.LittleEndian, &blockHeight)
	if err != nil {
		err = fmt.Errorf("error reading block height: %w", err)
		return
	}

	return
}

func (ds *DistributionService) WriteLastBlockHeight(blockHeight int64) error {
	cfg := config.GetConfig()
	err := os.MkdirAll(cfg.DataPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating data directory: %w", err)
	}

	// Construct the full path to the file
	filePath := cfg.DataPath + lastBlockHeightFileName

	// Open the file with write-only and create/truncate flags
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Write the int64 value
	err = binary.Write(file, binary.LittleEndian, blockHeight)
	if err != nil {
		return fmt.Errorf("error writing block height: %w", err)
	}

	return nil
}
