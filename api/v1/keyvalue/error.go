package keyvalue

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	// "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorNoSuchKey struct {
	Key string
}

func (e ErrorNoSuchKey) RetErr() error {
	// st := status.New(codes.NotFound, fmt.Sprintf("No such key: %s", e.Key))
	st := status.New(
		404,
		fmt.Sprintf("No such key: %s", e.Key),
	)

	msg := fmt.Sprintf(
		"No such key: %s",
		e.Key,
	)
	d := &errdetails.LocalizedMessage{
		Locale:  "en-US",
		Message: msg,
	}

	std, err := st.WithDetails(d)
	if err != nil {
		return err
	}
	return std.Err()
}
