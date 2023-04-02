package model

import (
	"time"
)

type Request struct {
	Id                int
	ParentCertificate *Certificate
	Certificate       *Certificate
	Datetime          time.Time
	Status            RequestStatus
}

type RequestStatus int64

const (
	PENDING RequestStatus = iota
	ACCEPTED
	REJECTED
)
