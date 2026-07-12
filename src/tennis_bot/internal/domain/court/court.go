package court

type Court struct {
	ID        int64
	Name      string
	OpenHour  string
	CloseHour string
	Address   string
	IsActive  bool
}
