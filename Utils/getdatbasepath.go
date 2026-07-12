package utils

import "os"

func GetDatabasePath() string {
	if envPath := os.Getenv("SQLITE_DB_PATH"); envPath != "" {
		return envPath
	}

	paths := []string{
		"/home/panda/Projects/SqLiteWireProtocol/db/test.db",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "db/test.db" // Fallback to relative path
}
