package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "fmt"
    "os"
    "log"
    "context"
    "net/http"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/awslabs/aws-lambda-go-api-proxy/gin"
    "github.com/gin-gonic/gin"
)

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}

// reference link: https://pkg.go.dev/github.com/aws/aws-sdk-go/aws/session#Session
func sessionInit(profile string, region string) *session.Session {
    sess, err := session.NewSessionWithOptions(session.Options{
        // Specify profile to load for the session's config
        Profile: "dsmqa",

        // Provide SDK Config options, such as Region.
        Config: aws.Config{
            Region: aws.String("cn-northwest-1"),
        },

        // Force enable Shared Config support
        SharedConfigState: session.SharedConfigEnable,
    })
    if err != nil {
        //log.Printf("session.NewSessionWithOptions failed.")
        exitErrorf("session.NewSessionWithOptions, %v", err)
    }
    return sess
}

func sessionInit_v2(region string) *session.Session {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("cn-northwest-1")},
    )
    if err != nil {
        //log.Printf("session.NewSessionWithOptions failed.")
        exitErrorf("session.NewSessionWithOptions, %v", err)
    }
    return sess
}

// reference link: https://pkg.go.dev/github.com/aws/aws-sdk-go/service/s3#S3
func s3bucketList(svc *s3.S3) {
    result, err := svc.ListBuckets(nil)
    if err != nil {
        exitErrorf("Unable to list buckets, %v", err)
    }

    fmt.Println("Buckets:")

    for _, b := range result.Buckets {
        fmt.Printf("* %s created on %s\n",
            aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
    }
}

func s3bucketListItems(svc *s3.S3, bucket string) string {
    resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
    if err != nil {
        exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
    }

    var str string = ""
    for _, item := range resp.Contents {
        str += fmt.Sprintln("Name:         ", *item.Key)
        str += fmt.Sprintln("Last modified:", *item.LastModified)
        str += fmt.Sprintln("Size:         ", *item.Size)
        str += fmt.Sprintln("Storage class:", *item.StorageClass)
        str += fmt.Sprintln("")
    }
    return str
}

func s3bucketUploadFile(sess *session.Session, bucket string, filename string) {
    file, err := os.Open(filename)
    if err != nil {
        exitErrorf("Unable to open file %q, %v", err)
    }

    defer file.Close()

    // http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
    uploader := s3manager.NewUploader(sess)

    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key: aws.String(filename),
        Body: file,
    })
    if err != nil {
        // Print the error and exit.
        exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
    }
}

/*
func main() {

    sess := sessionInit("dsmqa", "cn-northwest-1")
    // Create S3 service client
    svc := s3.New(sess)

    s3bucketList(svc)
    s3bucketListItems(svc, "qa-is-bucket")
    s3bucketUploadFile(sess, "qa-is-bucket", "go.mod")
}
*/

type MyReqPara struct {
    Bucket  string `json:bucket"`
}

/*
curl -XGET -H "Content-Type: application/json" \
-d '{"bucket": "qa-is-bucket"}' \
  https://otbxr6m1fg.execute-api.cn-northwest-1.amazonaws.com.cn/v1/api/listitems
*/
func gin_handler_api_listitems(c *gin.Context) {
    var reply string
    var para MyReqPara
    if err := c.BindJSON(&para); err != nil {
        reply = "Gin BindJSON failed."
    } else {
        sess := sessionInit_v2("cn-northwest-1")
        svc := s3.New(sess)
        reply = s3bucketListItems(svc, para.Bucket)
    }
    c.String(http.StatusOK, reply+"\n")
}

/*
curl -XPOST -F "file=@/tmp/test.txt" \
  -H "Content-Type: multipart/form-data" \
  https://otbxr6m1fg.execute-api.cn-northwest-1.amazonaws.com.cn/v1/api/upload
*/
func gin_handler_api_upload(c *gin.Context) {
    // var para MyReqPara
    //FormFile返回所提供的表单键的第一个文件
    f, err := c.FormFile("file")
    if err != nil {
        c.String(http.StatusOK, "Gin FormFile Failed.\n")
    } else {
        // SaveUploadedFile上传表单文件到指定的路径
        c.SaveUploadedFile(f, "/tmp/"+f.Filename)

        // if err := c.BindJSON(&para); err != nil {
        //     log.Printf("Gin BindJSON failed.")
        // } else {
        //     log.Printf("API upload, request parameters ->\n\tbucket: %s", para.Bucket)
        // }
        // 转存于S3 bucket
        sess := sessionInit_v2("cn-northwest-1")
        s3bucketUploadFile(sess, "qa-is-bucket", "/tmp/"+f.Filename)
        c.JSON(200, gin.H{
            "msg": f,
        })
    }
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
    r.GET("/api/listitems", gin_handler_api_listitems)
    r.POST("/api/upload", gin_handler_api_upload)

    ginLambda = ginadapter.New(r)

}

func APIGW_Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // If no name is provided in the HTTP request body, throw an error
    return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
    lambda.Start(APIGW_Handler)
}
