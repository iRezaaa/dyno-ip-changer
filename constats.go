package main

type ErrorType int

const (
	SendRequest ErrorType = iota
	Unmarshal
	Marshal
	Response
	Auth
	Unknown
	ReadFile
	WriteFile
)

type RequestType int

const (
	GetDomainList RequestType = iota
	CheckDomainIsBlocked
	UpdateDomainIP
	GetNewIPFromFile
)
