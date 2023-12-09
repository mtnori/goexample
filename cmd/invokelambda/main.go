package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"goexample/internal/proxyrolecreds"
	"log"
	"net"
	"net/http"
	"net/url"
)

func main() {
	// Create AWS SDK config
	// TODO: Use environment variables
	sdkConfig := NewSDKConfig("dummy", "dummy", "dummy")

	// Create Lambda client wrapper
	wrapper := NewFunctionWrapper(sdkConfig)

	// Invoke Lambda functions
	output1 := wrapper.Invoke("hello-world-sample", `{"key1":"aaa"}`, true)
	log.Println(string(output1.Payload))
	output2 := wrapper.Invoke("hello-world-sample", `{"key1":"bbb"}`, true)
	log.Println(string(output2.Payload))
	output3 := wrapper.Invoke("hello-world-sample", `{"key1":"ccc"}`, true)
	log.Println(string(output3.Payload))
}

func NewSDKConfig(proxyHost, proxyPort, proxyRoleName string) aws.Config {
	proxyUrl, _ := url.Parse("http://" + net.JoinHostPort(proxyHost, proxyPort))

	//httpClient := &http.Client{
	//	Transport: &http.Transport{
	//		Proxy: http.ProxyURL(proxyUrl),
	//	},
	//}

	customClient := awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
		tr.Proxy = http.ProxyURL(proxyUrl)
	})

	cfg, _ := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(
			proxyrolecreds.New(
				"http://169.254.169.254/latest/meta-data/iam/security-credentials/"+proxyRoleName,
				proxyrolecreds.WithHTTPClient(customClient),
			),
		),
		config.WithHTTPClient(customClient),
	)

	return cfg
}

func NewFunctionWrapper(sdkConfig aws.Config) FunctionWrapper {
	return FunctionWrapper{
		LambdaClient: lambda.NewFromConfig(sdkConfig),
	}
}

type FunctionWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper *FunctionWrapper) Invoke(functionName string, parameters any, getLog bool) *lambda.InvokeOutput {
	logType := types.LogTypeNone
	if getLog {
		logType = types.LogTypeTail
	}
	payload, err := json.Marshal(parameters)
	if err != nil {
		log.Panicf("Couldn't marshal parameters to JSON. Here's why %v\n", err)
	}

	invokeOutput, err := wrapper.LambdaClient.Invoke(context.Background(), &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		LogType:      logType,
		Payload:      payload,
	})
	if err != nil {
		log.Panicf("Couldn't invoke function %v. Here's why: %v\n", functionName, err)
	}
	return invokeOutput
}
