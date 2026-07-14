package config

type SiteConfig struct {
	Workers         []*WorkerConfig `json:"workers"`
	Name            string          `json:"name"`
	URL             string          `json:"url"`
	ProxyURL        string          `json:"proxy_url"`
	PagesBack       int             `json:"pages_back"`
	BlockedCooldown int             `json:"blocked_cooldown"`
	Enabled         bool            `json:"enabled"`
}

type WorkerConfig struct {
	Name     string `json:"name"`
	NumberOf int    `json:"number_of"`
	// milliseconds
	Cooldown int `json:"cooldown"`
}

type BrowserSettings struct {
	EngineSettings map[string]any `json:"engine_settings"`
	// milliseconds
	Timeout            int  `json:"timeout"`
	RandomizeUserAgent bool `json:"randomize_user_agent"`
}

// all in miliseconds
type ScraperSettings struct {
	ValidateAndUpdatePagesInterval int  `json:"validate_pages_interval"`
	ValidateAndUpdatePagesCooldown int  `json:"validate_pages_cooldown"`
	FailedPagesCooldown            int  `json:"failed_pages_cooldown"`
	FailedPagesWaitTime            int  `json:"failed_pages_wait_time"`
	Debug                          bool `json:"debug"`
	CliLogging                     bool `json:"cli_logging"`
}

type ApiConfig struct {
	Port               string     `json:"port"`
	ServerHost         string     `json:"server_host"`
	JWT                *JwtConfig `json:"jwt"`
	Debug              bool       `json:"debug"`
	PublicHost         string     `json:"public_host"`
	PublicFrontendHost string     `json:"public_frontend_host"`
}

type JwtConfig struct {
	Secret         string `json:"secret"`
	ExpirationTime int    `json:"expiration_time"`
}

type MailerConfig struct {
	SMTPHost    string `json:"smtp_host"`
	SMTPPort    int    `json:"smtp_port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	SenderEmail string `json:"sender_email"`
	SenderName  string `json:"sender_name"`
	UseTLS      bool   `json:"use_tls"`
	Debug       bool   `json:"debug"`
}

type NotifierConfig struct {
	WorkerCount int `json:"worker_count"`
}

type DBConfig struct {
	User       string `json:"user"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Debug      bool   `json:"debug"`
	CliLogging bool   `json:"cli_logging"`
}
type RabbitMQConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Vhost    string `json:"vhost"`
}

type PaymentMethodConfig struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type RedisConfig struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Password          string `json:"password"`
	DB                int    `json:"db"`
	DefaultTTLSeconds int    `json:"default_ttl_seconds"`
}
