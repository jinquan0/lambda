package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    lbd "github.com/aws/aws-sdk-go/service/lambda"
    "github.com/aws/aws-lambda-go/lambda"

    "encoding/json"
    "fmt"
    "os"
    _ "strconv"
)

type getItemsRequest struct {
    InstanceRegion    string
    InstanceIdList     []string
}

type getItemsResponseError struct {
    Message string `json:"message"`
}

type getItemsResponseData struct {
    Item string `json:"item"`
}

type getItemsResponseBody struct {
    Result string                 `json:"result"`
    Data   []getItemsResponseData `json:"data"`
    Error  getItemsResponseError  `json:"error"`
}

type getItemsResponseHeaders struct {
    ContentType string `json:"Content-Type"`
}

type getItemsResponse struct {
    StatusCode int                     `json:"statusCode"`
    Headers    getItemsResponseHeaders `json:"headers"`
    Body       getItemsResponseBody    `json:"body"`
}

func LambdaServiceClient() *lbd.Lambda {
    // Create Lambda service client
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    var client *lbd.Lambda
    client = lbd.New(sess, &aws.Config{Region: aws.String("cn-northwest-1")})
    return client
}

func HandleLambdaEvent() {
    client := LambdaServiceClient()

    // Get the 10 most recent items
    var request getItemsRequest
    request.InstanceRegion = "cn-northwest-1"
    request.InstanceIdList = append(request.InstanceIdList, "i-0c019bf947adb3fb6")


    payload, err := json.Marshal(request)
    if err != nil {
        fmt.Println("Error marshalling ec2start request")
        os.Exit(0)
    }

    result, err := client.Invoke(&lbd.InvokeInput{FunctionName: aws.String("ec2start"), Payload: payload})
    if err != nil {
        fmt.Println("Error calling ec2start")
        os.Exit(0)
    }

    result = result

    /*
    var resp getItemsResponse

    err = json.Unmarshal(result.Payload, &resp)
    if err != nil {
        fmt.Println("Error unmarshalling ec2start response")
        os.Exit(0)
    }

    // If the status code is NOT 200, the call failed
    if resp.StatusCode != 200 {
        fmt.Println("Error getting items, StatusCode: " + strconv.Itoa(resp.StatusCode))
        os.Exit(0)
    }

    // If the result is failure, we got an error
    if resp.Body.Result == "failure" {
        fmt.Println("Failed to get items")
        os.Exit(0)
    }

    // Print out items
    if len(resp.Body.Data) > 0 {
        for i := range resp.Body.Data {
            fmt.Println(resp.Body.Data[i].Item)
        }
    } else {
        fmt.Println("There were no items")
    }
    */
}


func main() {
        
    lambda.Start(HandleLambdaEvent)
}