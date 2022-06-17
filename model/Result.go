package model

const (
	SubmissionResultKey = "RESULT"
)

type Result struct {
	Verdict    string `json:"verdict"`
	Message    string `json:"msg"`
	TimeUsed   int64  `json:"time_used"`
	MemoryUsed int64  `json:"mem_used"`
}
