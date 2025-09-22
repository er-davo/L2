package config

// Config for grep
type Config struct {
	After       int
	Before      int
	Context     int
	CountOnly   bool
	IgnoreCase  bool
	InvertMatch bool
	Fixed       bool
	LineNum     bool
	Pattern     string
	File        string
}
