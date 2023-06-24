## Create IAM Role
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:CreateLogGroup",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws-cn:logs:cn-northwest-1:523497193792:log-group:/aws/lambda/apiauth-app-demo0:*"
            ]
        }
    ]
}


## Create Lambda function.
```
aws lambda create-function --function-name apiauth-app-demo0 \
  --zip-file fileb://deployment.zip \
  --runtime go1.x --handler main \
  --role arn:aws-cn:iam::523497193792:role/qa-is-apiauth-0 \
  --region cn-northwest-1 \
  --profile default
```

## Update Lambda function.
```
aws lambda update-function-code --function-name apiauth-app-demo0 \
  --zip-file fileb://deployment.zip \
  --region cn-northwest-1 \
  --profile default
```