package main

import (
	"factory_pattern/ops"
	"fmt"
)

func main() {
	baseOps, err := ops.CreateOps(map[string]interface{}{
		"OpsType": ops.BaseType,
		"PostUrl": "http://ops.cloud.com/send_heartbeat",
	})
	if err != nil {
		fmt.Println("create baseOps failed, err: ", err.Error())
		return
	}
	baseOps.DoConfigUpload() // Output: BaseOps: DoConfigUpload
	baseOps.DoUpdate()       // Output: BaseOps: DoUpdate
	baseOps.SendHeartbeat()  // Output: BaseOps: Send heartbeat

	specialOps, err := ops.CreateOps(map[string]interface{}{
		"OpsType":    ops.SpecialType,
		"SendConfig": true,
	})
	if err != nil {
		fmt.Println("create specialOps failed, err: ", err.Error())
		return
	}
	specialOps.DoConfigUpload() // Output: SpecialOps: DoConfigUpload
	specialOps.DoUpdate()       // Output: BaseOps: DoUpdate
	specialOps.SendHeartbeat()  // Output: SpecialOps: SendHeartbeat
}
