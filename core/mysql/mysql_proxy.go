package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"../common/zlog"
	"../conf"
	_ "github.com/go-sql-driver/mysql"
)

// _mysqlDb .
var _mysqlDb *sql.DB

// InitMysqlDB .
func InitMysqlDB() error {
	// 账号库初始化
	if conf.Conf.MysqlAccounts != "" {
		dbAccounts, err := sql.Open("mysql", conf.Conf.MysqlAccounts)
		if err != nil {
			zlog.Error("InitMysqlDB mysql err", zlog.String("Err", err.Error()))
			return errors.New("Connect envdevice db error")
		}
		dbAccounts.SetMaxOpenConns(500)
		dbAccounts.SetMaxIdleConns(500)
		if err = dbAccounts.Ping(); err != nil {
			zlog.Error("InitMysqlDB log err %v", zlog.String("Err", err.Error()))
			err = dbAccounts.Close()
			if err != nil {
				zlog.Error("Mysql db close err:%v", zlog.String("Err", err.Error()))
			}
			return errors.New("Connect envdevice db ping error")
		}
		_mysqlDb = dbAccounts
		zlog.Info("db_server connect envdevice db successful!")
	}
	return nil
}
func CloseMysqlDB() {
	if _mysqlDb != nil {
		err := _mysqlDb.Close()
		if err != nil {
			zlog.Error("Mysql db close err:%v", zlog.String("Err", err.Error()))
		}
	}
}

// 创建设备表
func CreateDeviceTable() error {
	table := "envdevice"
	if _mysqlDb != nil {
		var s string
		sh := fmt.Sprintf("SHOW TABLES LIKE '%s'", table)
		if err := _mysqlDb.QueryRow(sh).Scan(&s); err == nil { // 表已经存在
			return errors.New("CreateEventTable is Exist")
		} else if err == sql.ErrNoRows { // 表不存在
			_, err := _mysqlDb.Exec("CREATE TABLE " + table + " (" +
				"`event_id`  int(10) NOT NULL AUTO_INCREMENT ," +
				"`device_id`  varchar(300) ," +
				"`event_time`  varchar(300) ," +
				"PRIMARY KEY (`event_id`))" +
				"ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci " +
				"ROW_FORMAT=DYNAMIC;")
			if err != nil {
				return err
			} else {
				zlog.Info("CreateEventTable", zlog.String("table", table))
				return nil
			}
		} else {
			return errors.New(fmt.Sprintf("CreateEventTable err %v", err))
		}
	}
	return errors.New(fmt.Sprintf("Db connect failed."))
}

// 插入设备
func InsertDevice(id string, ts string) error {
	if _mysqlDb != nil {
		table := "envdevice"
		if _, err := _mysqlDb.Exec("INSERT INTO "+table+" (device_id, event_time"+
			") VALUES (?, ?)",
			id, ts); err == nil {
			return nil
		} else {
			//如果没有这个表 那么就创建下~
			errEvent := CreateDeviceTable()
			if errEvent == nil {
				if _, err := _mysqlDb.Exec("INSERT INTO "+table+" (device_id, event_time"+
					") VALUES (?, ?)",
					id, ts); err == nil {
					return nil
				} else {
					return errors.New(fmt.Sprintf("Insert device table %v", table))
				}
			} else {
				return errors.New(fmt.Sprintf("Create device table err %v", err))
			}
		}
	} else {
		return errors.New("InsertUerEvent dbLog is nil.")
	}
}

// 创建事件表
func CreateDeviceEventTable(table string) error {
	if _mysqlDb != nil {
		var s string
		sh := fmt.Sprintf("SHOW TABLES LIKE '%s'", table)
		if err := _mysqlDb.QueryRow(sh).Scan(&s); err == nil { // 表已经存在
			return errors.New("CreateEventTable is Exist")
		} else if err == sql.ErrNoRows { // 表不存在
			_, err := _mysqlDb.Exec("CREATE TABLE " + table + " (" +
				"`event_id`  int(10) NOT NULL AUTO_INCREMENT ," +
				"`device_id`  varchar(300) ," +
				"`event_time`  varchar(300) ," +
				"`event_data`  varchar(300) ," +
				"`event_seq`  int(10) ," +
				"PRIMARY KEY (`event_id`))" +
				"ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci " +
				"ROW_FORMAT=DYNAMIC;")
			if err != nil {
				return err
			} else {
				zlog.Info("CreateEventTable", zlog.String("table", table))
				return nil
			}
		} else {
			return errors.New(fmt.Sprintf("CreateEventTable err %v", err))
		}
	}
	return errors.New(fmt.Sprintf("Db connect failed."))
}

// 插入设备事件
func InsertDeviceEvent(id string, sequence int, datas string, ts string) error {
	if _mysqlDb != nil {
		tb := fmt.Sprintf("device_%v", id)
		if _, err := _mysqlDb.Exec("INSERT INTO "+tb+" (device_id, event_time, event_data, event_seq"+
			") VALUES (?, ?, ?, ?)",
			id, ts, datas, sequence); err == nil {
			return nil
		} else {
			//如果没有这个表 那么就创建下~
			err := CreateDeviceEventTable(tb)
			if err == nil {
				if _, err := _mysqlDb.Exec("INSERT INTO "+tb+" (device_id, event_time, event_data, event_seq"+
					") VALUES (?, ?, ?, ?)",
					id, ts, datas, sequence); err == nil {
					return nil
				} else {
					return errors.New(fmt.Sprintf("Insert device table %v", err.Error()))
				}
			} else {
				return err
			}
		}
	} else {
		return errors.New("InsertUerEvent dbLog is nil.")
	}
}

// 查询历史记录
func FindDeviceEvent(id string) [][]string {
	if _mysqlDb != nil {
		tb := fmt.Sprintf("device_%v", id)
		rows, err := _mysqlDb.Query("SELECT device_id, event_time,event_data " +
			"FROM " + tb + " ORDER BY event_seq DESC limit 10")
		if err != nil {
			zlog.Info("Query db", zlog.String("Err", err.Error()), zlog.String("table", tb))
			return nil
		}
		queryData := [][]string{}
		for rows.Next() {
			var device_id, event_time, event_data string
			err = rows.Scan(&device_id, &event_time, &event_data)
			if err == nil {
				queryData = append(queryData, []string{device_id, event_time, event_data})
				if len(queryData) >= 10 {
					break
				}
			} else {
				zlog.Info("Scan  query result", zlog.String("Err", err.Error()))
			}
		}
		return queryData
	}
	return nil
}

//查询设备
func FindDevice() []string {
	if _mysqlDb != nil {
		rows, err := _mysqlDb.Query("SHOW TABLES")
		if err != nil {
			zlog.Info("Query db", zlog.String("Err", err.Error()))
			return nil
		}
		queryData := []string{}
		for rows.Next() {
			var tableName string
			err = rows.Scan(&tableName)
			if err == nil {
				queryData = append(queryData, tableName)
			} else {
				zlog.Info("Scan  query result", zlog.String("Err", err.Error()))
			}
		}
		return queryData
	}
	return nil
}
