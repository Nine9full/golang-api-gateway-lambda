
## Lambda Custom Policy

```json
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"states:SendTaskSuccess",
				"states:SendTaskFailure",
				"states:SendTaskHeartbeat",
        "states:DescribeExecution"
			],
			"Resource": "arn:aws:states:ap-southeast-1:XXXXXXXXXXXX:stateMachine:lambda-golang-sqs-standart"
		}
	]
}
```

## Step Function Standart

```json
{
  "Comment": "A description of my state machine",
  "StartAt": "SQS SendMessage",
  "States": {
    "SQS SendMessage": {
      "Type": "Task",
      "Resource": "arn:aws:states:::aws-sdk:sqs:sendMessage.waitForTaskToken",
      "Parameters": {
        "MessageBody": {
          "input.$": "$",
          "MyTaskToken.$": "$$.Task.Token"
        },
        "QueueUrl": "https://sqs.ap-southeast-1.amazonaws.com/XXXXXXXXXXXX/proof-of-concept-sqs"
      },
      "Next": "Pass"
    },
    "Pass": {
      "Type": "Pass",
      "End": true
    }
  },
  "TimeoutSeconds": 600
}
```


## Step Function Express

```json
{
  "Comment": "A description of my state machine",
  "StartAt": "Step Functions StartExecution",
  "States": {
    "Step Functions StartExecution": {
      "Type": "Task",
      "Resource": "arn:aws:states:::states:startExecution",
      "Parameters": {
        "StateMachineArn": "arn:aws:states:ap-southeast-1:XXXXXXXXXXXX:stateMachine:lambda-golang-sqs-standart",
        "Input.$": "$"
      },
      "Next": "Wait X Seconds"
    },
    "Wait X Seconds": {
      "Type": "Wait",
      "Next": "DescribeExecution",
      "Seconds": 1
    },
    "DescribeExecution": {
      "Type": "Task",
      "Parameters": {
        "ExecutionArn.$": "$.ExecutionArn"
      },
      "Resource": "arn:aws:states:::aws-sdk:sfn:describeExecution",
      "Next": "Job Complete?"
    },
    "Job Complete?": {
      "Type": "Choice",
      "Choices": [
        {
          "Variable": "$.Status",
          "StringEquals": "FAILED",
          "Next": "Job Failed"
        },
        {
          "Variable": "$.Status",
          "StringEquals": "SUCCEEDED",
          "Next": "Pass: Job Success"
        },
        {
          "Variable": "$.Status",
          "StringEquals": "TIMED_OUT",
          "Next": "Job Timed out"
        },
        {
          "Variable": "$.Status",
          "StringEquals": "ABORTED",
          "Next": "Job Abort"
        }
      ],
      "Default": "Wait X Seconds"
    },
    "Pass: Job Success": {
      "Type": "Pass",
      "Next": "Job Success"
    },
    "Job Success": {
      "Type": "Succeed"
    },
    "Job Failed": {
      "Comment": "Placeholder for a state which handles the failure.",
      "Type": "Pass",
      "End": true
    },
    "Job Timed out": {
      "Type": "Pass",
      "End": true
    },
    "Job Abort": {
      "Type": "Pass",
      "End": true
    }
  }
}
```

## Step Function Express Trusted entities Policy

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": [
                    "states.amazonaws.com",
                    "apigateway.amazonaws.com"
                ]
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```


## API Gateway Mapping Request
```json
{
    "input": "$util.escapeJavaScript($input.body)",
    "name": "ExecutionFromAPIGateway",
    "stateMachineArn": "arn:aws:states:ap-southeast-1:XXXXXXXXXXXX:stateMachine:lambda-golang-sqs-express"
}
```

## API Gateway Mapping Response

```json
{
    "status": "$input.path('$.status')",
    "startDate": "$input.path('$.startDate')",
    "stopDate": "$input.path('$.stopDate')",
    "data": $util.parseJson($input.path('$.output')).Output
}
```


Body Test

```json
{
    "type": "W100", 
    "member": {
        "lineUserId": "U45d0514a18f60486b6bc7XXXXXXXXX",
        "cardList": [{"cardId": "312000252098"}]
    }
}
```