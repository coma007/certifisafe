package model

import "time"

type Request struct {
	Id                int32
	parentCertificate Certificate
	certificate       Certificate
	datetime          time.Time
	status            Status
}

type Status int64

const (
	Pending Status = iota
	Accepted
	Rejected
)
