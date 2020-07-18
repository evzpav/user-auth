package domain

type StorageSortType string

const (
	SortAsc  StorageSortType = "asc"
	SortDesc StorageSortType = "desc"
)

type StorageSort struct {
	Attribute string
	Type      StorageSortType
}
