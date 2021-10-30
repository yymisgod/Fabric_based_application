package main

import (
	//"bytes"
	"encoding/json"
	//"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

func PutModelWorker(APIstub shim.ChaincodeStubInterface, mw ModelWorker) ([]byte, bool) {

	b, err := json.Marshal(mw)
	if err != nil {
		return nil, false
	}

	//保存modelworker状态
	stringworkerid := strconv.Itoa(mw.Worker_id)
	err = APIstub.PutState(stringworkerid, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

//获得信息
func GetModelWorkerInfo(APIstub shim.ChaincodeStubInterface, Worker_id string) (ModelWorker, bool) {
	var mw ModelWorker

	b, err := APIstub.GetState(Worker_id)
	if err != nil {
		return mw, false
	}

	if b == nil {
		return mw, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &mw)
	if err != nil {
		return mw, false
	}

	// 返回结果
	return mw, true
}

func (t *ModelWorkerChaincode) InitLedger(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	// modelworkers := []ModelWorker{}
	modelworker := ModelWorker{}
	err := json.Unmarshal([]byte(args[0]), &modelworker)

	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// i := 0
	// for i < len(modelworkers) {
	// 	workerAsBytes, _ := json.Marshal(modelworkers[i])
	// 	stringworkerid := strconv.Itoa(modelworkers[i].Worker_id)
	// 	APIstub.PutState(stringworkerid, workerAsBytes)
	// 	i = i + 1
	// }

	_, bl := PutModelWorker(APIstub, modelworker)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}
	err = APIstub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ModelWorkerChaincode) QueryModelWorkerHistory(APIstub shim.ChaincodeStubInterface, query string) peer.Response {
	// 根据身份证号码查询状态
	workerAsBytes, err := APIstub.GetState(query)

	if err != nil {
		return shim.Error("出现错误")
	}
	if workerAsBytes == nil {
		return shim.Error("未找到该信息")
	}

	// 对查询到的状态进行反序列化
	var mw ModelWorker
	err = json.Unmarshal(workerAsBytes, &mw)
	if err != nil {
		return shim.Error("反序列化modelworker信息失败")
	}

	// 获取历史变更数据
	stringworkerid := strconv.Itoa(mw.Worker_id)
	iterator, err := APIstub.GetHistoryForKey(stringworkerid)
	if err != nil {
		return shim.Error("根据指定的工号查询对应的历史变更数据失败")
	}

	defer iterator.Close()

	// 迭代处理
	var historys []HistoryItem
	var hisMW ModelWorker
	for iterator.HasNext() {
		hisData, err := iterator.Next()

		if err != nil {
			return shim.Error("获取modelworker的历史变更数据失败")
		}

		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisMW)

		if hisData.Value == nil {
			var empty ModelWorker
			historyItem.ModelWorker = empty
		} else {
			historyItem.ModelWorker = hisMW
		}

		historys = append(historys, historyItem)

	}

	mw.Historys = historys
	// 返回
	result, err := json.Marshal(mw)
	// fmt.Println(result)
	if err != nil {
		return shim.Error("序列化modelworker信息时发生错误")
	}

	return shim.Success(result)
}

//增加
func (t *ModelWorkerChaincode) addModelWorker(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var mw ModelWorker

	err := json.Unmarshal([]byte(args[0]), &mw)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	stringworkerid := strconv.Itoa(mw.Worker_id)
	_, exist := GetModelWorkerInfo(APIstub, stringworkerid)
	if exist {
		return shim.Error("要添加的工号已存在")
	}

	_, bl := PutModelWorker(APIstub, mw)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err = APIstub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))

}

// 根据身份证号更新信息
// args:
func (mw *ModelWorkerChaincode) updateModelWorker(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	var info ModelWorker
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return shim.Error("反序列化modelworker信息失败")
	}

	stringworkerid := strconv.Itoa(info.Worker_id)
	result, bl := GetModelWorkerInfo(APIstub, stringworkerid)
	if !bl {
		return shim.Error("根据工号查询信息时发生错误")
	}

	result.Worker_id = info.Worker_id
	result.Worker_name = info.Worker_name
	result.Remark = info.Remark
	result.Del_flag = info.Del_flag
	result.Create_by = info.Create_by
	result.Create_time = info.Create_time
	result.Update_by = info.Update_by
	result.Update_time = info.Update_time
	result.Gid = info.Gid

	stringworkerid = strconv.Itoa(result.Worker_id)
	WorkerAsBytes, _ := json.Marshal(result)
	APIstub.PutState(stringworkerid, WorkerAsBytes)

	err = APIstub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
