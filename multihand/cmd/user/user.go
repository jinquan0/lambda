package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambdacontext"
)
 
type HandleRequest struct {
	Event string          `json:"event"`
	Body  json.RawMessage `json:"body"`
}

type HandleResponse struct {
	OK    bool   `json:"ok"`
	ReqID string `json:"req_id"`
}

type User struct {
	// Some dependencies
}

func (u User) Handle(ctx context.Context, req HandleRequest) (interface{}, error) {
	var reqID string
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		reqID = lc.AwsRequestID
	}

	select {
	case <-ctx.Done():
		return HandleResponse{OK: false, ReqID: reqID}, fmt.Errorf("request timeout: %w", ctx.Err())
	default:
	}

	switch req.Event {
	case "create":
		var dest CreateRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return u.Create(ctx, reqID, dest)
	case "delete":
		var dest DeleteRequest
		if err := json.Unmarshal(req.Body, &dest); err != nil {
			return nil, err
		}
		return u.Delete(ctx, reqID, dest)
	}

	return HandleResponse{OK: false, ReqID: reqID}, fmt.Errorf("%s is an unknown event", req.Event)
}