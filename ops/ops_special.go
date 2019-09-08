package ops

import (
	"errors"
	"fmt"
)

type specialOps struct {
	baseOps
	sendConfig bool
}

func NewSpecialOps(conf map[string]interface{}) (OpsImpl, error) {
	fmt.Println("SpecialOps: Create")
	sendConfig, ok := conf["SendConfig"]
	if !ok {
		return nil, errors.New("[SendConfig] has not been set in config map")
	}
	return &specialOps{
		sendConfig: sendConfig.(bool),
	}, nil
}

func (ops *specialOps) SendHeartbeat() error {
	fmt.Println("SpecialOps: SendHeartbeat")
	return nil
}

func (ops *specialOps) DoConfigUpload() error {
	fmt.Println("SpecialOps: DoConfigUpload")
	if ops.sendConfig {
		fmt.Println("SpecialOps: upload config")
	} else {
		fmt.Println("SpecialOps: no need to upload config")
	}
	return nil
}
