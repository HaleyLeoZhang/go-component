package xneo4j

type Config struct {
	Name string `yaml:"name" json:"name"` // 日志服务名  写入文件或者发到日志收集服务时使用
	DSN  string `yaml:"dsn" json:"dsn"`   // 链接地址 如 neo4j://localhost:7687
	// 账号信息
	UserName string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}
