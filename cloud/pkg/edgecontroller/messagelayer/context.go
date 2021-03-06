package messagelayer

import (
	"strings"

	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/common/constants"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
)

// MessageLayer define all functions that message layer must implement
type MessageLayer interface {
	Send(message model.Message) error
	Receive() (model.Message, error)
	Response(message model.Message) error
}

// ContextMessageLayer build on context
type ContextMessageLayer struct {
	SendModuleName       string
	SendRouterModuleName string
	ReceiveModuleName    string
	ResponseModuleName   string
}

// Send message
func (cml *ContextMessageLayer) Send(message model.Message) error {
	module := cml.SendModuleName
	// if message is rule/ruleendpoint type, send to router module.
	if isRouterMsg(message) {
		module = cml.SendRouterModuleName
	}
	beehiveContext.Send(module, message)
	return nil
}

func isRouterMsg(message model.Message) bool {
	resourceArray := strings.Split(message.GetResource(), constants.ResourceSep)
	return len(resourceArray) == 2 && (resourceArray[0] == model.ResourceTypeRule || resourceArray[0] == model.ResourceTypeRuleEndpoint)
}

// Receive message
func (cml *ContextMessageLayer) Receive() (model.Message, error) {
	return beehiveContext.Receive(cml.ReceiveModuleName)
}

// Response message
func (cml *ContextMessageLayer) Response(message model.Message) error {
	beehiveContext.Send(cml.ResponseModuleName, message)
	return nil
}

// NewContextMessageLayer create a ContextMessageLayer
func NewContextMessageLayer(config *v1alpha1.ControllerContext) MessageLayer {
	return &ContextMessageLayer{
		SendModuleName:       string(config.SendModule),
		SendRouterModuleName: string(config.SendRouterModule),
		ReceiveModuleName:    string(config.ReceiveModule),
		ResponseModuleName:   string(config.ResponseModule),
	}
}
