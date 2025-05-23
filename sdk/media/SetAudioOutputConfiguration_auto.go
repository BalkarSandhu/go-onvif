// Code generated : DO NOT EDIT.
// Copyright (c) 2022 Jean-Francois SMIGIELSKI
// Distributed under the MIT License

package media

import (
	"context"
	"github.com/juju/errors"
	"github.com/BalkarSandhu/go-onvif/onvif"
	"github.com/BalkarSandhu/go-onvif/sdk"
	"github.com/BalkarSandhu/go-onvif/media"
)

// Call_SetAudioOutputConfiguration forwards the call to dev.CallMethod() then parses the payload of the reply as a SetAudioOutputConfigurationResponse.
func Call_SetAudioOutputConfiguration(ctx context.Context, dev *onvif.Device, request media.SetAudioOutputConfiguration) (media.SetAudioOutputConfigurationResponse, error) {
	type Envelope struct {
		Header struct{}
		Body   struct {
			SetAudioOutputConfigurationResponse media.SetAudioOutputConfigurationResponse
		}
	}
	var reply Envelope
	if httpReply, err := dev.CallMethod(request); err != nil {
		return reply.Body.SetAudioOutputConfigurationResponse, errors.Annotate(err, "call")
	} else {
		err = sdk.ReadAndParse(ctx, httpReply, &reply, "SetAudioOutputConfiguration")
		return reply.Body.SetAudioOutputConfigurationResponse, errors.Annotate(err, "reply")
	}
}
