package sibylConfig

type SibylSystemConfig struct {
	TokenSize                 int64    `json:"toke_size"`
	Owners                    []int64  `json:"owners"`
	MaxPanic                  int64    `json:"max_panic"`
	MaxCacheTime              int64    `json:"max_cache_time"`
	DbUrl                     string   `json:"db_url"`
	DbName                    string   `json:"db_name"`
	UseSqlite                 bool     `json:"use_sqlite"`
	Port                      string   `json:"port"`
	BotToken                  string   `json:"bot_token"`
	BotAPIUrl                 string   `json:"api_url"`
	Debug                     bool     `json:"debug"`
	DropUpdates               bool     `json:"drop_updates"`
	OrdinaryPrefixes          []string `json:"cmd_prefixes"`
	CmdPrefixes               []rune   `json:"-"`
	BaseChats                 []int64  `json:"base_chats"`
	RateLimiterPunishmentTime int64    `json:"rate_limiter_punishment_time"`
	RateLimiterTimeout        int64    `json:"rate_limiter_timeout"`
	RateLimiterMaxMessages    int64    `json:"rate_limiter_max_messages"`
	RateLimiterMaxCache       int64    `json:"rate_limiter_max_cache"`
}
