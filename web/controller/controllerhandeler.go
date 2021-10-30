package controller

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-samples/modelworker/service"
	"io"
	"net/http"
	"strconv"
	"time"
)

var cuser User
var tempWorkerid string
var tempuser string

//var modifystring string = ""

type Ret struct {
	Description string
	Code        int
	Msg         string
}

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "login.html", nil)
}

func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
	}{
		CurrentUser: cuser,
	}
	ShowView(w, r, "index.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")
	tempuser = loginName
	var flag bool
	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			cuser = user
			flag = true
			break
		}
	}

	data := &struct {
		CurrentUser User
		Flag        bool
	}{
		CurrentUser: cuser,
		Flag:        false,
	}

	if flag {
		// 登录成功
		ShowView(w, r, "index.html", data)
	} else {
		// 登录失败
		data.Flag = true
		data.CurrentUser.LoginName = loginName
		ShowView(w, r, "login.html", data)
	}
}

func (app *Application) UpdateFailureInformation(w http.ResponseWriter, r *http.Request) {

	ret := new(Ret)
	msg, err := app.Setup.UpdateFailureInformation()
	if err != nil {
		ret.Code = 0
		ret.Msg = "区块链数据重传失败,原因:" + err.Error()
		ret_json, _ := json.Marshal(ret)
		io.WriteString(w, string(ret_json))
		fmt.Println(string(ret_json))
		fmt.Println("重传信息上链失败, Code: ", ret.Code, "结果:", ret.Msg)
	} else {
		fmt.Println("上链失败信息添加成功, 交易编号为: " + msg)
		ret.Code = 1
		ret.Msg = "区块链数据重传成功"
		ret_json, _ := json.Marshal(ret)
		io.WriteString(w, string(ret_json))
		fmt.Println(string(ret_json))
		fmt.Println("重传信息上链成功, Code: ", ret.Code, "结果", ret.Msg)
	}
}

//带历史记录
func (app *Application) QueryPage2(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "query2.html", data)
}

//查询劳模及其历史记录
func (app *Application) FindByID(w http.ResponseWriter, r *http.Request) {
	workerID := r.FormValue("workerid")
	tempWorkerid = workerID
	result, err := app.Setup.FindModelWorkerInfoByWorkerid(tempWorkerid)
	// fmt.Println(result)
	var mw = service.ModelWorker{}
	json.Unmarshal([]byte(result), &mw)

	data := &struct {
		MW          service.ModelWorker
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		MW:          mw,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	// fmt.Println(data.MW)

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true

		// return "query failed"
	}

	ShowView(w, r, "queryResult.html", data)

	// return "query success"
}

func (app *Application) QueryModelWorker(w http.ResponseWriter, r *http.Request) {
	workerID := r.FormValue("workerid")
	tempWorkerid = workerID
	result, err := app.Setup.FindModelWorkerInfoByWorkerid(tempWorkerid)
	if err != nil {
		warningString, _ := json.Marshal("查询失败,不存在该信息!")
		io.WriteString(w, string(warningString))
	}
	io.WriteString(w, string(result))

	// return "query success" b
}

// 显示添加信息页面
func (app *Application) AddModelWorkerShow(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
		//		Modifystring string
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		//		Modifystring: modifystring,
	}
	ShowView(w, r, "addModelWorker.html", data)
}

// 添加信息
func (app *Application) AddModelWorker(w http.ResponseWriter, r *http.Request) {
	timeObj := time.Now()
	var strtime = timeObj.Format("2006.01.02 15:04:05")
	// currentuser := tempuser
	currentuser := r.FormValue("updateBy")
	stringworkerid := r.FormValue("workerId")
	intworkerid, err := strconv.Atoi(stringworkerid)

	if err != nil {
		fmt.Errorf("字符串转int失败!")
		// return "0"
	}
	mw := service.ModelWorker{
		Worker_id:   intworkerid,
		Worker_name: r.FormValue("workerName"),
		Remark:      r.FormValue("remark"),
		Del_flag:    r.FormValue("delFlag"),
		Create_by:   r.FormValue("createBy"),
		Create_time: r.FormValue("createTime"),
		Update_by:   currentuser,
		Update_time: strtime,
		Gid:         r.FormValue("gid"),
	}

	fmt.Println(mw.Worker_id, mw.Worker_name, mw.Remark, mw.Del_flag, mw.Create_by, mw.Create_time, mw.Update_by, mw.Update_time, mw.Gid)
	msg, err := app.Setup.SaveModelWorker(mw)
	ret := new(Ret)
	if err != nil {
		fmt.Println(err.Error())
		// return "0"
		//modifystring = "添加信息失败"
		// app.AddModelWorkerShow(w, r)
		ret.Description = "链上不存在该信息,执行添加操作"
		ret.Code = 0
		ret.Msg = "数据上链失败,原因:" + err.Error()
		ret_json, _ := json.Marshal(ret)
		io.WriteString(w, string(ret_json))
		fmt.Println(string(ret_json))
		fmt.Println("添加信息发布失败, Code: ", ret.Code, "结果:", ret.Msg)
		//modifystring = ""

	} else {
		//fmt.Println("添加信息发布成功, 交易编号为: " + msg)
		// return "1"
		ret.Description = "链上不存在该信息,执行添加操作"
		ret.Code = 1
		ret.Msg = "数据上链成功"
		ret_json, _ := json.Marshal(ret)
		io.WriteString(w, string(ret_json))
		fmt.Println(string(ret_json))
		// r.Form.Set("workerid", stringworkerid)
		fmt.Println("添加信息发布成功, Code: ", ret.Code, "结果", ret.Msg, "交易编号为: "+msg, strtime)
		// app.FindByID(w, r)
		//modifystring = ""
	}

}

