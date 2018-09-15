package model

import "strconv"

// Snowflake is a discord data type representing a globally unique ID
type Snowflake uint64

// ParseSnowflake parses the input snowflake string
func ParseSnowflake(input string) (Snowflake, error) {
	val, err := strconv.ParseUint(input, 10, 64)
	return Snowflake(val), err
}

// Format stringifies a string.
func (s Snowflake) Format() string {
	return strconv.FormatUint(uint64(s), 10)
}
