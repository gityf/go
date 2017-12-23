package dbadapter

import (
	"database/sql"
	"reflect"
	"time"
	_ "github.com/go-sql-driver/mysql"
	logger "github.com/xlog4go"
	"common/statistic"
	"strings"
	"io/ioutil"
	"regexp"
	"encoding/json"
)

const (
	MYSQL_MAX_RETRY_TIMES     = 200
	MYSQL_RECONN_INTERVAL_MS  = 10
)

type MysqlAccess struct {
	mysqlClient map[string]*DBPool
	caller      map[string]string
}

func NewMysqlAccess() *MysqlAccess {
	var ma MysqlAccess
	ma.mysqlClient = make(map[string]*DBPool)
	return &ma
}

type DatabaseMysqlConfig struct {
	Database    string `json:"database"`
	Dsn         string `json:"dsn"`
	DbDriver    string `json:"dbdriver"`
	MaxOpenConn int    `json:"maxopenconn"`
	MaxIdleConn int    `json:"maxidleconn"`
	MaxLifeTime int    `json:"maxlifetime"`
}

type DatabaseConfig struct {
	Mysql []DatabaseMysqlConfig `json:"mysql"`
}

func IsDuplicatedMysqlError(errStr string) bool {
	return strings.Contains(errStr, "Error 1062: Duplicate")
}

func LoadConfigJson(path string, obj interface{}) (err error) {
	var bytes []byte
	bytes, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	re1 := regexp.MustCompile("^#.*$")
	lines := strings.Split(string(bytes), "\n")
	var loc []int
	var cnt string
	for _, line := range lines {
		loc = re1.FindStringIndex(line)
		if len(loc) == 0 {
			cnt += line
		}
	}

	err = json.Unmarshal([]byte(cnt), obj)
	return
}

func InitMysqlAccess(file string) (mysql *MysqlAccess, err error) {

	var dbcfg DatabaseConfig
	var dbPool *DBPool
	err = LoadConfigJson(file, &dbcfg)
	if err != nil {
		logger.Error("database Unmarshal config file failed. err=%s", err)
		return
	}

	mysql = &MysqlAccess{}
	mysql.mysqlClient = make(map[string]*DBPool)
	mysql.caller = make(map[string]string)

	for _, mysqlCfg := range dbcfg.Mysql {

		if dbPool, err = NewDBPool(mysqlCfg.Dsn, mysqlCfg.MaxOpenConn); err != nil {
			logger.Error("database NewDBPool(%v, %v) failed. err=%v", mysqlCfg.Dsn, mysqlCfg.MaxOpenConn, err)
			return
		}
		mysql.mysqlClient[mysqlCfg.Database] = dbPool
	}

	return
}

func (this *MysqlAccess) getPool(sid string) (dbpool *DBPool, exist bool) {
	var database string
	if database, exist = this.caller[sid]; exist {
		dbpool, exist = this.mysqlClient[database]
	}
	return
}

func (this *MysqlAccess) RegisterCaller(sid string, database string, limit int) {
	this.caller[sid] = database
	return
}

func (dbAccess *MysqlAccess) Insert(sid string, dbpool *DBPool) {
	dbAccess.mysqlClient[sid] = dbpool
}

func (dbAccess *MysqlAccess) getConn(sid string) (conn *sql.DB, connStr string, err error) {
	//	dbPool, exist := dbAccess.mysqlClient[dbAccess.GenId(sid, gid)]
	dbPool, exist := dbAccess.getPool(sid)
	if !exist {
		conn = nil
		err = Err_Pool_Empty
		return
	}
	connStr = dbPool.conn_str
	for ii := 0; ii < MYSQL_MAX_RETRY_TIMES; ii++ {
		conn, err = dbPool.GetConn()
		if err == nil {
			break
		}
		if err == Err_Pool_Empty {
			time.Sleep(time.Duration(MYSQL_RECONN_INTERVAL_MS) * time.Millisecond)
			continue
		}
	}
	return
}

func (dbAccess *MysqlAccess) releaseConn(sid string, conn *sql.DB) (err error) {
	//	dbPool, exist := dbAccess.mysqlClient[dbAccess.GenId(sid, gid)]
	dbPool, exist := dbAccess.getPool(sid)
	if !exist {
		return Err_Pool_Empty
	}
	err = dbPool.ReleaseConn(conn)
	return
}

func (dbAccess *MysqlAccess) closeDBConn(sid string, conn *sql.DB) (err error) {
	//	dbPool, exist := dbAccess.mysqlClient[dbAccess.GenId(sid, gid)]
	dbPool, exist := dbAccess.getPool(sid)
	if !exist {
		return Err_Pool_Empty
	}
	err = dbPool.CloseConn(conn)
	return
}

