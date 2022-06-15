package config

const (
	ErrInternal  = "1000"
	ErrWrongInfo = "1001"
)

type Language int8

const (
	C Language = iota
	CPP
	Java
	Python3
)
