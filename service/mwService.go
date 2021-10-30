package service

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hyperledger/fabric-samples/modelworker/mysqlconnect"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"strconv"
	"time"
)

func (t *ServiceSetup) InitLedger() (string, error) {

	DataBase, err := mysqlconnect.ConnectMysql()
	if err != nil {
		fmt.Println("数据库连接失败")
	} else {
		fmt.Println("数据库连接成功")
	}
	fmt.Println("开始进行初始化")
	//Modelworkers := []ModelWorker{}
	DataModelWorker := mysqlconnect.QueryData(DataBase) //变量DataModelWorker是DataModelWorker类型的切片
	//loopAmount := len(DataModelWorker)
	loopAmount := 100
	var errorCode []int
	var retryTime = 100
	var countTime = 101
	for j := 0; j < loopAmount; j++ {

		timeObj := time.Now()
		var strtime = timeObj.Format("2006.01.02 15:04:05")
		TempDataModelWorker := DataModelWorker[j]
		DWorker_id := TempDataModelWorker.Worker_id
		DWorker_name := TempDataModelWorker.Worker_name
		DRemark := TempDataModelWorker.Remark
		DDel_flag := TempDataModelWorker.Del_flag
		DCreate_by := TempDataModelWorker.Create_by
		DCreate_time := TempDataModelWorker.Create_time
		DUpdate_by := TempDataModelWorker.Update_by
		//DUpdate_time := TempDataModelWorker.Update_time
		DGid := TempDataModelWorker.Gid
		var Dmodelworker = ModelWorker{
			Worker_id:   DWorker_id,
			Worker_name: DWorker_name,
			Remark:      DRemark,
			Del_flag:    DDel_flag,
			Create_by:   DCreate_by,
			Create_time: DCreate_time,
			Update_by:   DUpdate_by,
			Update_time: strtime,
			Gid:         DGid}
		//Modelworkers = append(Modelworkers, Dmodelworker)
		_, err := t.InitModelWorker(Dmodelworker)
		errorString := strconv.Itoa(j + 1)
		countTimeString := strconv.Itoa(countTime - retryTime)
		if err != nil {
			fmt.Println("第"+errorString+"条数据初始化失败,错误为:", err.Error())
			if retryTime > 0 {
				fmt.Println("第" + countTimeString + "次" + "重传第" + errorString + "条数据")
				j = j - 1
				retryTime = retryTime - 1
			} else {
				errorCode = append(errorCode, j)
				fmt.Println("背书节点背书100次失败,网络可能出现问题,该信息为第" + errorString + "条")
				//return "背书节点背书100次失败,网络可能出现问题"
			}
		} else {
			fmt.Println("第" + errorString + "条数据初始化成功")
			retryTime = 100
		}
		k := j % 100
		if k == 0 && j != 0 {
			fmt.Println("已初始化100条信息")
		}
	}
	//fmt.Println("数据初始化成功")
	if len(errorCode) != 0 {
		fmt.Println(errorCode)
	}
	return "数据初始化成功", nil
}

func (t *ServiceSetup) FindModelWorkerInfoByWorkerid(Worker_id string) ([]byte, error) { //controller记得改函数名

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "QueryModelWorkerHistory", Args: [][]byte{[]byte(Worker_id)}}

	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	// fmt.Println(response.Payload)
	return response.Payload, nil
}

