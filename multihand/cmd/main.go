package main

import (
	"github.com/aws/aws-lambda-go/lambda"
    "context"
    "encoding/json"
    "fmt"

    "github.com/aws/aws-lambda-go/lambdacontext"
    _ "context"
    _ "fmt"
    "time"


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


func main() {
	lambda.Start(User{}.Handle)
}