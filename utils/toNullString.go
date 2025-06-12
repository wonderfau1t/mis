package utils

import "database/sql"

func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
