package conf

import "time"

type Net struct {
	MaxOpenRequests int           `yaml:"max_open_requests"`
	DialTimeout     time.Duration `yaml:"dial_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	TLS             TLS           `yaml:"TLS"`
	SASL            SASL          `yaml:"SASL"`
	KeepAlive       time.Duration `yaml:"keep_alive"`
}

type TLS struct {
	Enable bool `yaml:"enable"`
}

type SASL struct {
	Enable    bool   `yaml:"enable"`
	Handshake bool   `yaml:"hand_shake"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
}
