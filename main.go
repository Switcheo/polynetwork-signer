package main

import (
	"fmt"
	"github.com/Switcheo/polynetwork-signer/app"
	config "github.com/Switcheo/polynetwork-signer/config"
	"github.com/ethereum/go-ethereum/common"
	polycommon "github.com/polynetwork/poly/common"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

const ArgsExpected = 4

func main() {
	if len(os.Args) != ArgsExpected + 1 {
		panic("Need " + strconv.Itoa(ArgsExpected) + " args")
	}

	// Setup

	err := godotenv.Load(".env")
	confileFilePath := os.Getenv("CONFIG_FILE_PATH")
	conf := config.NewConfig(confileFilePath)
	polySdk, err := app.SetUpPoly(conf)
	if err != nil {
		panic("SDK init error: " + err.Error())
	}
	polySigner, err := app.SetupPolySigner(conf, polySdk)
	if err != nil {
		panic("signer init error: " + err.Error())
	}

	// Parse chain id
	chainID, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		panic("unable to parse chain id")
	}

	// Parse data
	txData := common.Hex2Bytes(os.Args[2])

	// Parse height
	height, err := strconv.ParseUint(os.Args[3], 10, 32)
	if err != nil {
		panic("unable to parse height")
	}

	// Parse proof
	proof := common.Hex2Bytes(os.Args[4])

	// Get rest of params
	relayerAddress := common.Hex2Bytes(polySigner.Address.ToHexString())
	headerOrCrossChainMsg := []byte{}

	// Sign
	tx, err := app.SignImportOuterTransfer(*polySdk, chainID, txData, uint32(height), proof, relayerAddress, headerOrCrossChainMsg, polySigner)
	if err != nil {
		panic("Signing error: " + err.Error())
	}
	sink := polycommon.NewZeroCopySink(nil)
	err = tx.Serialization(sink)
	if err != nil {
		panic("Serialization error: " + err.Error())
	}

	// print
	fmt.Println(common.Bytes2Hex(sink.Bytes()))

	os.Exit(1)
}
