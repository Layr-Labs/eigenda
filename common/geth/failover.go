package geth

import (
	"net/url"
	"sync"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

type FailoverController struct {
	mu             *sync.RWMutex
	numberRpcFault uint64
	UrlDomains     []string

	Logger logging.Logger
}

func NewFailoverController(logger logging.Logger, rpcUrls []string) (*FailoverController, error) {
	urlDomains := make([]string, len(rpcUrls))
	for i := 0; i < len(urlDomains); i++ {
		url, err := url.Parse(rpcUrls[i])
		if err != nil {
			return nil, err
		}
		urlDomains[i] = url.Hostname()
	}
	return &FailoverController{
		Logger:     logger.With("component", "FailoverController"),
		mu:         &sync.RWMutex{},
		UrlDomains: urlDomains,
	}, nil
}

// ProcessError attributes the error and updates total number of fault for RPC
// It returns if RPC should immediately give up
func (f *FailoverController) ProcessError(err error, rpcIndex int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err == nil {
		return false
	}

	urlDomain := ""
	if rpcIndex >= len(f.UrlDomains) || rpcIndex < 0 {
		f.Logger.Error("[FailoverController]", "err", "rpc index is outside of known url")
	} else {
		urlDomain = f.UrlDomains[rpcIndex]
	}

	nextEndpoint, action := f.handleError(err, urlDomain)

	if nextEndpoint == NewRPC {
		f.numberRpcFault += 1
	}

	return action == Return
}

func (f *FailoverController) GetTotalNumberRpcFault() uint64 {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.numberRpcFault
}
