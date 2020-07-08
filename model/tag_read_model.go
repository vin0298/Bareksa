package model

type TagReadModel struct {
	Uuid string
	Name string
	Id   int
}

type TagReadNoPKModel struct {
	Uuid string
	Name string
}
