package xlog

type Config struct {
	Name         string `yaml:"name" json:"name"`                   // 日志服务名  写入文件或者发到日志收集服务时使用
	Stdout       bool   `yaml:"stdout" json:"stdout"`               // 是否输出
	Dir          string `yaml:"dir" json:"dir"`                     // 日志根目录
	MaxAge       int    `yaml:"max_age" json:"max_age"`             // 日志保留天数
	RotationHour int    `yaml:"rotation_hour" json:"rotation_hour"` // 分割小时数，默认为1小时
}
