
package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
)

//type MyReqPara struct {
//    Filename  string `json:filename"`
//}

/*
参考链接: https://blog.csdn.net/weixin_46618592/article/details/125570527

curl -XPOST -F "file=@/tmp/test.txt" \
  -H "Content-Type: multipart/form-data" \
  http://127.0.0.1:8080/api/upload
*/
func gin_handler_api_upload(c *gin.Context) {
    //FormFile返回所提供的表单键的第一个文件
    f, err := c.FormFile("file")
    if err != nil {
        c.String(http.StatusOK, "Gin FormFile Failed.\n")
    } else {
        //SaveUploadedFile上传表单文件到指定的路径
        c.SaveUploadedFile(f, "./"+f.Filename)
        c.JSON(200, gin.H{
            "msg": f,
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