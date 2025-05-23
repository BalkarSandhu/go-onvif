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

// Call_GetCACertificates forwards the call to dev.CallMethod() then parses the payload of the reply as a GetCACertificatesResponse.
func Call_GetCACertificates(ctx context.Context, dev *onvif.Device, request device.GetCACertificates) (device.GetCACertificatesResponse, error) {
	type Envelope struct {
		Header struct{}
		Body   struct {
			GetCACertificatesResponse device.GetCACertificatesResponse
		}
	}
	var reply Envelope
	if httpReply, err := dev.CallMethod(request); err != nil {
		return reply.Body.GetCACertificatesResponse, errors.Annotate(err, "call")
	} else {
		err = sdk.ReadAndParse(ctx, httpReply, &reply, "GetCACertificates")
		return reply.Body.GetCACertificatesResponse, errors.Annotate(err, "reply")
	}
}
