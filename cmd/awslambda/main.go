package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

// eventType は Lambda に送られるイベント種別を表す列挙型
type eventType int

const (
	unknown eventType = iota
	sqs
	direct
)

// Message は Lambda に送られるメッセージボディを表す構造体
type Message struct {
	Name string
}

// Event は Lambda に送られるイベント内容を表す構造体
// Unmarshaler インターフェースを満たす
type Event struct {
	EventType   eventType
	SQSEvent    events.SQSEvent
	DirectEvent Message
}

// UnmarshalJSON は Lambda に送られた JSON 文字列を Event 型に変換する
func (event *Event) UnmarshalJSON(data []byte) error {
	log.Println("call UnmarshalJSON")

	switch event.getEventType(data) {
	case sqs:
		log.Println("switch SQS event type")

		var sqsEvent events.SQSEvent
		if err := json.Unmarshal(data, &sqsEvent); err != nil {
			return err
		}
		event.mapSQSEvent(sqsEvent)
	case direct:
		log.Println("switch direct event type")

		var directEvent Message
		if err := json.Unmarshal(data, &directEvent); err != nil {
			return err
		}
		event.mapDirectEvent(directEvent)
	case unknown:
		return errors.New("unsupported event type")
	}

	return nil
}

// mapSQSEvent は SQS からのイベント内容を Event 構造体にマッピングするための処理
func (event *Event) mapSQSEvent(sqsEvent events.SQSEvent) {
	event.EventType = sqs
	event.SQSEvent = sqsEvent
}

// mapSQSEvent は 直接呼び出しのイベント内容を Event 構造体にマッピングするための処理
func (event *Event) mapDirectEvent(directEvent Message) {
	log.Println("call mapDirectEvent")

	event.EventType = direct
	event.DirectEvent = directEvent
}

// getEventType はイベント種別を返却する
func (event *Event) getEventType(data []byte) eventType {
	// JSON 文字列か確認する
	// JSON 文字列でなければ、イベント種別は不明と判定する
	temp := make(map[string]any)
	if err := json.Unmarshal(data, &temp); err != nil {
		return unknown
	}

	// Records フィールドが存在するかチェックする
	// 存在しなければ、イベント種別を直接実行と判定する
	rs, rsOK := temp["Records"].([]any)
	if !rsOK {
		return direct
	}

	// Records フィールドの一つ目の要素が存在するかチェックする
	// 存在しなければ、イベント種別は不明と判定する
	r, rOK := rs[0].(map[string]any)
	if !rOK {
		return unknown
	}

	// イベントソース名を取得する
	var eventSource string
	if es, ok := r["EventSource"]; ok {
		eventSource = es.(string)
	} else if es, ok := r["eventSource"]; ok {
		eventSource = es.(string)
	}

	// イベントソース名をイベント種別へ変換する
	// 現在は SQS メッセージのみに対応
	switch eventSource {
	case "aws:sqs":
		return sqs
	default:
		return unknown
	}
}

// HandleRequest は Lambda のハンドラ関数
func HandleRequest(_ context.Context, event Event) (string, error) {
	switch event.EventType {
	case sqs:
		log.Println("handle SQS event")
		return "from SQS", nil
	case direct:
		log.Println("handle direct event")
		message := fmt.Sprintf("Hello %s!", event.DirectEvent.Name)
		return message, nil
	}

	return "", errors.New("lambda error")
}

func main() {
	lambda.Start(HandleRequest)
}
