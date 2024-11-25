package helpers

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func NewNullString(value string) sql.NullString {
	if value != "" {
		return sql.NullString{
			String: value,
			Valid:  true,
		}
	}
	return sql.NullString{}
}

func NewNullUUID(value string) uuid.NullUUID {
	if value != "" {
		id, err := uuid.Parse(value)
		if err == nil {
			return uuid.NullUUID{
				UUID:  id,
				Valid: true,
			}
		}
	}
	return uuid.NullUUID{}
}

func NewNullDate(date string) sql.NullTime {
	if date == "" {
		return sql.NullTime{}
	}
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  parsedDate,
		Valid: true,
	}
}

func NewNullInt16(value int) sql.NullInt16 {
	if value == 0 {
		return sql.NullInt16{}
	}
	return sql.NullInt16{
		Int16: int16(value),
		Valid: true,
	}
}

func NewNullInt32(value int) sql.NullInt32 {
	if value == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: int32(value),
		Valid: true,
	}
}

func NewNullInt64(value int) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(value),
		Valid: true,
	}
}

func NewNullFloat64(value float64) sql.NullFloat64 {
	if value == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: value,
		Valid:   true,
	}
}
