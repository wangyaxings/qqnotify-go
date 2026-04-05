package smokeenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadFirst(paths ...string) (string, error) {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if err := LoadFile(path); err != nil {
				return "", err
			}
			return path, nil
		}
	}

	return "", fmt.Errorf("no smoke config file found in candidate paths")
}

func LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open smoke config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid env line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set env %s: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan smoke config file: %w", err)
	}

	return nil
}
