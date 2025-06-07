// Code generated : DO NOT EDIT.
// Copyright (c) 2022 Jean-Francois SMIGIELSKI
// Distributed under the MIT License

package media

import (
	"context"
	"fmt"
	"go-onvif/media"
	"go-onvif/onvif"
	"go-onvif/sdk"

	"github.com/juju/errors"
)

// Call_GetSnapshotUri forwards the call to dev.CallMethod() then parses the payload of the reply as a GetSnapshotUriResponse.
func Call_GetSnapshotUri(ctx context.Context, dev *onvif.Device, request media.GetSnapshotUri) (media.GetSnapshotUriResponse, error) {
	type Envelope struct {
		Header struct{}
		Body   struct {
			GetSnapshotUriResponse media.GetSnapshotUriResponse
		}
	}
	var reply Envelope
	httpReply, err := dev.CallMethod(request);
	fmt.Println(httpReply, err)
	if err != nil {
		fmt.Println(err)
		return reply.Body.GetSnapshotUriResponse, errors.Annotate(err, "call")
	} else {
		fmt.Println(httpReply)
		err = sdk.ReadAndParse(ctx, httpReply, &reply, "GetSnapshotUri")
		return reply.Body.GetSnapshotUriResponse, errors.Annotate(err, "reply")
	}
}
