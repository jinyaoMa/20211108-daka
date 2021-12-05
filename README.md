# 前后端开发环境

- Node [https://nodejs.org/dist/v16.13.0/node-v16.13.0-x64.msi](https://nodejs.org/dist/v16.13.0/node-v16.13.0-x64.msi)
- Go [https://golang.org/dl/go1.17.3.windows-amd64.msi](https://golang.org/dl/go1.17.3.windows-amd64.msi)

# 测试

- 运行`curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin`
- 运行`go install github.com/swaggo/swag/cmd/swag@latest`
- 运行`npm run setup:all`
- 运行`npm run serve:front`
- 运行`npm run serve:back`
- 浏览器打开：`http://localhost:8080`
- Swagger打开：`http://localhost:8081/swagger/index.html`

# 数据库

- `/database/init.go`
  - `MAIN` - 登录 user 表
  - `STORES` - 所有店铺链接
  - `OFFICE_EXCEPT_INDEX` - Office 用户不能用的链接，对应`STORES`下标

# 网站端口修改

- `/client/src/main.js`
  - `axios.defaults.baseURL = "http://localhost:8081"` - 替换 8081
- `/server/init.go`
  - `Addr: ":8081",` - 替换 8081

# 部署

- 运行`npm run build:all`
- 确保`/client`文件夹和`main.exe`执行文件在同一目录下，运行`main.exe`
