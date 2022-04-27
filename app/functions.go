package app

import (
	"context"
	ccmkeeper "github.com/Switcheo/polynetwork-cosmos/x/ccm/keeper"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/polynetwork/poly-go-sdk"
	polycommon "github.com/polynetwork/poly/common"
	ccmcommon "github.com/polynetwork/poly/native/service/cross_chain_manager/common"
	ccmcosmos "github.com/polynetwork/poly/native/service/cross_chain_manager/cosmos"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"strconv"
)

const DefaultArgsLength = 4
const TendermintArgsLength = 4

// CreateCrossChainTx Forms the correct signed tx to be broadcasted by taking in the user args
// userArgs: <chainID> <txData> <height> <proof>
func CreateCrossChainTx(polySdk *sdk.PolySdk, polySigner *sdk.Account, userArgs []string) []byte {
	EnsureArgsLength(userArgs, DefaultArgsLength)
	// Parse chain id
	chainID, err := strconv.ParseUint(userArgs[0], 10, 64)
	if err != nil {
		panic("unable to parse chain id")
	}

	// Parse data
	txData := common.Hex2Bytes(userArgs[1])

	// Parse height
	height, err := strconv.ParseUint(userArgs[2], 10, 32)
	if err != nil {
		panic("unable to parse height")
	}

	// Parse proof
	proof := common.Hex2Bytes(userArgs[3])

	// Get rest of params
	relayerAddress := common.Hex2Bytes(polySigner.Address.ToHexString())
	headerOrCrossChainMsg := []byte{}

	// Sign
	tx, err := SignImportOuterTransfer(*polySdk, chainID, txData, uint32(height), proof, relayerAddress, headerOrCrossChainMsg, polySigner)
	if err != nil {
		panic("Signing error: " + err.Error())
	}
	sink := polycommon.NewZeroCopySink(nil)
	err = tx.Serialization(sink)
	if err != nil {
		panic("Serialization error: " + err.Error())
	}
	return sink.Bytes()
}

// CreateCrossChainTxTendermint Forms the correct signed tx to be broadcasted by querying provided rpc URL for the details
// userArgs: <rpcURL> <chainID> <height> <ccmKeyHash>
func CreateCrossChainTxTendermint(polySdk *sdk.PolySdk, polySigner *sdk.Account, userArgs []string) []byte {
	EnsureArgsLength(userArgs, TendermintArgsLength)

	// parse rpc
	rpcEndpoint := userArgs[0]

	// Parse chain id
	chainID, err := strconv.ParseUint(userArgs[1], 10, 64)
	if err != nil {
		panic("unable to parse chain id")
	}

	// parse header height
	headerHeight, err := strconv.ParseInt(userArgs[2], 10, 64)
	if err != nil {
		panic("unable to parse header height")
	}

	// initialize rpc client
	rpcClient, err := rpchttp.New(rpcEndpoint, "/websocket")
	if err != nil {
		panic("failed to init rpc client:" + err.Error())
	}

	// Get Header
	header, err := GetCosmosHeader(rpcClient, headerHeight)
	if err != nil {
		panic("unable to get header from rpc" + rpcEndpoint)
	}
	cdc := NewCodecForRelayer()
	headerOrCrossChainMsg, err := cdc.MarshalBinaryBare(*header)
	if err != nil {
		panic(err)
	}

	// parse ccm key hash
	ccmKeyHash := common.Hex2Bytes(userArgs[3])
	// Get cosmos proof value
	proofHeight := headerHeight - 1
	res, err := rpcClient.ABCIQueryWithOptions(context.Background(), "/store/ccm/key", ccmkeeper.GetCrossChainTxKey(ccmKeyHash),
		client.ABCIQueryOptions{Prove: true, Height: proofHeight})
	if err != nil {
		panic(err)
	}
	proof, err := res.Response.GetProofOps().Marshal()
	if err != nil {
		panic(err)
	}

	kp := merkle.KeyPath{}
	kp = kp.AppendKey([]byte("ccm"), merkle.KeyEncodingURL)
	kp = kp.AppendKey(res.Response.Key, merkle.KeyEncodingURL)
	proofValue, _ := cdc.MarshalBinaryBare(&ccmcosmos.CosmosProofValue{
		Kp:    kp.String(),
		Value: res.Response.GetValue(),
	})

	txParam := new(ccmcommon.MakeTxParam)
	_ = txParam.Deserialization(polycommon.NewZeroCopySource(res.Response.GetValue()))

	// Sign
	relayerAddress := common.Hex2Bytes(polySigner.Address.ToHexString())
	tx, err := SignImportOuterTransfer(*polySdk, chainID, proofValue, uint32(headerHeight), proof, relayerAddress, headerOrCrossChainMsg, polySigner)
	if err != nil {
		panic("Signing error: " + err.Error())
	}
	sink := polycommon.NewZeroCopySink(nil)
	err = tx.Serialization(sink)
	if err != nil {
		panic("Serialization error: " + err.Error())
	}
	return sink.Bytes()
}

func EnsureArgsLength(args []string, expected int) {
	lenArgs := len(args)
	if lenArgs != expected {
		panic("Expecting " + strconv.Itoa(expected) +
			" args, only received " + strconv.Itoa(lenArgs) + " args")
	}
}
