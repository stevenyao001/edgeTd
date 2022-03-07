package tdengine

import (
	"fmt"
	"github.com/stevenyao001/edgeCommon/logger"
	"github.com/stevenyao001/edgeCommon/tdengine"
	"time"
)

type UserTd struct {
	*tdengine.TdEngine
	MegnetStatus bool `json:"megnet_status"`
	Ia           int  `json:"ia"`
	Ep           int  `json:"ep"`
}

type UserTdGK struct {
	*tdengine.TdEngine
	Num       int `json:"num"`
	Status    int `json:"status"`
	InstantEp int `json:"instant_ep"`
	MegnetStatus bool `json:"megnet_status"`
	Ia           int  `json:"ia"`
	Ep           int  `json:"ep"`
}

type Input struct {
	Ts         int64       `json:"ts"`
	Properties *Properties `json:"properties"`
}

type Properties struct {
	MegnetStatus bool `json:"megnet_status"`
	Num          int  `json:"num"`
	Status       int  `json:"status"`
	Ia           int  `json:"ia"`
	Ep           int  `json:"ep"`
	InstantEp    int  `json:"instant_ep"`
}

func NewUserTd() *UserTd {
	tde := &UserTd{
		TdEngine: &tdengine.TdEngine{
			Db:        nil,
			InsName:   "rootcloud",
			DbName:    "test",
			TableName: "megnet_001",
		},
	}

	tde.Conn()

	return tde
}

/*func (u *UserTd) Insert(uid int, name string) (ts int64, err error) {
	ts = time.Now().UnixMilli()
	sqls := fmt.Sprintf("insert into `%s` (uid,name,ts) values (%d,'%s',%d)", u.TableName, uid, name, ts)
	_, err = u.Db.Exec(sqls)
	if err != nil {
		logger.ErrorLog("UserTde-Insert", "执行sql报错", "", err)
		return
	}

	return
}

func (u *UserTd) Find() (res []UserTd, err error) {
	sqls := fmt.Sprintf("select uid,name,ts from %s limit 10", u.TableName)
	rows, err := u.Db.Query(sqls)
	if err != nil {
		logger.ErrorLog("UserTde-Find", "执行sql报错", "", err)
		return
	}
	defer rows.Close()

	res = make([]UserTd, 0, 0)
	for rows.Next() {
		tmp := UserTd{}
		err = rows.Scan(&tmp.Uid, &tmp.Name, &tmp.Ts2)
		if err != nil {
			logger.ErrorLog("UserTde-Find", "赋值报错", "", err)
			return
		}
		res = append(res, tmp)
	}

	return
}*/

func (u *UserTd) Find() error {
	sqls := fmt.Sprintf("select ep, ia, megnet_status from megnet_001 limit 10")
	rows, err := u.Db.Query(sqls)
	if err != nil {
		logger.ErrorLog("UserTde-Find", "执行sql报错", "", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		tmp := UserTd{}
		err = rows.Scan(&tmp.Ep, &tmp.Ia, &tmp.MegnetStatus)
		if err != nil {
			logger.ErrorLog("UserTde-Find", "赋值报错", "", err)
			return err
		}
		fmt.Println("megnet_001---:", tmp)
	}

	sqlsGK := fmt.Sprintf("select instant_ep, num, status, ep, ia, megnet_status from megnet_gk_001 limit 10")
	rowsGK, err := u.Db.Query(sqlsGK)
	if err != nil {
		logger.ErrorLog("UserTde-Find", "执行sql报错", "", err)
		return err
	}
	defer rowsGK.Close()

	for rowsGK.Next() {
		tmp := UserTdGK{}
		err = rowsGK.Scan(&tmp.InstantEp, &tmp.Num, &tmp.Status,&tmp.Ep, &tmp.Ia, &tmp.MegnetStatus)
		if err != nil {
			logger.ErrorLog("UserTde-Find", "执行sql报错", "", err)
			return err
		}
		fmt.Println("megnet_gk_001---:", tmp)
	}

	return nil
}

func (u *UserTd) Insert(data *Input) error {
	ts := time.Now().UnixMilli()
	sqls := fmt.Sprintf("insert into megnet_001 (ts, ep, ia, megnet_status) values (%d, %d, %d, %t)", ts, data.Properties.Ep, data.Properties.Ia, data.Properties.MegnetStatus)
	_, err := u.Db.Exec(sqls)
	if err != nil {
		logger.ErrorLog("UserTde-Insert", "执行sql报错", "", err)
		return err
	}

	tsGK := time.Now().UnixMilli()
	sqlsGK := fmt.Sprintf("insert into megnet_gk_001 (ts, megnet_id, instant_ep, num, status, ep, ia, megnet_status) values (%d, %d, %d, %d, %d, %d, %d, %t)", tsGK, ts, data.Properties.InstantEp, data.Properties.Num, data.Properties.Status, data.Properties.Ep, data.Properties.Ia, data.Properties.MegnetStatus)
	_, err = u.Db.Exec(sqlsGK)
	if err != nil {
		logger.ErrorLog("UserTde-Insert", "执行sql报错", "", err)
		return err
	}

	return nil
}
