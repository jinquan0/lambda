
package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "mime/multipart"
)

/*
type SomeStruct struct {
    Token  string `json:token"`
}

type SomeRequest struct {
  File        *multipart.FileHeader `form:"file"`
  StructField SomeStruct            `form:"structField"`
}

// 参考链接: 
// https://blog.csdn.net/weixin_46618592/article/details/125570527
// https://www.appsloveworld.com/go/19/how-to-upload-multipart-file-and-json-in-go-with-gin-gonic

// curl -XPOST -F "file=@/tmp/test.txt" \
//   -F 'structField={ "token": "1234567890" }' \
//   -H "Content-Type: multipart/form-data" \
//   http://127.0.0.1:8080/api/upload

func gin_handler_api_upload(c *gin.Context)  {
    var req SomeRequest
    if err := c.ShouldBind(&req); err != nil {
        c.String(http.StatusOK, "Gin ShouldBind Failed.\n")
    } else {
        log.Printf("token: %s", req.StructField.Token)
        c.SaveUploadedFile(req.File, "./"+req.File.Filename)
        //log.Println(req)
        c.JSON(200, gin.H{
            "msg": req.File,
        })
    }
}
*/

type SomeRequest struct {
  Token  string `form:"token"`
  File   *multipart.FileHeader `form:"file"`
}

// 参考链接: 
// https://blog.csdn.net/weixin_46618592/article/details/125570527
// https://www.appsloveworld.com/go/19/how-to-upload-multipart-file-and-json-in-go-with-gin-gonic

// curl -XPOST -F 'token=1234567890' -F 'file=@/tmp/test.txt' -H 'Content-Type: multipart/form-data' http://127.0.0.1:8080/api/upload
// 或者
// curl -XPOST -F 'token="1234567890"' -F 'file=@/tmp/test.txt' -H 'Content-Type: multipart/form-data' http://127.0.0.1:8080/api/upload

func gin_handler_api_upload(c *gin.Context)  {
    var req SomeRequest
    if err := c.ShouldBind(&req); err != nil {
        c.String(http.StatusOK, "Gin ShouldBind Failed.\n")
    } else {
        log.Printf("token: %s", req.Token)
        c.SaveUploadedFile(req.File, "./"+req.File.Filename)
        c.JSON(200, gin.H{
            "msg": req.File,
        })
    }
}

func gin_handler_healthy(c *gin.Context) {
    // c.JSON(200, gin.H{
    //     "message": "hello, lambda function is ok",
    // })
    
    c.String(http.StatusOK, "hello, lambda function is ok\n")
}

func main() {

    // stdout and stderr are sent to AWS CloudWatch Logs
    log.Printf("Gin cold start")
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hi, its me, the root path[/]\n")
    })

    r.GET("/healthy", gin_handler_healthy)
    r.POST("/api/upload", gin_handler_api_upload)

    r.Run(":8080")
}