package helpers

type SkipCondition struct {
	Method string
	Path   string
}

func GetSkipConditions() []SkipCondition {
	return []SkipCondition{
		{Method: "POST", Path: "/login"},
		{Method: "POST", Path: "/customer"},
		{Method: "GET", Path: "/verify"},
		{Method: "GET", Path: "/swagger/*"},
		{Method: "GET", Path: "/customers"},
	}
}
