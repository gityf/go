package dbadapter

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logger "github.com/xlog4go"
)

type DBPool struct {
	pool     chan *sql.DB
	size     int
	conn_str string
	timeout  time.Duration
}

var (
	Err_Pool_Empty = errors.New("poolempty")
	DbGradeTimeout = 800
)

func NewPool(connect_str string, max_size, min_size int) (pool *sql.DB, err error) {
	if pool, err = sql.Open("mysql", connect_str); err == nil {
		pool.SetMaxIdleConns(min_size)
		pool.SetMaxOpenConns(max_size)
	}
	return
}

func NewDBPool(connect_str string, pool_size int) (pool *DBPool, err error) {
	pool = new(DBPool)
	pool.pool = make(chan *sql.DB, pool_size)
	pool.size = pool_size
	pool.conn_str = connect_str
	for i := 0; i < pool.size; i++ {
		pool.pool <- nil
	}
	pool.timeout = time.Duration(DbGradeTimeout)

	return
}

func (pool *DBPool) GetConn() (conn *sql.DB, err error) {
	if pool.pool == nil {
		return nil, errors.New("pool is not initialized")
	}

	select {
	case conn = <-pool.pool:
		if conn == nil {
			conn, err = sql.Open("mysql", pool.conn_str)
			if err == nil {
				conn.SetMaxIdleConns(1)
				conn.SetMaxOpenConns(1)
				//go-1.4不支持
				//conn.SetConnMaxLifetime(time.Minute * 3) //根据DBA建议,设置3分钟
				logger.Debug("DBPool: get nil conn from pool, open one: %v", conn)
			} else {
				logger.Error("DBPool: errorcode while open new conn")
			}
		}
	case <-time.After(pool.timeout):
		logger.Error("DBPool: get conn from pool time out")
		err = Err_Pool_Empty
	}
	return
}

func (pool *DBPool) ReleaseConn(conn *sql.DB) (err error) {
	if pool.pool == nil {
		if conn != nil {
			logger.Error("DBPool: close conn due to pool.pool is nil: %v", conn)
			doCloseConn(conn)
		}
		return
	}
	select {
	case pool.pool <- conn:
	default:
		if conn != nil {
			logger.Error("DBPool: close conn due to pool.pool is full in ReleaseConn: %v", conn)
			doCloseConn(conn)
		}
	}
	return
}

func (pool *DBPool) CloseConn(conn *sql.DB) (err error) {
	if pool.pool == nil {
		if conn != nil {
			logger.Debug("DBPool: close conn due to pool.pool is nil: %v", conn)
			doCloseConn(conn)
		}
		return
	}
	select {
	case pool.pool <- nil:
		if conn != nil {
			logger.Debug("DBPool: close conn nornally: %v", conn)
			doCloseConn(conn)
		}
	default:
		if conn != nil {
			logger.Error("DBPool: close conn due to pool.pool is full in CloseConn: %v", conn)
			doCloseConn(conn)
		}
	}
	return
}

func doCloseConn(conn *sql.DB) (err error) {
	if conn != nil {
		err = conn.Close()
		if err != nil {
			logger.Error("DBPool: errorcode while close conn: %v", err)
		}
	}
	return
}
