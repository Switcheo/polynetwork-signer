/*
* Copyright (C) 2020 The poly network Authors
* This file is part of The poly network library.
*
* The poly network is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The poly network is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
* You should have received a copy of the GNU Lesser General Public License
* along with The poly network . If not, see <http://www.gnu.org/licenses/>.
 */
package app

import (
	"github.com/Switcheo/polynetwork-signer/config"
	sdk "github.com/polynetwork/poly-go-sdk"
	"github.com/polynetwork/poly/common"
)

func SetUpPoly(conf *config.Config) (*sdk.PolySdk, error) {
	rpcAddr := conf.RpcEndPoints.PolyNetwork
	polySdk := sdk.NewPolySdk()
	polySdk.NewRpcClient().SetAddress(rpcAddr)
	hdr, err := polySdk.GetHeaderByHeight(0)
	if err != nil {
		return nil, err
	}
	polySdk.SetChainId(hdr.ChainID)
	return polySdk, nil
}

func SetupPolySigner(conf *config.Config, polySdk *sdk.PolySdk) (*sdk.Account, error) {
	var wallet *sdk.Wallet
	var err error
	walletFile := conf.BroadCaster.PolyNetwork.FileName
	walletPassword := conf.BroadCaster.PolyNetwork.Password
	if !common.FileExisted(walletFile) {
		wallet, err = polySdk.CreateWallet(walletFile)
		if err != nil {
			return nil, err
		}
	} else {
		wallet, err = polySdk.OpenWallet(walletFile)
		if err != nil {
			return nil, err
		}
	}
	signer, err := wallet.GetDefaultAccount([]byte(walletPassword))
	if err != nil || signer == nil {
		signer, err = wallet.NewDefaultSettingAccount([]byte(walletPassword))
		if err != nil {
			return nil, err
		}

		err = wallet.Save()
		if err != nil {
			return nil, err
		}
	}
	return signer, nil
}
