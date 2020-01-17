package conf

import (
	"github.com/srlemon/contrib/config"
	"github.com/srlemon/userDetail/utils"
)

// ProjectConfig
type ProjectConfig struct {
	config.Config `yaml:"-"`
	ServerAddr    string    `yaml:"serverAddr"`
	ServerPort    string    `yaml:"serverPort"`
	Sync          bool      `yaml:"sync"`        // true: 数据库同步
	Mode          string    `yaml:"mode"`        // mode: 方式
	UserKeyDir    string    `yaml:"userKeyDir"`  // 私钥地址
	UserPubDir    string    `yaml:"userPubDir"`  // 公钥地址
	AdminKeyDir   string    `yaml:"adminKeyDir"` // 私钥地址
	AdminPubDir   string    `yaml:"adminPubDir"` // 公钥地址
	Db            *Database `yaml:"db"`
	RDB           *RedisDB  `yaml:"rdb"`
	IsTLS         bool      `yaml:"isTLS"` // true: 开启https
	TLS           *TLS      `yaml:"tls"`
}

// RedisDB 缓存的连接参数
type RedisDB struct {
	Host     string
	Port     string
	DB       int
	Password string
}

// TLS
type TLS struct {
	Cert string `json:"cert"` // 证书文件
	Key  string `json:"key"`  // Key 文件
}

// Database 数据库连接参数
type Database struct {
	Host         string
	Port         string
	Driver       string
	DatabaseName string
	Username     string
	Password     string
	MaxIdleConn  int
	MaxOpenConn  int
}

const (
	ModeProduction = "production"
)

var (
	ProjectSetting = &ProjectConfig{
		ServerAddr: utils.PubGetEnvString("SERVER_ADDR", "127.0.0.1"),
		ServerPort: utils.PubGetEnvString("SERVER_PORT", "8060"),
		Sync:       false,
		Mode:       "production",
		IsTLS:      false,
		Db: &Database{
			Driver:       utils.PubGetEnvString("DB_DRIVER", "postgres"),
			Host:         utils.PubGetEnvString("DB_HOST", "127.0.0.1"),
			Port:         utils.PubGetEnvString("DB_PORT", "65432"),
			DatabaseName: utils.PubGetEnvString("DB_NAME", "business"),
			Username:     utils.PubGetEnvString("DB_USERNAME", "business"),
			Password:     utils.PubGetEnvString("DB_PASSWORD", "business"),
			MaxIdleConn:  200,
			MaxOpenConn:  2000,
		},
		RDB: &RedisDB{
			Host:     utils.PubGetEnvString("RDB_HOST", "127.0.0.1"),
			Port:     utils.PubGetEnvString("RDB_PORT", "6379"),
			DB:       0,
			Password: utils.PubGetEnvString("RDB_PASSWORD", "business"),
		},
		UserKeyDir:  "./conf/user.key",
		UserPubDir:  "./conf/user.pub",
		AdminKeyDir: "./conf/admin.key",
		AdminPubDir: "./conf/admin.pub",
	}
)

// InitConfig 初始化配置文件
func InitConfig() (err error) {
	if err = config.LoadConfiguration("./conf/project.config.yaml", ProjectSetting, ProjectSetting); err != nil {
		return
	}
	if ProjectSetting.IsTLS {
		ProjectSetting.TLS = &TLS{
			Cert: "./conf/serve.cert",
			Key:  "./conf/serve.key",
		}
	}
	if err = ProjectSetting.Save(nil); err != nil {
		return
	}
	return
}
