package shared

type Config struct {
	Superuser struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"superuser"`
	StartTimer int `yaml:"start_timer" mapstructure:"start_timer"`
}
