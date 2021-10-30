package web

import (
	"fmt"
	"github.com/hyperledger/fabric-samples/modelworker/web/controller"
	"net/http"
)

// 启动Web服务并指定路由信息
func WebStart(app controller.Application) {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 指定路由信息(匹配请求)
	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login", app.Login)

	http.HandleFunc("/index", app.Index)

	http.HandleFunc("/addModelWorkerInfo", app.AddModelWorkerShow) // 显示添加信息页面
	http.HandleFunc("/addModelWorker", app.AddModelWorker)         // 提交信息请求

	http.HandleFunc("/queryPage2", app.QueryPage2) // 转至根据身份证号码查询信息页面
	http.HandleFunc("/query2", app.FindByID)       // 根据身份证号码查询信息（历史记录）

	http.HandleFunc("/modifyPage", app.ModifyShow) // 修改信息页面
	http.HandleFunc("/modify", app.Modify)         //  修改信息

	http.HandleFunc("/addormodify", app.AddOrModify)
	http.HandleFunc("/queryjson", app.QueryModelWorker)

	http.HandleFunc("/updatefailureinformation", app.UpdateFailureInformation)

	fmt.Println("启动Web服务, 监听端口号为: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v", err)
	}

}
