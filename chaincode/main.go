package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type ModelWorkerChaincode struct {
}

func (mw *ModelWorkerChaincode) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (mw *ModelWorkerChaincode) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {

	function, args := APIstub.GetFunctionAndParameters()
	if function == "InitLedger" { //初始化账本
		return mw.InitLedger(APIstub, args)
	} else if function == "addModelWorker" { //新建劳模
		return mw.addModelWorker(APIstub, args)
	} else if function == "QueryModelWorkerHistory" { //根据身份证号查询带历史记录
		return mw.QueryModelWorkerHistory(APIstub, args[0])
	} else if function == "updateModelWorker" { //更新劳模信息
		return mw.updateModelWorker(APIstub, args)
	}

	return shim.Error("不存在该方法")
}

func main() {
	err := shim.Start(new(ModelWorkerChaincode))
	if err != nil {
		fmt.Printf("创建链码失败，错误 :%s", err.Error())
	}
}
