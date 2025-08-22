package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Interface interface {
}

type Config struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Pass            string        `mapstructure:"pass"`
	Db              string        `mapstructure:"db"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	TimeZone        string        `mapstructure:"timezone"`
	MaxIdleConn     int           `mapstructure:"max_idle_conn"`
	MaxOpenConn     int           `mapstructure:"max_open_conn"`
	MaxLifeTimeConn time.Duration `mapstructure:"max_lifetime_conn"`
	MaxIdleTimeConn time.Duration `mapstructure:"max_idletime_conn"`
}

type Postgres struct {
	config *Config
	Client *gorm.DB
	exitCh chan struct{}
}

func New(c *Config) (*Postgres, error) {
	if err := c.check(); err != nil {
		return nil, fmt.Errorf("invalid config:%s", err)
	}

	url := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s TimeZone=%s",
		c.Host, c.Port, c.User, c.Db, c.Pass, c.SSLMode, c.TimeZone)
	
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConn)
	db.DB().SetMaxOpenConns(c.MaxOpenConn)
	db.DB().SetConnMaxIdleTime(c.MaxIdleTimeConn)
	db.DB().SetConnMaxLifetime(c.MaxLifeTimeConn)

	client := &Postgres{
		config: c,
		Client: db,
		exitCh: make(chan struct{}),
	}
	return client, nil
}

func (c *Config) check() error {
	if "" == c.User || "" == c.Pass || "" == c.Host || "" == c.Db {
		return fmt.Errorf("config error")
	}

	if 0 == c.Port {
		c.Port = 5432
	}

	if 0 == c.MaxIdleConn || 0 == c.MaxOpenConn {
		c.MaxIdleConn = 80
		c.MaxOpenConn = 80
	}

	if c.MaxIdleConn != c.MaxOpenConn {
		c.MaxOpenConn = c.MaxIdleConn
	}

	if c.MaxLifeTimeConn == time.Duration(0) {
		c.MaxLifeTimeConn = 600 * time.Second
	}

	if c.MaxIdleTimeConn == time.Duration(0) {
		c.MaxIdleTimeConn = 600 * time.Second
	}

	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	if c.TimeZone == "" {
		c.TimeZone = "Asia/Shanghai"
	}

	return nil
}

func (postgres *Postgres) GetConfig() *Config {
	return postgres.config
}

func (postgres *Postgres) Close() error {
	close(postgres.exitCh)
	return postgres.Client.Close()
}