package internal

var testUsers = map[string]int64{
	"jo.example@example.com": 1001,
	"vip@example.com":        1002,
}

func LookupUserByEmail(email string) int64 {
	if id, ok := testUsers[email]; ok {
		return id
	}
	return 0 // не найден
}
