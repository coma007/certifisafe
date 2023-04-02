package model

import "time"

type Request struct {
	Id                int32
	ParentCertificate *Certificate
	Certificate       *Certificate
	Datetime          time.Time
	Status            RequestStatus
}

type RequestStatus int64

const (
	Pending RequestStatus = iota
	Accepted
	Rejected
)
