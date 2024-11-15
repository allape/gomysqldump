# mysqldump in Golang

## Dependencies

### Alpine

```shell
apk add mariadb-client
```

### macOS

```shell
brew install mysql-client

echo "export PATH=\"\$PATH:/usr/local/opt/mysql-client/bin\"" >> ~/.zshrc

source ~/.zshrc
```

## Example

see [mysqldump_test.go](mysqldump_test.go)
