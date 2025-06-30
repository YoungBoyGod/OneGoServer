package config

import "fmt"

// GetDsn 获取数据库连接字符串
func (c *DatabaseConfig) GetDsn() string {
	switch c.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
			c.Username, c.Password, c.Host, c.Port, c.DBName, c.Charset, c.ParseTime, c.Loc)
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
			c.Host, c.Username, c.Password, c.DBName, c.Port, c.Loc)
	case "sqlite":
		if c.DBName == ":memory:" {
			return ":memory:"
		}
		return c.DBName + ".db"
	default:
		return ""
	}
}
