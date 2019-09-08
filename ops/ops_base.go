package ops

import (
	"errors"
	"fmt"
)

type baseOps struct {
	postUrl string
}

func NewBaseOps(conf map[string]interface{}) (OpsImpl, error) {
	fmt.Println("BaseOps: Create")
	postUrl, ok := conf["PostUrl"]
	if !ok {
		return nil, errors.New("[postUrl] has not been set in config map")
	}
	return &baseOps{
		postUrl: postUrl.(string),
	}, nil
}

func (base *baseOps) SendHeartbeat() error {
	fmt.Println("BaseOps: Send heartbeat")
	fmt.Println("BaseOps: send to url: ", base.postUrl)
	return nil
}

func (base *baseOps) DoUpdate() error {
	fmt.Println("BaseOps: DoUpdate")
	return nil
}

func (base *baseOps) DoConfigUpload() error {
	fmt.Println("BaseOps: DoConfigUpload")
	return nil
}
