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
	sdk "github.com/polynetwork/poly-go-sdk"
	polytypes "github.com/polynetwork/poly/core/types"
)

// sourceChainID = eth chain id
// txData = rawdata
// height = any height that is > tx height + confirmation height
// proof = gotten from the eccd contract using proofKey (where key := crosstx.txIndex; keyBytes, err := eth.MappingKeyAt(key, "01"); proofKey := hexutil.Encode(keyBytes);
// relayerAddress = ethcommon.Hex2Bytes(this.polySigner.Address.ToHexString())
// headerOrCrossChainMsg = not sure what this is yet but commitProof for eth is empty bytes ( []byte{} ), neo has something
// signer = ethcommon.Hex2Bytes(this.polySigner.Address.ToHexString()),
func SignImportOuterTransfer(polySdk sdk.PolySdk, sourceChainId uint64, txData []byte, height uint32, proof []byte,
	relayerAddress []byte, HeaderOrCrossChainMsg []byte, signer *sdk.Account) (*polytypes.Transaction, error) {
	tx, err := polySdk.Native.Ccm.NewImportOuterTransferTransaction(
		sourceChainId, txData, height, proof, relayerAddress, HeaderOrCrossChainMsg)
	if err != nil {
		return nil, err
	}
	err = polySdk.SignToTransaction(tx, signer)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

