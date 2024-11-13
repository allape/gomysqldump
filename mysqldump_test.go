package gomysqldump

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"time"
)

func TestMySQLDump(t *testing.T) {
	// docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=Root_123456 -p 3306:3306 mysql
	repo, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: "root:Root_123456@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local",
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.CreateTemp(cwd, "dumped_*.sql")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()

	err = MySQLDump(file, repo, &Config{Timeout: 10 * time.Minute, Logger: log.Default()})
	if err != nil {
		_ = os.Remove(file.Name())
		t.Fatal(err)
	}

	t.Log("Dumped to", file.Name())
}