func (dbAccess *MysqlAccess) Exec(sid string, execSql string) (num int64, err error) {
	var connErr error
	var conn *sql.DB
	defer func() {
		if r := recover(); r != nil {
			statistic.IncPanicCount()
			logger.Error("Exec recover panic: %v", r)
			dbAccess.closeDBConn(sid, conn)
			return
		}
		if connErr != nil {
			return
		}
		if err == nil {
			dbAccess.releaseConn(sid, conn)
			return
		}
		//不是主键冲突错误,关闭连接
		if IsDuplicatedMysqlError(err.Error()) {
			dbAccess.releaseConn(sid, conn)
		} else {
			dbAccess.closeDBConn(sid, conn)
		}
	}()
	conn, _, connErr = dbAccess.getConn(sid)
	if connErr != nil {
		err = connErr
		logger.Error("getconn err:%v", err)
		return 0, err
	}
	var result sql.Result
	result, err = conn.Exec(execSql)

	if err != nil {
		//		statistic.CreditInfoPointer.Load().(*statistic.CreditStatistic).IncReqMysqlExecErr()
		logger.Debug("execErr sql:[%v], err:%v", execSql, err)
		return
	}
	num, err = result.RowsAffected()
	return
}

// Params stores the Params
type Params map[string]interface{}

// ParamsList stores paramslist
type ParamsList []interface{}

// query data to [][]interface
func (dbAccess *MysqlAccess) Query(sid string, querySql string, container *[]ParamsList, needCols ...string) (num int64, err error) {
	var connErr error
	var conn *sql.DB
	defer func() {
		if r := recover(); r != nil {
			statistic.IncPanicCount()
			logger.Error("recover panic: %v", r)
			dbAccess.closeDBConn(sid, conn)
			return
		}
		if connErr != nil {
			return
		}
		if err == nil {
			dbAccess.releaseConn(sid, conn)
			return
		}
		dbAccess.closeDBConn(sid, conn)
	}()
	conn, _, connErr = dbAccess.getConn(sid)
	if connErr != nil {
		err = connErr
		logger.Error("getconn err:%v", err)
		return
	}
	num, err = dbAccess.query(conn, querySql, container, needCols)
	return
}

// query data to []map[string]interface{}
func (dbAccess *MysqlAccess) QueryMap(sid string, querySql string, container *[]Params, needCols ...string) (num int64, err error) {
	var connErr error
	var conn *sql.DB
	defer func() {
		if r := recover(); r != nil {
			statistic.IncPanicCount()
			logger.Error("Query recover panic: %v", r)
			dbAccess.closeDBConn(sid, conn)
			return
		}
		if connErr != nil {
			return
		}
		if err == nil {
			dbAccess.releaseConn(sid, conn)
			return
		}
		dbAccess.closeDBConn(sid, conn)
	}()
	conn, _, connErr = dbAccess.getConn(sid)
	if connErr != nil {
		err = connErr
		logger.Error("getconn err:%v", err)
		return
	}
	num, err = dbAccess.query(conn, querySql, container, needCols)
	return
}

func (dbAccess *MysqlAccess) query(conn *sql.DB, querySql string, container interface{}, needCols []string) (int64, error) {
	//tBegin := common.NowInNs()
	var err error
	var rows *sql.Rows
	rows, err = conn.Query(querySql)
	if err != nil {
		logger.Error("queryErr sql:[%v], err:%v", querySql, err)
		return 0, err
	}

	defer func() {

		rows.Close()
	}()
	var (
		maps  []Params
		lists []ParamsList
		list  ParamsList
	)

	typ := 0
	switch container.(type) {
	case *[]Params:
		typ = 1
	case *[]ParamsList:
		typ = 2
	case *ParamsList:
		typ = 3
	default:
		return 0, nil
	}
	var (
		refs   []interface{}
		cnt    int64
		cols   []string
		indexs []int
	)

	for rows.Next() {
		if cnt == 0 {
			columns, err := rows.Columns()
			if err != nil {
				return 0, err
			}
			/*for _, xx := range columns {
				logger.Warn("columns: [%v]", xx)
			}*/
			if len(needCols) > 0 {
				indexs = make([]int, 0, len(needCols))
			} else {
				indexs = make([]int, 0, len(columns))
			}

			cols = columns
			refs = make([]interface{}, len(cols))
			for i := range refs {
				var ref sql.NullString
				refs[i] = &ref

				if len(needCols) > 0 {
					for _, c := range needCols {
						if c == cols[i] {
							indexs = append(indexs, i)
						}
					}
				} else {
					indexs = append(indexs, i)
				}
			}
		}

		if err := rows.Scan(refs...); err != nil {
			return 0, err
		}

		switch typ {
		case 1:
			params := make(Params, len(cols))
			for _, i := range indexs {
				ref := refs[i]
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					params[cols[i]] = value.String
				} else {
					params[cols[i]] = nil
				}
			}
			maps = append(maps, params)
		case 2:
			params := make(ParamsList, 0, len(cols))
			for _, i := range indexs {
				ref := refs[i]
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					params = append(params, value.String)
				} else {
					params = append(params, nil)
				}
			}
			lists = append(lists, params)
		case 3:
			for _, i := range indexs {
				ref := refs[i]
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					list = append(list, value.String)
				} else {
					list = append(list, nil)
				}
			}
		}

		cnt++
	}

	switch v := container.(type) {
	case *[]Params:
		*v = maps
	case *[]ParamsList:
		*v = lists
	case *ParamsList:
		*v = list
	}

	return cnt, nil
}
