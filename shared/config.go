package shared

type Config struct {
	Superuser struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"superuser"`
	StartTimer int `yaml:"start_timer" mapstructure:"start_timer"`
	DNS        struct {
		PrimaryPort  int `yaml:"primary_port" mapstructure:"primary_port"`
		FallbackPort int `yaml:"fallback_port" mapstructure:"fallback_port"`
	} `yaml:"dns"`
	DomainManager struct {
		Port int `yaml:"port"`
	} `yaml:"domain_manager"`
	Mail struct {
		PrimaryPort  int `yaml:"primary_port" mapstructure:"primary_port"`
		FallbackPort int `yaml:"fallback_port" mapstructure:"fallback_port"`
		AdminPort    int `yaml:"admin_port" mapstructure:"admin_port"`
		WebmailPort  int `yaml:"webmail_port" mapstructure:"webmail_port"`
	} `yaml:"mail"`
}
