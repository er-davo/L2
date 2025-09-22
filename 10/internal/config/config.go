package config

// Config of sorting
type Config struct {
	Column       int    // -k
	Delimiter    string // N
	Numeric      bool   // -n
	Reverse      bool   // -r
	Unique       bool   // -u
	Month        bool   // -M
	IgnoreBlanks bool   // -b
	CheckSorted  bool   // -c
	Human        bool   // -h
}
