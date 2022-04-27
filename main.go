package main

import (
	"fmt"
	"github.com/Switcheo/polynetwork-signer/app"
	config "github.com/Switcheo/polynetwork-signer/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"os"
)

func main() {
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

	cmd := os.Args[1]
	userArgs := os.Args[2:]

	var signedTxBytes []byte

	switch cmd {
	case "create_crosschain_tx":
		signedTxBytes = app.CreateCrossChainTx(polySdk, polySigner, userArgs)
	case "create_crosschain_tx_tendermint":
		signedTxBytes = app.CreateCrossChainTxTendermint(polySdk, polySigner, userArgs)
	}

	hexTx := common.Bytes2Hex(signedTxBytes)

	// print
	fmt.Printf(hexTx)

	os.Exit(0)
}
