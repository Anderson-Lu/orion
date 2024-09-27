package xcontext

import (
	"google.golang.org/grpc/metadata"
)

const (
	KEY_X_UIT_REQUEST_ID = "x-uit-request-id"
	KEY_X_UIT_CLIENT_IP  = "x-uit-client-ip"
	KEY_X_UIT_UID        = "x-uit-uid"
)

type TracerHeader struct {
	RequestId string
	ClientIP  string
	Uid       string
}

func filterTraceMD(md metadata.MD) metadata.MD {
	return nil
}
