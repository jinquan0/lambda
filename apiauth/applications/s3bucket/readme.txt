## IAM Role policy[qa-is-apiauth]

{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "lambdaPermission",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:CreateLogGroup",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws-cn:logs:cn-northwest-1:260547014985:log-group:/aws/lambda/edison-apiauth-authorizer:*",
                "arn:aws-cn:logs:cn-northwest-1:260547014985:log-group:/aws/lambda/edison-apiauth-app-demo0:*",
                "arn:aws-cn:logs:cn-northwest-1:260547014985:log-group:/aws/lambda/edison-apiauth-app-s3:*"
            ]
        },
        {
            "Sid": "s3bucketPermission",
            "Effect": "Allow",
            "Action": [
                "s3:*"
            ],
            "Resource": [
                "arn:aws-cn:s3:::qa-is-bucket",
                "arn:aws-cn:s3:::qa-is-bucket/*"
            ]
        },
        {
            "Sid": "ec2Permission",
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeNetworkInterfaces",
                "ec2:CreateNetworkInterface",
                "ec2:DeleteNetworkInterface",
                "ec2:DescribeInstances",
                "ec2:AttachNetworkInterface"
            ],
            "Resource": "*"
        }
    ]
}


## create function
aws lambda create-function --function-name edison-apiauth-app-s3 \
  --zip-file fileb://deployment.zip \
  --runtime go1.x --handler main \
  --role arn:aws-cn:iam::260547014985:role/qa-is-apiauth \
  --region cn-northwest-1 \
  --profile dsmqa


## update function
aws lambda update-function-code --function-name edison-apiauth-app-s3 \
  --zip-file fileb://deployment.zip \
  --region cn-northwest-1 \
  --profile dsmqa



## upload single local file to AWS S3 Bucket[qa-is-bucket]
curl -XPOST -F "file=@/tmp/test.txt" \
  -H "Content-Type: multipart/form-data" \
  https://otbxr6m1fg.execute-api.cn-northwest-1.amazonaws.com.cn/v1/api/upload

curl -XGET -H "Content-Type: application/json" \
-d '{"bucket": "qa-is-bucket"}' \
  https://otbxr6m1fg.execute-api.cn-northwest-1.amazonaws.com.cn/v1/api/listitems



curl -XPOST -F 'token="79faf82271944fe38c4f1d99be71bc9c"' \
  -F 'file=@/tmp/test.txt' \
  -F 'auth={"token": "79faf82271944fe38c4f1d99be71bc9c"}' \
  -H 'Content-Type: multipart/form-data' \
  https://otbxr6m1fg.execute-api.cn-northwest-1.amazonaws.com.cn/v1/api/upload