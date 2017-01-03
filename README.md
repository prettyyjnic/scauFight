#scauFight 华农正方系统抢课用
## 现在支持抢中文课，A系列课程
### 使用方法
    1. 修改bin目录下的帐号密码，正方的入口地址
    2. 复制config.ini.sample 为 config.ini 并且修改账号密码
    3. 修改run.go 中 courses 数组为你想要抢课的数组，
    4. 运行 `go run run.go` ，根据提示输入验证码（由于session可能丢失，所以可能为不定次数）
    