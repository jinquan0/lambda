package user

import (
	"context"
	"fmt"
)

type DeleteRequest struct {
	ID string `json:"id"`
}

type DeleteResponse struct {
	OK    bool   `json:"ok"`
	ReqID string `json:"req_id"`
}

func (u User) Delete(_ context.Context, reqID string, req DeleteRequest) (DeleteResponse, error) {
	if req.ID == "" {
		return DeleteResponse{OK: false, ReqID: reqID}, fmt.Errorf("the id key is missing")
	}

	return DeleteResponse{OK: true, ReqID: reqID}, nil
}