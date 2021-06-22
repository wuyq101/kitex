package loadbalance

import (
	"context"

	"github.com/cloudwego/kitex/pkg/discovery"
)

// Picker picks an instance for next RPC call.
type Picker interface {
	Next(ctx context.Context, request interface{}) discovery.Instance
}

// Loadbalancer generates pickers for the given service discovery result.
type Loadbalancer interface {
	GetPicker(discovery.Result) Picker
	Name() string // unique key
}

// Rebalancer is a kind of Loadbalancer that performs rebalancing when the result of service discovery changes.
type Rebalancer interface {
	Rebalance(discovery.Change)
	Delete(discovery.Change)
}

// SynthesizedPicker synthesizes a picker using a next function.
type SynthesizedPicker struct {
	NextFunc func(ctx context.Context, request interface{}) discovery.Instance
}

// Next implements the Picker interface.
func (sp *SynthesizedPicker) Next(ctx context.Context, request interface{}) discovery.Instance {
	return sp.NextFunc(ctx, request)
}

// SynthesizedLoadbalancer synthesizes a loadbalancer using a GetPickerFunc function.
type SynthesizedLoadbalancer struct {
	GetPickerFunc func(discovery.Result) Picker
	NameFunc      func() string
}

// GetPicker implements the Loadbalancer interface.
func (slb *SynthesizedLoadbalancer) GetPicker(entry discovery.Result) Picker {
	if slb.GetPickerFunc == nil {
		return nil
	}
	return slb.GetPickerFunc(entry)
}

// Name implements the Loadbalancer interface.
func (slb *SynthesizedLoadbalancer) Name() string {
	if slb.NameFunc == nil {
		return ""
	}
	return slb.NameFunc()
}
