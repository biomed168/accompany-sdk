# 跨端sdk

核心业务逻辑使用golang进行开发，编译后嵌入各个端

## 思路

### pc & macos

electron， 如果在electron环境，可使用sqlite3

### web

go 编译 wasm, 在web使用， 数据库采用sqlite3编译的wasm与indexdb嵌入


### ios

go 编译 ios， 使用flutter channel进行互相通讯

### android

go 编译 aar ， 使用flutter channel进行原生与通讯
