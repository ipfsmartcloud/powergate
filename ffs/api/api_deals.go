package api

import (
	"fmt"

	"github.com/textileio/powergate/deals"
)

// ListStorageDealRecords lists storage deals for this FFS instance according to the provided options.
func (i *API) ListStorageDealRecords(opts ...deals.ListDealRecordsOption) ([]deals.StorageDealRecord, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	c := deals.ListDealRecordsConfig{}
	for _, opt := range opts {
		opt(&c)
	}
	finalAddrs, err := i.finalAddresses(c.FromAddrs)
	if err != nil {
		return nil, fmt.Errorf("getting final addrs: %v", err)
	}
	recs, err := i.drm.ListStorageDealRecords(
		deals.WithFromAddrs(finalAddrs...),
		deals.WithAscending(c.Ascending),
		deals.WithDataCids(c.DataCids...),
		deals.WithIncludeFinal(c.IncludeFinal),
		deals.WithIncludePending(c.IncludePending),
	)
	if err != nil {
		return nil, fmt.Errorf("calling ListStorageDealRecords: %v", err)
	}
	return recs, nil
}

// ListRetrievalDealRecords returns a list of retrieval deals for this FFS instance according to the provided options.
func (i *API) ListRetrievalDealRecords(opts ...deals.ListDealRecordsOption) ([]deals.RetrievalDealRecord, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	c := deals.ListDealRecordsConfig{}
	for _, opt := range opts {
		opt(&c)
	}
	finalAddrs, err := i.finalAddresses(c.FromAddrs)
	if err != nil {
		return nil, fmt.Errorf("getting final addrs: %v", err)
	}
	recs, err := i.drm.ListRetrievalDealRecords(
		deals.WithFromAddrs(finalAddrs...),
		deals.WithAscending(c.Ascending),
		deals.WithDataCids(c.DataCids...),
	)
	if err != nil {
		return nil, fmt.Errorf("calling dm.ListRetrievalDealRecords: %v", err)
	}
	return recs, nil
}

func (i *API) finalAddresses(fromAddrs []string) ([]string, error) {
	instanceAddrs := make([]string, 0, len(i.cfg.Addrs))
	instanceAddrsFilter := make(map[string]struct{})
	for _, addrInfo := range i.cfg.Addrs {
		instanceAddrs = append(instanceAddrs, addrInfo.Addr)
		instanceAddrsFilter[addrInfo.Addr] = struct{}{}
	}

	var finalAddrs []string
	if len(fromAddrs) > 0 {
		for _, addr := range fromAddrs {
			if _, ok := instanceAddrsFilter[addr]; !ok {
				return nil, fmt.Errorf("address %s is not managed by this ffs instance", addr)
			}
			finalAddrs = append(finalAddrs, addr)
		}
	} else {
		finalAddrs = instanceAddrs
	}
	return finalAddrs, nil
}