func (t *ServiceSetup) InitModelWorker(mw ModelWorker) (string, error) {

	eventID := "eventInitLedger"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(mw)
	if err != nil {
		return "", fmt.Errorf("指定的mw对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "InitLedger", Args: [][]byte{b, []byte(eventID)}}
	//respone, err := t.Client.Execute(req)
	_, err = t.Client.Execute(req)
	if err != nil {
		return "初始化失败", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "初始化失败", err
	}

	//return string(respone.TransactionID), nil
	return "", nil
}

//创建新的
func (t *ServiceSetup) SaveModelWorker(mw ModelWorker) (string, error) {

	eventID := "eventAddModelWorker"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(mw)
	if err != nil {
		return "", fmt.Errorf("指定的mw对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addModelWorker", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	// Db, err := mysqlconnect.ConnectMysql()
	// if err != nil {
	// 	return ("数据库连接失败"), err
	// } else {
	// 	fmt.Println("数据库连接成功")

	// 	stringworkerid := strconv.Itoa(mw.Worker_id)
	// 	tempString := []string{stringworkerid, mw.Worker_name, mw.Remark, mw.Del_flag, mw.Create_by, mw.Create_time, mw.Update_by, mw.Update_time, mw.Gid}
	// 	mysqlconnect.AddRecord(Db, tempString)
	// 	Db.Close()
	// 	return string(respone.TransactionID), nil
	// }
	return string(respone.TransactionID), nil
}

//更改信息
func (t *ServiceSetup) ModifyModelWorker(mw ModelWorker) (string, error) {
	eventID := "eventModifyMW"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(mw)
	if err != nil {
		return "", fmt.Errorf("指定的mw对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "updateModelWorker", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)

	if err != nil {
		return "", err
	}
	// Db, err := mysqlconnect.ConnectMysql()
	// if err != nil {
	// 	fmt.Println("数据库连接失败")
	// } else {
	// 	fmt.Println("数据库连接成功")
	// }

	// stringworkerid := strconv.Itoa(mw.Worker_id)
	// tempString := []string{stringworkerid, mw.Worker_name, mw.Remark, mw.Del_flag, mw.Create_by, mw.Create_time, mw.Update_by, mw.Update_time, mw.Gid}
	// mysqlconnect.UpdateRecord(Db, tempString)
	// Db.Close()
	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) UpdateFailureInformation() (string, error) {

	DataBase, err := mysqlconnect.ConnectMysql()
	if err != nil {
		fmt.Println("数据库连接失败")
	} else {
		fmt.Println("数据库连接成功")
	}

	var errorCode []int
	var retryTime = 100
	var countTime = 101
	//Modelworkers := []ModelWorker{}
	lenOfFailure := mysqlconnect.QueryFailureCount(DataBase)
	if lenOfFailure == 0 {
		fmt.Println("warning:失败个数为0！")
	}
	DataModelWorker := mysqlconnect.QueryDataFailure(DataBase) //变量DataModelWorker是DataModelWorker类型的切片
	for j := 0; j < lenOfFailure; j++ {
		TempDataModelWorker := DataModelWorker[j]
		DWorker_id := TempDataModelWorker.Worker_id
		DWorker_name := TempDataModelWorker.Worker_name
		DRemark := TempDataModelWorker.Remark
		DDel_flag := TempDataModelWorker.Del_flag
		DCreate_by := TempDataModelWorker.Create_by
		DCreate_time := TempDataModelWorker.Create_time
		DUpdate_by := TempDataModelWorker.Update_by
		DUpdate_time := TempDataModelWorker.Update_time
		DGid := TempDataModelWorker.Gid
		var Dmodelworker = ModelWorker{
			Worker_id:   DWorker_id,
			Worker_name: DWorker_name,
			Remark:      DRemark,
			Del_flag:    DDel_flag,
			Create_by:   DCreate_by,
			Create_time: DCreate_time,
			Update_by:   DUpdate_by,
			Update_time: DUpdate_time,
			Gid:         DGid}
		//Modelworkers = append(Modelworkers, Dmodelworker)

		t.InitModelWorker(Dmodelworker)
		errorString := strconv.Itoa(j + 1)
		countTimeString := strconv.Itoa(countTime - retryTime)
		if err != nil {
			fmt.Println("第"+errorString+"条数据上链失败,错误为:", err.Error())
			if retryTime > 0 {
				fmt.Println("第" + countTimeString + "次" + "重传第" + errorString + "条数据")
				j = j - 1
				retryTime = retryTime - 1
			} else {
				errorCode = append(errorCode, j)
				fmt.Println("背书节点背书100次失败,网络可能出现问题,该信息为第" + errorString + "条")
				//return "背书节点背书100次失败,网络可能出现问题"
			}
		} else {
			fmt.Println("第" + errorString + "条数据初始化成功")
			retryTime = 100
		}
		k := j % 100
		if k == 0 && j != 0 {
			fmt.Println("已上链100条信息")
		}
	}
	//fmt.Println("数据初始化成功")
	if len(errorCode) != 0 {
		fmt.Println(errorCode)
	}
	return "失败数据上链成功", nil
}
