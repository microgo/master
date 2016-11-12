package utils

import (
	"encoding/json"
	"fmt"
	"gopkg.in/jmcvetta/napping.v3"
	"time"
)

type Payload struct {
	Text string `json:"text"`
}

func SendLogger(hook string, a ...interface{}) error {
	content := time.Now().String() + fmt.Sprint(a)
	payload := Payload{Text: content}
	_, err := napping.Post(hook, &payload, nil, nil)
	return err
}

func LogInfo(a ...interface{}) {
	content := fmt.Sprint(a)
	fmt.Println("[INFO]", time.Now(), content)
}

func LogWarning(a ...interface{}) {
	content := fmt.Sprint(a)
	fmt.Println("[WARN]", time.Now(), content)
}

func LogError(a ...interface{}) {
	content := fmt.Sprint(a)
	fmt.Println("[ERROR]", time.Now(), content)
}

func LogJson(a interface{}) {
	bytes, _ := json.Marshal(a)
	fmt.Println(string(bytes))
}
