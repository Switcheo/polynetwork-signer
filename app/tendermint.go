package app

import (
	"context"
	"fmt"
	headersynctypes "github.com/Switcheo/polynetwork-cosmos/x/headersync/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
	"strings"
)

type CosmosHeader struct {
	Header  tmtypes.Header
	Commit  *tmtypes.Commit
	Valsets []*CosmosValidator
}

type CosmosValidator struct {
	Address          tmtypes.Address    `json:"address"`
	PubKey           cryptotypes.PubKey `json:"pub_key"`
	VotingPower      int64              `json:"voting_power"`
	ProposerPriority int64              `json:"proposer_priority"`
}

func GetCosmosHeader(rpcClient *rpchttp.HTTP, height int64) (*CosmosHeader, error) {
	rc, err := rpcClient.Commit(context.TODO(), &height)
	if err != nil {
		return nil, fmt.Errorf("failed to get Commit of height %d: %v", height, err)
	}
	vSet, err := getValidators(rpcClient, height)
	if err != nil {
		return nil, fmt.Errorf("failed to get Validators of height %d: %v", height, err)
	}
	return &CosmosHeader{
		Header:  *rc.Header,
		Commit:  rc.Commit,
		Valsets: vSet,
	}, nil
}

func getValidators(rpcClient *rpchttp.HTTP, height int64) ([]*CosmosValidator, error) {
	page := 1
	perPage := 100
	vSet := make([]*CosmosValidator, 0)
	for {
		res, err := rpcClient.Validators(context.TODO(), &height, &page, &perPage)
		if err != nil {
			if strings.Contains(err.Error(), "page should be within") {
				return vSet, nil
			}
			return nil, err
		}
		// In case tendermint don't give relayer the right error
		if len(res.Validators) == 0 {
			return vSet, nil
		}

		for i := range res.Validators {
			pk, err := cryptocodec.FromTmPubKeyInterface(res.Validators[i].PubKey)
			if err != nil {
				panic(err)
			}
			vSet = append(vSet, &CosmosValidator{
				Address:          res.Validators[i].Address,
				PubKey:           pk,
				VotingPower:      res.Validators[i].VotingPower,
				ProposerPriority: res.Validators[i].ProposerPriority,
			})
		}
		page++
	}
}

func NewCodecForRelayer() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	headersynctypes.RegisterCodec(cdc)
	cryptocodec.RegisterCrypto(cdc)
	return cdc
}
