package mysqlconnect

import (
	// "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DataModelWorker struct {
	audit_id    int    `db:"audit_id"`
	Worker_id   int    `db:"worker_id"`
	Worker_name string `db:"worker_name"`
	Remark      string `db:"remark"`
	Del_flag    string `db:"del_flag"`
	Create_by   string `db:"create_by"`
	Create_time string `db:"create_time"`
	Update_by   string `db:"update_by"`
	Update_time string `db:"update_time"`
	Gid         string `db:"gid"`
}

type DataFailedInformation struct {
	id          int `db:"id"`
	upload_flag int `db:"upload_flag"`
}

var (
	//userName string = "root"
	userName string = "qkl"
	//password string = "root"
	password string = "CDfUp2vgXchOatuY"
	//ipAddrees string = "192.168.0.132"
	ipAddrees string = "10.10.10.2" //mysql数据库服务器ip
	port      int    = 3306
	//dbName    string = "ModelWorker"
	dbName  string = "laomo_dev_new"
	charset string = "utf8"
	// tableName string = "ModelWorkers"
	// tempID int = 1
)

//连接数据库
func ConnectMysql() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddrees, port, dbName, charset)
	Db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
		return Db, err
	}
	return Db, err
}

//查询信息
func QueryData(Db *sqlx.DB) []DataModelWorker {
	var modelworkers []DataModelWorker
	//sqlStr := "select worker_id, worker_name, IFNULL(remark,'') as remark, del_flag, IFNULL(create_by,'') as create_by, IFNULL(create_time,'') as create_time, IFNULL(update_by,'') as update_by, IFNULL(update_time,'') as update_time, IFNULL(gid,'') as gid from ModelWorkers"
	sqlStr := "select worker_id, worker_name, IFNULL(remark,'') as remark, del_flag, IFNULL(create_by,'') as create_by, IFNULL(create_time,'') as create_time, IFNULL(update_by,'') as update_by, IFNULL(update_time,'') as update_time, IFNULL(gid,'') as gid from worker_auditremark"
	err := Db.Select(&modelworkers, sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		// return
	}
	return modelworkers
}

func QueryFailureCount(Db *sqlx.DB) int {
	var failureCount int
	//sqlStr := "select count(*) as id from FailureInformation where status = 0" //成功失败记录表
	sqlStr := "select count(*) from worker_docking_record where upload_flag = 0" //成功失败记录表
	err := Db.Get(&failureCount, sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return 0
	}
	return failureCount
}

func QueryDataFailure(Db *sqlx.DB) []DataModelWorker {
	var modelworkers []DataModelWorker
	//sqlStr := "worker_id, worker_name, IFNULL(remark,'') as remark, del_flag, IFNULL(create_by,'') as create_by, IFNULL(create_time,'') as create_time, IFNULL(update_by,'') as update_by, IFNULL(update_time,'') as update_time, IFNULL(gid,'') as gid from ModelWorkers where audit_id in (select id from FailureInformation where status = 0)"
	sqlStr := "worker_id, worker_name, IFNULL(remark,'') as remark, del_flag, IFNULL(create_by,'') as create_by, IFNULL(create_time,'') as create_time, IFNULL(update_by,'') as update_by, IFNULL(update_time,'') as update_time, IFNULL(gid,'') as gid from worker_auditremark where audit_id in (select audit_id from worker_docking_record where upload_flag = 0)"
	err := Db.Select(&modelworkers, sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		// return
	}
	return modelworkers
}

//创建个人
func AddRecord(Db *sqlx.DB, args []string) {

	tempID := 0

	workerid := args[0]
	workername := args[1]
	remark := args[2]
	delflag := args[3]
	createby := args[4]
	createtime := args[5]
	updateby := args[6]
	updatetime := args[7]
	gid := args[8]

	result, err := Db.Exec("insert into ModelWorkers values(?,?,?,?,?,?,?,?,?,?)", tempID, workerid /*workerid*/, workername, /*workername*/
		remark /*remark*/, delflag /*delflag*/, createby /*createby*/, createtime /*createtime*/, updateby /*updateby*/, updatetime /*updatetime*/, gid /*gid*/)
	//result, err := Db.Exec("insert into worker_auditremark values(?,?,?,?,?,?,?,?,?,?)", tempID, workerid /*workerid*/, workername, /*workername*/
	//	remark /*remark*/, delflag /*delflag*/, createby /*createby*/, createtime /*createtime*/, updateby /*updateby*/, updatetime /*updatetime*/, gid /*gid*/)
	if err != nil {
		fmt.Printf("data insert failed, error:[%v]", err.Error())
		return
	}

	id, _ := result.LastInsertId()

	fmt.Printf("insert success, last id:[%d]\n", id)
}

