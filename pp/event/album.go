package event

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stratosnet/sds/framework/core"
	"github.com/stratosnet/sds/msg/header"
	"github.com/stratosnet/sds/msg/protos"
	"github.com/stratosnet/sds/pp/client"
	"github.com/stratosnet/sds/pp/peers"
	"github.com/stratosnet/sds/pp/setting"
	"github.com/stratosnet/sds/pp/types"
	"github.com/stratosnet/sds/utils"
)

// CreateAlbum
func CreateAlbum(albumName, albumBlurb, albumCoverHash, albumType, reqID string, files []*protos.FileInfo, isPrivate bool, w http.ResponseWriter) {
	aType := protos.AlbumType_ALL
	if albumType == "1" {
		aType = protos.AlbumType_VIDEO
	} else if albumType == "2" {
		aType = protos.AlbumType_MUSIC
	} else if albumType == "3" {
		aType = protos.AlbumType_OTHER
	}
	if setting.CheckLogin() {
		peers.SendMessage(client.PPConn, types.ReqCreateAlbumData(albumName, albumBlurb, albumCoverHash, reqID, aType, files, isPrivate), header.ReqCreateAlbum)
		storeResponseWriter(reqID, w)
	} else {
		notLogin(w)
	}

}

// ReqCreateAlbum ReqCreateAlbum
func ReqCreateAlbum(ctx context.Context, conn core.WriteCloser) {
	peers.TransferSendMessageToSPServer(core.MessageFromContext(ctx))
}

// RspCreateAlbum  RspCreateAlbum
func RspCreateAlbum(ctx context.Context, conn core.WriteCloser) {
	var target protos.RspCreateAlbum
	if types.UnmarshalData(ctx, &target) {
		if target.P2PAddress == setting.P2PAddress {
			if target.Result.State == protos.ResultState_RES_SUCCESS {
				fmt.Println("action  successfully", target.Result.Msg, "AlbumId", target.AlbumId)
			} else {
				fmt.Println("action  failed", target.Result.Msg)
			}
			putData(target.ReqId, HTTPCreateAlbum, &target)
		} else {
			peers.TransferSendMessageToClient(target.P2PAddress, core.MessageFromContext(ctx))
		}
	}
}

// FindMyAlbum FindMyAlbum
func FindMyAlbum(reqID string, page, number uint64, albumType, keyword string, w http.ResponseWriter) {
	aType := protos.AlbumType_ALL
	if albumType == "1" {
		aType = protos.AlbumType_VIDEO
	} else if albumType == "2" {
		aType = protos.AlbumType_MUSIC
	} else if albumType == "3" {
		aType = protos.AlbumType_OTHER
	}
	if setting.CheckLogin() {
		peers.SendMessage(client.PPConn, types.ReqFindMyAlbumData(aType, reqID, page, number, keyword), header.ReqFindMyAlbum)
		storeResponseWriter(reqID, w)
	} else {
		notLogin(w)
	}
}

// ReqFindMyAlbum
func ReqFindMyAlbum(ctx context.Context, conn core.WriteCloser) {
	peers.TransferSendMessageToSPServer(core.MessageFromContext(ctx))
}

// RspFindMyAlbum
func RspFindMyAlbum(ctx context.Context, conn core.WriteCloser) {
	var target protos.RspFindMyAlbum
	if types.UnmarshalData(ctx, &target) {
		if target.P2PAddress == setting.P2PAddress {
			if target.Result.State == protos.ResultState_RES_SUCCESS {
				for _, info := range target.AlbumInfo {
					fmt.Println("AlbumId:", info.AlbumId)
					fmt.Println("AlbumName:", info.AlbumName)
					fmt.Println("AlbumBlurb:", info.AlbumBlurb)
					fmt.Println("AlbumVisit:", info.AlbumVisit)
					fmt.Println("AlbumTime:", info.AlbumTime)
				}
			} else {
				fmt.Println("action  failed", target.Result.Msg)
			}
			putData(target.ReqId, HTTPFindMyAlbum, &target)
		} else {
			peers.TransferSendMessageToClient(target.P2PAddress, core.MessageFromContext(ctx))
		}
	}
}

// EditAlbum
func EditAlbum(albumID, albumCoverHash, albumName, albumBlurb, reqID string, changeFiles []*protos.FileInfo, isPrivate bool, w http.ResponseWriter) {

	if setting.CheckLogin() {
		peers.SendMessage(client.PPConn, types.ReqEditAlbumData(albumID, albumCoverHash, albumName, albumBlurb, reqID, changeFiles, isPrivate), header.ReqEditAlbum)
		storeResponseWriter(reqID, w)
	} else {
		notLogin(w)
	}
}

