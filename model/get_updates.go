package model

type GetUpdatesRq struct {
	Offset  int64 `json:"offset"`
	Timeout int   `json:"timeout"`
}

type GetUpdatesRs []Update