//更新个人
func UpdateRecord(Db *sqlx.DB, args []string) {

	workerid := args[0]
	workername := args[1]
	remark := args[2]
	delflag := args[3]
	createby := args[4]
	createtime := args[5]
	updateby := args[6]
	updatetime := args[7]
	gid := args[8]

	var ModelWorker1 DataModelWorker

	//sqlStr := "select worker_id, worker_name, IFNULL(remark,'') as remark, del_flag, IFNULL(create_by,'') as create_by, IFNULL(create_time,'') as create_time, IFNULL(update_by,'') as update_by, IFNULL(update_time,'') as update_time, IFNULL(gid,'') as gid from ModelWorkers where worker_id = ?"
	sqlStr := "select worker_id, worker_name, IFNULL(remark,'') as remark, del_flag, IFNULL(create_by,'') as create_by, IFNULL(create_time,'') as create_time, IFNULL(update_by,'') as update_by, IFNULL(update_time,'') as update_time, IFNULL(gid,'') as gid from worker_auditremark where worker_id = ?"
	tempDb, _ := ConnectMysql()
	err := tempDb.Get(&ModelWorker1, sqlStr, workerid)
	tempDb.Close()

	if err != nil {
		fmt.Println("get data failed, error:[%v]", err.Error())
	}

	//
	// tempstring := []string{workerid, workername, remark, delflag, createby, createtime, updateby, updatetime, gid}
	// AddRecord(Db, tempstring)

	if ModelWorker1.Worker_name != workername {
		sqlStr := "update ModelWorkers set `worker_name` = '" + workername + "' where `worker_id` = '" + workerid + "'"
		//sqlStr := "update worker_auditremark set `worker_name` = '" + workername + "' where `worker_id` = '" + workerid + "'"
		result, err := Db.Exec(sqlStr)
		if err != nil {
			fmt.Printf("update worker_name failed, error:[%v]", err.Error())
			return
		}
		num, _ := result.RowsAffected()
		fmt.Printf("update success, affected rows:[%d]\n", num)
	}

	if ModelWorker1.Remark != remark {
		//sqlStr := "update ModelWorkers set `remark` = '" + remark + "' where `worker_id` = '" + workerid + "'"
		sqlStr := "update worker_auditremark set `remark` = '" + remark + "' where `worker_id` = '" + workerid + "'"
		result, err := Db.Exec(sqlStr)
		if err != nil {
			fmt.Printf("update remark failed, error:[%v]", err.Error())
			return
		}
		num, _ := result.RowsAffected()
		fmt.Printf("update success, affected rows:[%d]\n", num)

	}

	if ModelWorker1.Del_flag != delflag {
		//sqlStr := "update ModelWorkers set `del_flag` = '" + delflag + "' where `worker_id` = '" + workerid + "'"
		sqlStr := "update worker_auditremark set `remark` = '" + remark + "' where `worker_id` = '" + workerid + "'"
		result, err := Db.Exec(sqlStr)
		if err != nil {
			fmt.Printf("update status failed, error:[%v]", err.Error())
			return
		}
		num, _ := result.RowsAffected()
		fmt.Printf("update success, affected rows:[%d]\n", num)
	}

	if ModelWorker1.Create_by != createby {
		//sqlStr := "update ModelWorkers set `create_by` = '" + createby + "' where `worker_id` = '" + workerid + "'"
		sqlStr := "update worker_auditremark set `create_by` = '" + createby + "' where `worker_id` = '" + workerid + "'"
		result, err := Db.Exec(sqlStr)
		if err != nil {
			fmt.Printf("update create_by failed, error:[%v]", err.Error())
			return
		}
		num, _ := result.RowsAffected()
		fmt.Printf("update success, affected rows:[%d]\n", num)
	}

	if ModelWorker1.Create_time != createtime {
		//sqlStr := "update ModelWorkers set `create_time` = '" + createtime + "' where `worker_id` = '" + workerid + "'"
		sqlStr := "update worker_auditremark `create_time` = '" + createtime + "' where `worker_id` = '" + workerid + "'"
		result, err := Db.Exec(sqlStr)
		if err != nil {
			fmt.Printf("update create_time failed, error:[%v]", err.Error())
			return
		}
		num, _ := result.RowsAffected()
		fmt.Printf("update success, affected rows:[%d]\n", num)
	}

	if ModelWorker1.Update_by != updateby {
		//sqlStr := "update ModelWorkers set `update_by` = '" + updateby + "' where `worker_id` = '" + workerid + "'"
		sqlStr := "update worker_auditremark set `update_by` = '" + updateby + "' where `worker_id` = '" + workerid + "'"
		result, err := Db.Exec(sqlStr)
		if err != nil {
			fmt.Printf("update update_by failed, error:[%v]", err.Error())
			return
		}
		num, _ := result.RowsAffected()
		fmt.Printf("update success, affected rows:[%d]\n", num)
	}

	if ModelWorker1.Update_time != updatetime {
		//sqlStr := "update ModelWorkers set `update_time` = '" + updatetime + "' where `worker_id` = '" + workerid + "'"
		sqlStr := "update worker_auditremark set `update_time` = '" + updatetime + "' where `worker_id` = '" + workerid + "'"
		result, err := Db.Exec(sqlStr)
		if err != nil {
			fmt.Printf("update update_time failed, error:[%v]", err.Error())
			return
		}
		num, _ := result.RowsAffected()
		fmt.Printf("update success, affected rows:[%d]\n", num)
	}

	if ModelWorker1.Gid != gid {
		//sqlStr := "update ModelWorkers set `gid` = '" + gid + "' where `worker_id` = '" + workerid + "'"
		sqlStr := "update worker_auditremark set `gid` = '" + gid + "' where `worker_id` = '" + workerid + "'"
		result, err := Db.Exec(sqlStr)
		if err != nil {
			fmt.Printf("update gid failed, error:[%v]", err.Error())
			return
		}
		num, _ := result.RowsAffected()
		fmt.Printf("update success, affected rows:[%d]\n", num)
	}

	// num, _ := result.RowsAffected()
	// fmt.Printf("update success, affected rows:[%d]\n", num)
}
