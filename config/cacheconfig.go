package config

//DatabaseConfiguration holds the configuration for the database
type CacheConfiguration struct {
	//Memory to be reserved for caching
	Memory int
	//Cache replacement stagergy to be used. Default: LRU
	CRS string
}
