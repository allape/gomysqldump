package gomysqldump

import (
	"context"
	"fmt"
	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	Timeout time.Duration
	Logger  *log.Logger
	OnArgs  func(args []string) []string
}

func MySQLDump(writer io.Writer, dsn *mysqldriver.Config, cfg *Config) (int64, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Minute
	}

	if cfg.Logger == nil {
		cfg.Logger = log.New(io.Discard, "[backup]", log.LstdFlags)
	}

	var err error

	timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	host := "127.0.0.1"
	port := "3306"
	addr := dsn.Addr
	if addr != "" {
		host, port, err = net.SplitHostPort(addr)
		if err != nil {
			return 0, err
		}
	}

	file, err := os.CreateTemp(os.TempDir(), "dump-*.sql")
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}()

	args := []string{
		"--host",
		host,
		"--port",
		port,
		"--user",
		dsn.User,
		fmt.Sprintf("--password=%s", dsn.Passwd),
		"--databases",
		dsn.DBName,
		"--result-file",
		file.Name(),
	}

	if cfg.OnArgs != nil {
		args = cfg.OnArgs(args)
	}

	// apk add mariadb-client
	// brew install mysql-client, export PATH="$PATH:/usr/local/opt/mysql-client/bin"
	cmd := exec.CommandContext(
		timeoutCtx,
		"mysqldump",
		args...,
	)

	output, err := cmd.CombinedOutput()

	cfg.Logger.Println(string(output))

	if err != nil {
		return 0, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return io.Copy(writer, file)
}

func MySQLDumpFromDSNString(writer io.Writer, dsnStr string, cfg *Config) (int64, error) {
	dsn, err := mysqldriver.ParseDSN(dsnStr)
	if err != nil {
		return 0, err
	}
	return MySQLDump(writer, dsn, cfg)
}

func MySQLDumpFromDialector(writer io.Writer, dialector *mysql.Dialector, cfg *Config) (int64, error) {
	return MySQLDump(writer, dialector.DSNConfig, cfg)
}

func MySQLDumpFromGORM(writer io.Writer, db *gorm.DB, cfg *Config) (int64, error) {
	return MySQLDumpFromDialector(writer, db.Config.Dialector.(*mysql.Dialector), cfg)
}
