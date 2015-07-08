package database

// OrderDirection ...
type OrderDirection int

// String ...
func (o OrderDirection) String() (out string) {
	switch o {
	case OrderAscending:
		out = "ASC"
	case OrderDescending:
		out = "DESC"
	}
	return
}

const (
	// OrderAscending ...
	OrderAscending OrderDirection = iota
	// OrderDescending ...
	OrderDescending
)