// ReqEditAlbum
func ReqEditAlbum(ctx context.Context, conn core.WriteCloser) {
	peers.TransferSendMessageToSPServer(core.MessageFromContext(ctx))
}

// RspEditAlbum
func RspEditAlbum(ctx context.Context, conn core.WriteCloser) {
	var target protos.RspEditAlbum
	if types.UnmarshalData(ctx, &target) {
		if target.P2PAddress == setting.P2PAddress {
			if target.Result.State == protos.ResultState_RES_SUCCESS {
				fmt.Println("action  successfully", target.Result.Msg)
			} else {
				fmt.Println("action  failed", target.Result.Msg)
			}
			putData(target.ReqId, HTTPEditAlbum, &target)
		} else {
			peers.TransferSendMessageToClient(target.P2PAddress, core.MessageFromContext(ctx))
		}
	}
}

// AlbumContent
func AlbumContent(albumID, reqID string, w http.ResponseWriter) {
	if setting.CheckLogin() {
		peers.SendMessage(client.PPConn, types.ReqAlbumContentData(albumID, reqID), header.ReqAlbumContent)
		storeResponseWriter(reqID, w)
	} else {
		notLogin(w)
	}
}

// ReqAlbumContent
func ReqAlbumContent(ctx context.Context, conn core.WriteCloser) {
	peers.TransferSendMessageToSPServer(core.MessageFromContext(ctx))
}

// RspAlbumContent
func RspAlbumContent(ctx context.Context, conn core.WriteCloser) {
	var target protos.RspAlbumContent
	if types.UnmarshalData(ctx, &target) {
		if target.P2PAddress == setting.P2PAddress {
			if target.Result.State == protos.ResultState_RES_SUCCESS {
				fmt.Println("AlbumInfo", target.AlbumInfo)
				fmt.Println("FileInfo", target.FileInfo)
			} else {
				fmt.Println("action  failed", target.Result.Msg)
			}
			putData(target.ReqId, HTTPAlbumContent, &target)
		} else {
			peers.TransferSendMessageToClient(target.P2PAddress, core.MessageFromContext(ctx))
		}
	}
}

// SearchAlbum
func SearchAlbum(keyword, albumType, sortType, reqID string, page, number uint64, w http.ResponseWriter) {
	aType := protos.AlbumType_ALL
	if albumType == "1" {
		aType = protos.AlbumType_VIDEO
	} else if albumType == "2" {
		aType = protos.AlbumType_MUSIC
	} else if albumType == "3" {
		aType = protos.AlbumType_OTHER
	}
	sType := protos.AlbumSortType_LATEST
	if sortType == "1" {
		sType = protos.AlbumSortType_VISITS
	}

	if setting.CheckLogin() {
		peers.SendMessage(client.PPConn, types.ReqSearchAlbumData(keyword, reqID, aType, sType, page, number), header.ReqSearchAlbum)
		storeResponseWriter(reqID, w)
	} else {
		notLogin(w)
	}
}

// ReqSearchAlbum
func ReqSearchAlbum(ctx context.Context, conn core.WriteCloser) {
	peers.TransferSendMessageToSPServer(core.MessageFromContext(ctx))
}

// RspSearchAlbum
func RspSearchAlbum(ctx context.Context, conn core.WriteCloser) {
	var target protos.RspSearchAlbum
	if types.UnmarshalData(ctx, &target) {
		if target.P2PAddress == setting.P2PAddress {
			if target.Result.State == protos.ResultState_RES_SUCCESS {
				fmt.Println("AlbumInfo", target.AlbumInfo)
			} else {
				fmt.Println("action  failed", target.Result.Msg)
			}
			putData(target.ReqId, HTTPSearchAlbum, &target)
		} else {
			peers.TransferSendMessageToClient(target.P2PAddress, core.MessageFromContext(ctx))
		}
	}
}

// CollectionAlbum
func CollectionAlbum(albumID, reqID string, isCollection bool, w http.ResponseWriter) {
	if setting.CheckLogin() {

		peers.SendMessage(client.PPConn, types.ReqCollectionAlbumData(albumID, reqID, isCollection), header.ReqCollectionAlbum)
		storeResponseWriter(reqID, w)
	} else {
		notLogin(w)
	}
}

// ReqCollectionAlbum
func ReqCollectionAlbum(ctx context.Context, conn core.WriteCloser) {
	peers.TransferSendMessageToSPServer(core.MessageFromContext(ctx))
}

