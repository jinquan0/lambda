package main

import (
        "log"
        "context"
        "net/http"
        "github.com/aws/aws-lambda-go/events"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/awslabs/aws-lambda-go-api-proxy/gin"
        "github.com/gin-gonic/gin"
        // "github.com/tidwall/gjson"
        "fmt"
)

type MyReqPara struct {
    Idx    int `json:"idx"`
    Content  string `json:content"`
}

/*
curl -XPOST -H "Content-Type: application/json" \
-d '{"idx": 1, "content": "test string... ..."}' \
https://w1xad31w0m.execute-api.cn-northwest-1.amazonaws.com.cn/v1/api/0
*/
func gin_handler_api_0(c *gin.Context) {
	var reply string
	var para MyReqPara
	if err := c.BindJSON(&para); err != nil {
		reply = "Gin BindJSON failed."
	} else {
		reply = fmt.Sprintf("API 0, request parameters ->\n\tidx: %d\n\tcontent: %s", para.Idx, para.Content)
	}
	c.String(http.StatusOK, reply+"\n")
}

func gin_handler_healthy(c *gin.Context) {
	/*
	c.JSON(200, gin.H{
		"message": "hello, lambda function is ok",
	})
	*/
	c.String(http.StatusOK, "hello, lambda function is ok\n")
}

var ginLambda *ginadapter.GinLambda

func init() {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Gin cold start")
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi, its me, the root path[/]\n")
	})

	r.GET("/healthy", gin_handler_healthy)
	r.POST("/api/0", gin_handler_api_0)

	ginLambda = ginadapter.New(r)

}

func APIGW_Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(APIGW_Handler)
}