package mongo

// Config represents MongoDB connection configuration
type Config struct {
	// URI is the MongoDB connection string
	URI string
	// Database is the database name to connect to
	Database string
	// Username for authentication (optional if included in URI)
	Username string
	// Password for authentication (optional if included in URI)
	Password string
	// AuthDatabase is the authentication database (optional)
	AuthDatabase string
	// ReplicaSet is the replica set name (optional)
	ReplicaSet string
	// MaxPoolSize is the maximum number of connections in the connection pool
	MaxPoolSize uint64
	// MinPoolSize is the minimum number of connections in the connection pool
	MinPoolSize uint64
	// ConnectTimeout is the timeout for connecting to MongoDB in seconds
	ConnectTimeout int
	// SocketTimeout is the timeout for socket operations in seconds
	SocketTimeout int
	// ServerSelectionTimeout is the timeout for server selection in seconds
	ServerSelectionTimeout int
	// HeartbeatInterval is the interval between server monitoring in seconds
	HeartbeatInterval int
	// LocalThreshold is the acceptable latency difference for selecting servers in milliseconds
	LocalThreshold int
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		URI:                    "mongodb://localhost:47017",
		MaxPoolSize:            100,
		MinPoolSize:            0,
		ConnectTimeout:         30,
		SocketTimeout:          30,
		ServerSelectionTimeout: 30,
		HeartbeatInterval:      10,
		LocalThreshold:         15,
	}
}
