package cache

import "github.com/op/go-logging"

var log = logging.MustGetLogger("cache_manager")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type Manager struct {
}
