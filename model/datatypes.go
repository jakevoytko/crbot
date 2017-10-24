package model

import "strconv"

type Snowflake uint64

func ParseSnowflake(input string) (Snowflake, error) {
	val, err := strconv.ParseUint(input, 10, 64)
	return Snowflake(val), err
}

func (s Snowflake) Format() string {
	return strconv.FormatUint(uint64(s), 10)
}
