package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	wrapper := NewFunctionWrapper()
	output1 := wrapper.Invoke("hello-world-sample", `{"key1":"aaaa"}`, true)
	log.Println(string(output1.Payload))
	output2 := wrapper.Invoke("hello-world-sample", `{"key1":"bbb"}`, true)
	log.Println(string(output2.Payload))
	output3 := wrapper.Invoke("hello-world-sample", `{"key1":"ccc"}`, true)
	log.Println(string(output3.Payload))
}

type OptionFunc func(*endpointcreds.Options)

func HTTPClient(client endpointcreds.HTTPClient) OptionFunc {
	return func(args *endpointcreds.Options) {
		args.HTTPClient = client
	}
}

func NewFunctionWrapper() FunctionWrapper {
	//proxyCredentials := getCredentials()
	//
	customClient := awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
		proxyURL, err := url.Parse("http://PROXY_HOST:8080")
		if err != nil {
			log.Fatal(err)
		}
		tr.Proxy = http.ProxyURL(proxyURL)
	})
	//
	//cfg, _ := config.LoadDefaultConfig(context.Background(),
	//	config.WithCredentialsProvider(
	//		credentials.NewStaticCredentialsProvider(
	//			proxyCredentials.AccessKeyId,
	//			proxyCredentials.SecretAccessKey,
	//			proxyCredentials.Token)),
	//	config.WithHTTPClient(customClient),
	//)

	cfg, _ := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(
			endpointcreds.New(
				"http://169.254.169.254/latest/meta-data/iam/security-credentials/ROLE_NAME",
				HTTPClient(customClient),
			),
		),
	)

	return FunctionWrapper{
		LambdaClient: lambda.NewFromConfig(cfg),
	}
}

type FunctionWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper FunctionWrapper) Invoke(functionName string, parameters any, getLog bool) *lambda.InvokeOutput {
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

type ProxyCredentials struct {
	AccessKeyId     string    `json:"AccessKeyId"`
	SecretAccessKey string    `json:"SecretAccessKey"`
	Token           string    `json:"Token"`
	Expiration      time.Time `json:"Expiration"`
}

//func getCredentials() ProxyCredentials {
//	req, _ := http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/iam/security-credentials/invoke-lambda-role", nil)
//
//	proxyUrl, _ := url.Parse("PROXY URL")
//	client := &http.Client{
//		Transport: &http.Transport{
//			Proxy: http.ProxyURL(proxyUrl),
//		},
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Panicf("Couldn't access proxy server %v\n", err)
//	}
//	defer resp.Body.Close()
//
//	var proxyCredentials ProxyCredentials
//	if err := json.NewDecoder(resp.Body).Decode(&proxyCredentials); err != nil {
//		log.Panicf("Couldn't decode metadata response %v\n", err)
//	}
//
//	//log.Print(proxyCredentials.AccessKeyId)
//	//log.Print(proxyCredentials.SecretAccessKey)
//	//log.Print(proxyCredentials.Token)
//	log.Print(proxyCredentials.Expiration)
//
//	return proxyCredentials
//}
