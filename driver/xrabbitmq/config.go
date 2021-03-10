package xrabbitmq

// 配置文件
type Config struct {
	Name     string `json:"name" yaml:"name"` // consumer tag需要
	UserName string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
}
