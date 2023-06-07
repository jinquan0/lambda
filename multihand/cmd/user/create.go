package user

import (
	"context"
	"fmt"
	"time"
)

type CreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateResponse struct {
	OK    bool   `json:"ok"`
	ID    int64  `json:"id"`
	ReqID string `json:"req_id"`
}

func (u User) Create(_ context.Context, reqID string, req CreateRequest) (CreateResponse, error) {
	if req.FirstName == "" {
		return CreateResponse{OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("the first_name is missing")
	}
	if req.LastName == "" {
		return CreateResponse{OK: false, ID: 0, ReqID: reqID}, fmt.Errorf("the last_name is missing")
	}

	return CreateResponse{OK: true, ID: time.Now().UnixNano(), ReqID: reqID}, nil
}