package main

import (
	// "encoding/json"
	"fmt"
	//_ "github.com/hyperledger/fabric-samples/modelworker/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hyperledger/fabric-samples/modelworker/mysqlconnect"
	"github.com/hyperledger/fabric-samples/modelworker/sdkInit"
	"github.com/hyperledger/fabric-samples/modelworker/service"
	"github.com/hyperledger/fabric-samples/modelworker/web"
	"github.com/hyperledger/fabric-samples/modelworker/web/controller"
	// "github.com/robfig/cron"
	"os"
)

const (
	configFile = "config.yaml"
	SimpleCC   = "mwcc"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID:     "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hyperledger/fabric-samples/modelworker/fixtures/channel-artifacts/channel.tx",

		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: "orderer.example.com",

		ChaincodeID:     SimpleCC,
		ChaincodeGoPath: "/opt/gopath",
		ChaincodePath:   "github.com/hyperledger/fabric-samples/modelworker/chaincode/",
		UserName:        "User1",
	}

	//初始化SDK
	sdk, err := sdkInit.SetupSDK(configFile)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	//关闭sdk节省资源
	defer sdk.Close()

	//创建频道
	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//安装并实例化链码
	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(channelClient)

	serviceSetup := service.ServiceSetup{
		ChaincodeID: SimpleCC,
		Client:      channelClient,
	}

	mainDb, err := mysqlconnect.ConnectMysql()
	if err != nil {
		fmt.Println("数据库连接测试失败")
	} else {
		fmt.Println("数据库连接测试成功")
	}
	mainDb.Close()

	//msg, err := serviceSetup.InitLedger()
	_, err = serviceSetup.InitLedger()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		//fmt.Println(msg)
	}

	// result, err := serviceSetup.FindModelWorkerInfoByWorkerid("13")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	var mw service.ModelWorker
	// 	json.Unmarshal(result, &mw)
	// 	fmt.Println("根据身份证号码查询信息成功：")
	// 	fmt.Println(mw)
	// }

	// //修改信息
	// info := service.ModelWorker{
	// 	Worker_id:   13,
	// 	Worker_name: "张三",
	// 	Remark:      "aa",
	// 	Del_flag:    "1",
	// 	Create_by:   "lifeng",
	// 	Create_time: "2020-11-11 12:00:00",
	// 	Update_by:   "3",
	// 	Update_time: "2020-11-03 00:00:00",
	// 	Gid:         "6fdbaba6-248b-44c3-8a43-ff3de273f9e9",
	// }

	// msg, err = serviceSetup.ModifyModelWorker(info)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("信息操作成功, 交易编号为: " + msg)
	// }

	// result, err = serviceSetup.FindModelWorkerInfoByWorkerid("13")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	var mw service.ModelWorker
	// 	json.Unmarshal(result, &mw)
	// 	fmt.Println("根据身份证号码查询信息成功：")
	// 	fmt.Println(mw)
	// }

	app := controller.Application{
		Setup: &serviceSetup,
	}
	web.WebStart(app)
	fmt.Println("Web Start")

	// c := cron.New(cron.WithSeconds())
	// spec := "00 00 00 * * ?" //cron表达式，每天24点一次
	// c.AddFunc(spec, func() {
	// 	app.UpdateFailureInformation()
	// })
	// c.Start()
	// select {}
}

// 	//修改信息
// 	info := service.ModelWorker{
// 		Name:            "王五",
// 		EntityID:        "370000190000000002",
// 		WorkPlace:       "XX大学研究生院",
// 		ModelWorkerType: "type3",
// 		Status:          "status3",
// 		Auditor:         "lf",
// 		AuditDate:       "2020.09.23",
// 	}

// 	msg, err = serviceSetup.ModifyModelWorker(info)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	} else {
// 		fmt.Println("信息操作成功, 交易编号为: " + msg)
// 	}

// 	result, err = serviceSetup.QueryModelWorker("370000190000000000")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	} else {
// 		var mw service.ModelWorker
// 		json.Unmarshal(result, &mw)
// 		fmt.Println("根据身份证号码查询信息成功：")
// 		fmt.Println(mw)
// 	}

// 	modelworker1 := service.ModelWorker{
// 		Name:            "赵四",
// 		EntityID:        "370000190000000006",
// 		WorkPlace:       "XX大学研究生院",
// 		ModelWorkerType: "type3",
// 		Status:          "status3",
// 		Auditor:         "lf",
// 		AuditDate:       "2020.10.29",
// 	}

// 	msg, err = serviceSetup.SaveModelWorker(modelworker1)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	} else {
// 		fmt.Println("信息发布成功, 交易编号为: " + msg)
// 	}

// 	result, err = serviceSetup.QueryModelWorker("370000190000000006")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	} else {
// 		var mw service.ModelWorker
// 		json.Unmarshal(result, &mw)
// 		fmt.Println("根据身份证号码查询信息成功：")
// 		fmt.Println(mw)
// 	}

// 	app := controller.Application{
// 		Setup: &serviceSetup,
// 	}
// 	web.WebStart(app)
// 	fmt.Println("Web Start")
// }
