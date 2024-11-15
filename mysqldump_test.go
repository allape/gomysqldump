package gomysqldump

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

// Start a MySQL container with the following command before running the tests:
// docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=Root_123456 -p 3306:3306 mysql

const TestData = "testdata"

func CreateTmpFile() (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	folder := path.Join(cwd, TestData)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.Mkdir(folder, 0755); err != nil {
			return nil, err
		}
	}

	file, err := os.CreateTemp(folder, fmt.Sprintf("dumped_%d_*.sql", time.Now().Unix()))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func TestMySQLDumpFromGORM(t *testing.T) {
	repo, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: "root:Root_123456@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local",
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	file, err := CreateTmpFile()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = MySQLDumpFromGORM(file, repo, &Config{Timeout: 10 * time.Minute, Logger: log.Default()})
	if err != nil {
		_ = os.Remove(file.Name())
		t.Fatal(err)
	}

	t.Log("Dumped to", file.Name())
}

func TestMySQLDumpFromDSNString(t *testing.T) {
	file, err := CreateTmpFile()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = MySQLDumpFromDSNString(
		file,
		"root:Root_123456@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local",
		&Config{
			Timeout: 10 * time.Minute,
			Logger:  log.Default(),
		},
	)
	if err != nil {
		_ = os.Remove(file.Name())
		t.Fatal(err)
	}

	t.Log("Dumped to", file.Name())
}
