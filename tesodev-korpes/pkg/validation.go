package pkg

import "strconv"

func LimitOffsetValidation(l string, o string) (int64, int64) {

	defaultLimit := int64(25)
	defaultOffset := int64(0)
	
	offset, err := strconv.ParseInt(o, 10, 64)
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	limit, err := strconv.ParseInt(l, 10, 64)
	if err != nil || limit <= 0 || limit > 100 {
		limit = defaultLimit
	}

	if l == "" {
		limit = defaultLimit
	}

	if o == "" {
		offset = defaultOffset
	}

	return limit, offset
}
