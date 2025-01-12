package stratoschain

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stratosnet/sds/framework/core"
	"github.com/stratosnet/sds/utils/crypto/ed25519"
	utiltypes "github.com/stratosnet/sds/utils/types"
	pottypes "github.com/stratosnet/stratos-chain/x/pot/types"
	registertypes "github.com/stratosnet/stratos-chain/x/register/types"
	sdstypes "github.com/stratosnet/stratos-chain/x/sds/types"
	"github.com/tendermint/tendermint/libs/bech32"
)

// Stratos-chain 'pot' module
func BuildVolumeReportMsg(traffic []*core.Traffic, reporterAddress, reporterOwnerAddress []byte, epoch uint64, reportReference string) (sdktypes.Msg, error) {
	aggregatedVolume := make(map[string]uint64)
	for _, trafficReccord := range traffic {
		aggregatedVolume[trafficReccord.P2PAddress] += trafficReccord.Volume
	}

	var nodesVolume []pottypes.SingleNodeVolume
	for p2pAddressString, volume := range aggregatedVolume {
		_, p2pAddressBytes, err := bech32.DecodeAndConvert(p2pAddressString)
		if err != nil {
			return nil, err
		}
		p2pAddress := sdktypes.AccAddress(p2pAddressBytes[:])
		nodesVolume = append(nodesVolume, pottypes.SingleNodeVolume{
			NodeAddress: p2pAddress,
			Volume:      sdktypes.NewIntFromUint64(volume),
		})
	}

	return pottypes.NewMsgVolumeReport(nodesVolume, reporterAddress, sdktypes.NewIntFromUint64(epoch), reportReference, reporterOwnerAddress), nil
}

// Stratos-chain 'register' module
func BuildCreateResourceNodeMsg(networkID, token, moniker, nodeType string, pubKey []byte, stakeAmount int64, ownerAddress utiltypes.Address) sdktypes.Msg {
	if nodeType == "" {
		nodeType = registertypes.STORAGE.Type()
	}
	return registertypes.NewMsgCreateResourceNode(
		networkID,
		ed25519.PubKeyBytesToPubKey(pubKey),
		sdktypes.NewInt64Coin(token, stakeAmount),
		ownerAddress[:],
		registertypes.Description{
			Moniker: moniker,
		},
		nodeType,
	)
}

func BuildCreateIndexingNodeMsg(networkID, token, moniker string, pubKey []byte, stakeAmount int64, ownerAddress utiltypes.Address) sdktypes.Msg {
	return registertypes.NewMsgCreateIndexingNode(
		networkID,
		ed25519.PubKeyBytesToPubKey(pubKey),
		sdktypes.NewInt64Coin(token, stakeAmount),
		ownerAddress[:],
		registertypes.Description{
			Moniker: moniker,
		},
	)
}

func BuildRemoveResourceNodeMsg(nodeAddress, ownerAddress utiltypes.Address) sdktypes.Msg {
	return registertypes.NewMsgRemoveResourceNode(
		nodeAddress[:],
		ownerAddress[:],
	)
}

func BuildRemoveIndexingNodeMsg(nodeAddress, ownerAddress utiltypes.Address) sdktypes.Msg {
	return registertypes.NewMsgRemoveIndexingNode(
		nodeAddress[:],
		ownerAddress[:],
	)
}

func BuildIndexingNodeRegistrationVoteMsg(candidateNetworkAddress, candidateOwnerAddress, voterNetworkAddress, voterOwnerAddress utiltypes.Address, voteOpinion bool) sdktypes.Msg {
	return registertypes.NewMsgIndexingNodeRegistrationVote(
		candidateNetworkAddress[:],
		candidateOwnerAddress[:],
		registertypes.VoteOpinionFromBool(voteOpinion),
		voterNetworkAddress[:],
		voterOwnerAddress[:],
	)
}

// Stratos-chain 'sds' module
func BuildFileUploadMsg(fileHash, reporterAddress, uploaderAddress []byte) sdktypes.Msg {
	return sdstypes.NewMsgUpload(
		fileHash,
		reporterAddress,
		uploaderAddress,
	)
}

func BuildPrepayMsg(token string, amount int64, senderAddress []byte) sdktypes.Msg {
	return sdstypes.NewMsgPrepay(
		senderAddress,
		sdktypes.NewCoins(sdktypes.NewInt64Coin(token, amount)),
	)
}
