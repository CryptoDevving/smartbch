package api

import (
	"encoding/json"
	"runtime"
	"sync/atomic"
	"time"

	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/tendermint/tendermint/libs/log"

	stakingtypes "github.com/smartbch/smartbch/staking/types"
)

const (
	StatusUpdateInterval = 60 // seconds
)

type Stats struct {
	NumGoroutine     int    `json:"numGoroutine"`
	NumGC            uint32 `json:"numGC"`
	MemAllocMB       uint64 `json:"memAllocMB"`
	MemSysMB         uint64 `json:"memSysMB"`
	OsMemTotalMB     uint64 `json:"osMemTotalMB"`
	OsMemUsedMB      uint64 `json:"osMemUsedMB"`
	OsMemCachedMB    uint64 `json:"osMemCachedMB"`
	OsMemFreeMB      uint64 `json:"osMemFreeMB"`
	OsMemActiveMB    uint64 `json:"osMemActiveMB"`
	OsMemInactiveMB  uint64 `json:"osMemInactiveMB"`
	OsMemSwapTotalMB uint64 `json:"osMemSwapTotalMB"`
	OsMemSwapUsedMB  uint64 `json:"osMemSwapUsedMB"`
	OsMemSwapFreeMB  uint64 `json:"osMemSwapFreeMB"`
	NumEthCall       uint64 `json:"numEthCall"`
}

type DebugAPI interface {
	GetStats() Stats
	GetSeq(addr gethcmn.Address) hexutil.Uint64
	NodeInfo() json.RawMessage
	ValidatorOnlineInfos() json.RawMessage
}

type debugAPI struct {
	logger         log.Logger
	ethAPI         *ethAPI
	lastUpdateTime int64
	stats          Stats
}

func newDebugAPI(ethAPI *ethAPI, logger log.Logger) DebugAPI {
	return &debugAPI{
		logger: logger,
		ethAPI: ethAPI,
	}
}

func (api *debugAPI) GetSeq(addr gethcmn.Address) hexutil.Uint64 {
	api.logger.Debug("debug_getSeq")
	return hexutil.Uint64(api.ethAPI.backend.GetSeq(addr))
}

func (api *debugAPI) NodeInfo() json.RawMessage {
	api.logger.Debug("debug_nodeInfo")
	nodeInfo := api.ethAPI.backend.NodeInfo()
	bytes, _ := json.Marshal(nodeInfo)
	return bytes
}

func (api *debugAPI) GetStats() Stats {
	api.logger.Debug("debug_getStats")

	now := time.Now().Unix()
	lastUpdateTime := atomic.LoadInt64(&api.lastUpdateTime)
	if now > lastUpdateTime+StatusUpdateInterval {
		if atomic.CompareAndSwapInt64(&api.lastUpdateTime, lastUpdateTime, now) {
			api.updateStats()
		}
	}

	return api.stats
}

func (api *debugAPI) updateStats() {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	api.stats.NumGoroutine = runtime.NumGoroutine()
	api.stats.NumGC = memStats.NumGC
	api.stats.MemAllocMB = toMB(memStats.Alloc)
	api.stats.MemSysMB = toMB(memStats.Sys)

	osMemStats, err := memory.Get()
	if err == nil {
		api.stats.OsMemTotalMB = toMB(osMemStats.Total)
		api.stats.OsMemUsedMB = toMB(osMemStats.Used)
		api.stats.OsMemCachedMB = toMB(osMemStats.Cached)
		api.stats.OsMemFreeMB = toMB(osMemStats.Free)
		api.stats.OsMemActiveMB = toMB(osMemStats.Active)
		api.stats.OsMemInactiveMB = toMB(osMemStats.Inactive)
		api.stats.OsMemSwapTotalMB = toMB(osMemStats.SwapTotal)
		api.stats.OsMemSwapUsedMB = toMB(osMemStats.SwapUsed)
		api.stats.OsMemSwapFreeMB = toMB(osMemStats.SwapFree)
	}

	api.stats.NumEthCall = atomic.LoadUint64(&api.ethAPI.numCall)
}

func toMB(n uint64) uint64 {
	return n / 1024 / 1024
}

/* Validator Online Info */

type ValidatorOnlineInfosToMarshal struct {
	StartHeight int64                  `json:"start_height"`
	OnlineInfos []*OnlineInfoToMarshal `json:"online_infos"`
}

type OnlineInfoToMarshal struct {
	ValidatorConsensusAddress gethcmn.Address `json:"validator_consensus_address"`
	SignatureCount            int32           `json:"signature_count"`
	HeightOfLastSignature     int64           `json:"height_of_last_signature"`
}

func castValidatorOnlineInfos(infos stakingtypes.ValidatorOnlineInfos) ValidatorOnlineInfosToMarshal {
	infosToMarshal := ValidatorOnlineInfosToMarshal{
		StartHeight: infos.StartHeight,
		OnlineInfos: make([]*OnlineInfoToMarshal, len(infos.OnlineInfos)),
	}
	for i, onlineInfo := range infos.OnlineInfos {
		infosToMarshal.OnlineInfos[i] = &OnlineInfoToMarshal{
			ValidatorConsensusAddress: onlineInfo.ValidatorConsensusAddress,
			SignatureCount:            onlineInfo.SignatureCount,
			HeightOfLastSignature:     onlineInfo.HeightOfLastSignature,
		}
	}
	return infosToMarshal
}

func (api *debugAPI) ValidatorOnlineInfos() json.RawMessage {
	api.logger.Debug("debug_validatorsOnlineInfo")
	onlineInfos := api.ethAPI.backend.ValidatorOnlineInfos()
	onlineInfosToMarshal := castValidatorOnlineInfos(onlineInfos)
	bytes, _ := json.Marshal(onlineInfosToMarshal)
	return bytes
}
