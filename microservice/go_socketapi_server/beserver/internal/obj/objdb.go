package obj

type ObjectDb struct {
	TestObj map[int]Test
}

func InitObjectDb() *ObjectDb {
	// init database

	return &ObjectDb{}

}
