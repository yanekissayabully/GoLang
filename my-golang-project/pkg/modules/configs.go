package modules

import "time"

// PostgreSQL конфиг
type PostgreSQL struct {
    Host         string
    Port         string
    Username     string
    Password     string
    DBName       string
    SSLMode      string
    ExecTimeout  time.Duration // Добавим таймаут, который был в методичке
}