// RspCollectionAlbum
func RspCollectionAlbum(ctx context.Context, conn core.WriteCloser) {
	var target protos.RspCollectionAlbum
	if types.UnmarshalData(ctx, &target) {
		if target.P2PAddress == setting.P2PAddress {
			if target.Result.State == protos.ResultState_RES_SUCCESS {
				fmt.Println("action  successfully", target.Result.Msg)
			} else {
				fmt.Println("action  failed", target.Result.Msg)
			}
			putData(target.ReqId, HTTPCollectionAlbum, &target)
		} else {
			peers.TransferSendMessageToClient(target.P2PAddress, core.MessageFromContext(ctx))
		}
	}
}

// AbstractAlbum
func AbstractAlbum(reqID string, w http.ResponseWriter) {
	if setting.CheckLogin() {

		peers.SendMessage(client.PPConn, types.ReqAbstractAlbumData(reqID), header.ReqAbstractAlbum)
		storeResponseWriter(reqID, w)
	} else {
		notLogin(w)
	}
}

// ReqAbstractAlbum
func ReqAbstractAlbum(ctx context.Context, conn core.WriteCloser) {
	peers.TransferSendMessageToSPServer(core.MessageFromContext(ctx))
}

// RspAbstractAlbum
func RspAbstractAlbum(ctx context.Context, conn core.WriteCloser) {
	var target protos.RspAbstractAlbum
	if types.UnmarshalData(ctx, &target) {
		if target.P2PAddress == setting.P2PAddress {
			if target.Result.State == protos.ResultState_RES_SUCCESS {
				fmt.Println("action  successfully", target.Result.Msg)
				utils.DebugLog("target", target)
			} else {
				fmt.Println("action  failed", target.Result.Msg)
			}
			putData(target.ReqId, HTTPAbstractAlbum, &target)
		} else {
			peers.TransferSendMessageToClient(target.P2PAddress, core.MessageFromContext(ctx))
		}
	}
}

// MyCollectionAlbum
func MyCollectionAlbum(albumType, reqID string, page, number uint64, keyword string, w http.ResponseWriter) {
	aType := protos.AlbumType_ALL
	if albumType == "1" {
		aType = protos.AlbumType_VIDEO
	} else if albumType == "2" {
		aType = protos.AlbumType_MUSIC
	} else if albumType == "3" {
		aType = protos.AlbumType_OTHER
	}
	if setting.CheckLogin() {
		peers.SendMessage(client.PPConn, types.ReqMyCollectionAlbumData(aType, reqID, page, number, keyword), header.ReqMyCollectionAlbum)
		storeResponseWriter(reqID, w)
	} else {
		notLogin(w)
	}
}

// ReqMyCollectionAlbum
func ReqMyCollectionAlbum(ctx context.Context, conn core.WriteCloser) {
	peers.TransferSendMessageToSPServer(core.MessageFromContext(ctx))
}

// RspMyCollectionAlbum
func RspMyCollectionAlbum(ctx context.Context, conn core.WriteCloser) {
	var target protos.RspMyCollectionAlbum
	if types.UnmarshalData(ctx, &target) {
		if target.P2PAddress == setting.P2PAddress {
			if target.Result.State == protos.ResultState_RES_SUCCESS {
				fmt.Println("action  successfully", target.Result.Msg)
				utils.DebugLog("target", target)
			} else {
				fmt.Println("action  failed", target.Result.Msg)
			}
			putData(target.ReqId, HTTPMyCollectionAlbum, &target)
		} else {
			peers.TransferSendMessageToClient(target.P2PAddress, core.MessageFromContext(ctx))
		}
	}
}

// DeleteAlbum
func DeleteAlbum(albumID, reqID string, w http.ResponseWriter) {
	if setting.CheckLogin() {
		peers.SendMessage(client.PPConn, types.ReqDeleteAlbumData(albumID, reqID), header.ReqDeleteAlbum)
		storeResponseWriter(reqID, w)
	} else {
		notLogin(w)
	}
}

// ReqDeleteAlbum
func ReqDeleteAlbum(ctx context.Context, conn core.WriteCloser) {
	peers.TransferSendMessageToSPServer(core.MessageFromContext(ctx))
}

// RspDeleteAlbum
func RspDeleteAlbum(ctx context.Context, conn core.WriteCloser) {
	var target protos.RspDeleteAlbum
	if types.UnmarshalData(ctx, &target) {
		if target.P2PAddress == setting.P2PAddress {
			if target.Result.State == protos.ResultState_RES_SUCCESS {
				fmt.Println("action  successfully", target.Result.Msg)
				utils.DebugLog("target", target)
			} else {
				fmt.Println("action  failed", target.Result.Msg)
			}
			putData(target.ReqId, HTTPDeleteAlbum, &target)
		} else {
			peers.TransferSendMessageToClient(target.P2PAddress, core.MessageFromContext(ctx))
		}
	}
}
