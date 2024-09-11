package conf

type Log struct {
	Level      string `json:"Level,omitempty"`
	Output     string `json:"Output,omitempty"`
	MaxBackups int    `json:"MaxBackups,omitempty"`
	MaxSize    int    `json:"MaxSize,omitempty"`
	MaxAge     int    `json:"MaxAge,omitempty"`
}

func newLog() Log {
	return Log{
		Level:      "info",
		Output:     "",
		MaxBackups: 3,
		MaxSize:    100,
		MaxAge:     28,
	}
}
