// Code generated : DO NOT EDIT.
// Copyright (c) 2022 Jean-Francois SMIGIELSKI
// Distributed under the MIT License

package device

import (
	"context"
	"github.com/juju/errors"
	"github.com/BalkarSandhu/go-onvif/onvif"
	"github.com/BalkarSandhu/go-onvif/sdk"
	"github.com/BalkarSandhu/go-onvif/device"
)

// Call_SendAuxiliaryCommand forwards the call to dev.CallMethod() then parses the payload of the reply as a SendAuxiliaryCommandResponse.
func Call_SendAuxiliaryCommand(ctx context.Context, dev *onvif.Device, request device.SendAuxiliaryCommand) (device.SendAuxiliaryCommandResponse, error) {
	type Envelope struct {
		Header struct{}
		Body   struct {
			SendAuxiliaryCommandResponse device.SendAuxiliaryCommandResponse
		}
	}
	var reply Envelope
	if httpReply, err := dev.CallMethod(request); err != nil {
		return reply.Body.SendAuxiliaryCommandResponse, errors.Annotate(err, "call")
	} else {
		err = sdk.ReadAndParse(ctx, httpReply, &reply, "SendAuxiliaryCommand")
		return reply.Body.SendAuxiliaryCommandResponse, errors.Annotate(err, "reply")
	}
}