// 修改/添加新信息
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request) {
	// 根据证书编号与姓名查询信息
	workerid := tempWorkerid
	result, err := app.Setup.FindModelWorkerInfoByWorkerid(workerid)

	var mw = service.ModelWorker{}
	json.Unmarshal([]byte(result), &mw)

	data := &struct {
		ModelWorker service.ModelWorker
		CurrentUser User
		Msg         string
		Flag        bool
		//		Modifystring string
	}{
		ModelWorker: mw,
		CurrentUser: cuser,
		Flag:        true,
		Msg:         "",
		//		Modifystring: modifystring,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "modify.html", data)
}

func (app *Application) Modify(w http.ResponseWriter, r *http.Request) {
	timeObj := time.Now()
	var strtime = timeObj.Format("2006.01.02 15:04:05")
	// currentuser := tempuser
	currentuser := r.FormValue("updateBy")
	// workerid := tempWorkerid
	workerid := r.FormValue("workerId")

	//intworkerid, err := strconv.Atoi(workerid)
	intworkerid, err := strconv.Atoi(workerid)
	if err != nil {
		fmt.Errorf("字符串转int失败!")
		// return "0"
	}
	mw := service.ModelWorker{
		Worker_id:   intworkerid,
		Worker_name: r.FormValue("workerName"),
		Remark:      r.FormValue("remark"),
		Del_flag:    r.FormValue("delFlag"),
		Create_by:   r.FormValue("createBy"),
		Create_time: r.FormValue("createTime"),
		Update_by:   currentuser,
		Update_time: strtime,
		Gid:         r.FormValue("gid"),
	}

	fmt.Println(mw.Worker_id, mw.Worker_name, mw.Remark, mw.Del_flag, mw.Create_by, mw.Create_time, mw.Update_by, mw.Update_time, mw.Gid)
	msg, err := app.Setup.ModifyModelWorker(mw)
	ret := new(Ret)
	if err != nil {
		fmt.Println(err.Error())
		ret.Description = "链上已存在该信息,执行修改操作"
		ret.Code = 0
		ret.Msg = "区块链数据修改失败,原因:" + err.Error()
		ret_json, _ := json.Marshal(ret)
		io.WriteString(w, string(ret_json))
		// fmt.Println(string(ret_json))
		fmt.Println("修改信息发布失败, Code: ", ret.Code, "结果:", ret.Msg, strtime)
	} else {
		//fmt.Println("修改信息发布成功, 交易编号为: " + msg)
		ret.Description = "链上已存在该信息,执行修改操作"
		ret.Code = 1
		ret.Msg = "区块链数据修改成功"
		ret_json, _ := json.Marshal(ret)
		io.WriteString(w, string(ret_json))
		// fmt.Println(string(ret_json))
		fmt.Println("修改信息发布成功, Code: ", ret.Code, "结果", ret.Msg, "交易编号为: "+msg)
		//r.Form.Set("workerid", workerid)
		//app.FindByID(w, r)
	}
}

func (app *Application) AddOrModify(w http.ResponseWriter, r *http.Request) {
	timeObj := time.Now()
	var strtime = timeObj.Format("2006.01.02 15:04:05")
	inputWorkerid := r.FormValue("workerId")
	fmt.Println(r.FormValue("workerId"), r.FormValue("workerName"), r.FormValue("remark"),
		r.FormValue("delFlag"), r.FormValue("createBy"), r.FormValue("createTime"),
		r.FormValue("updateBy"), r.FormValue("gid"), strtime)
	intinputWorkerid, _ := strconv.Atoi(inputWorkerid)
	result, _ := app.Setup.FindModelWorkerInfoByWorkerid(inputWorkerid)
	// if err != nil {
	// 	errorstring := "查询操作失败,原因为" + err.Error()
	// 	io.WriteString(w, errorstring)
	// 	// return "0"
	// }
	ret := new(Ret)

	var mw = service.ModelWorker{}
	json.Unmarshal([]byte(result), &mw)
	fmt.Println("####################################################")
	fmt.Println(mw.Worker_id)
	fmt.Println("####################################################")
	if intinputWorkerid == 0 {
		fmt.Println("信息导入失败,传入id为0")
		ret.Description = "信息导入失败,传入id为0"
		ret.Code = 0
		ret.Msg = "区块链数据操作失败,原因:传入id为0"
		ret_json, _ := json.Marshal(ret)
		io.WriteString(w, string(ret_json))
	} else if mw.Worker_id == 0 {
		fmt.Println("链上不存在该信息,执行添加操作")
		app.AddModelWorker(w, r)
		fmt.Println("####################################################")
		fmt.Println("####################################################")
	} else {
		fmt.Println("链上已存在该信息,执行修改操作")
		app.Modify(w, r)
		fmt.Println("####################################################")
		fmt.Println("####################################################")
	}
}
