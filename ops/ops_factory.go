package ops

import (
	"errors"
	"fmt"
	"strings"
)

type opsTypeEnum string

const (
	BaseType    opsTypeEnum = "BaseType"
	SpecialType opsTypeEnum = "SpecialType"
)

type OpsImpl interface {
	SendHeartbeat() error
	DoUpdate() error
	DoConfigUpload() error
}

type OpsFactory func(conf map[string]interface{}) (OpsImpl, error)

var opsFactories = make(map[opsTypeEnum]OpsFactory)

func init() {
	RegisterOps(BaseType, NewBaseOps)
	RegisterOps(SpecialType, NewSpecialOps)
}

func RegisterOps(opsType opsTypeEnum, factory OpsFactory) {
	if factory == nil {
		panic(fmt.Sprintf("Ops factory %s does not exist", string(opsType)))
	}
	_, ok := opsFactories[opsType]
	if ok {
		fmt.Printf("Ops factory %s has been registered already\n", string(opsType))
	} else {
		fmt.Printf("Register ops factory %s\n", string(opsType))
		opsFactories[opsType] = factory
	}
}

func CreateOps(conf map[string]interface{}) (OpsImpl, error) {
	opsType, ok := conf["OpsType"]
	if !ok {
		fmt.Println("No ops type in config map. Use base ops type as default.")
		opsType = BaseType
	}
	OpsFactory, ok := opsFactories[opsType.(opsTypeEnum)]
	if !ok {
		availableOps := make([]string, len(opsFactories))
		for k, _ := range opsFactories {
			availableOps = append(availableOps, string(k))
		}
		return nil, errors.New(fmt.Sprintf("Invalid ops type. Must be one of: %s", strings.Join(availableOps, ", ")))
	}
	fmt.Println("Create ops: ", opsType)
	return OpsFactory(conf)
}
