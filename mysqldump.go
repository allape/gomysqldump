package gomysqldump

import (
	"context"
	"fmt"
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

func MySQLDump(file *os.File, db *gorm.DB, cfg *Config) error {
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

	config := db.Config.Dialector.(*mysql.Dialector).Config

	timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	host := "127.0.0.1"
	port := "3306"
	addr := config.DSNConfig.Addr
	if addr != "" {
		host, port, err = net.SplitHostPort(addr)
		if err != nil {
			return err
		}
	}

	args := []string{
		"--host",
		host,
		"--port",
		port,
		"--user",
		config.DSNConfig.User,
		fmt.Sprintf("--password=%s", config.DSNConfig.Passwd),
		"--databases",
		config.DSNConfig.DBName,
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
		return err
	}

	return nil
}
