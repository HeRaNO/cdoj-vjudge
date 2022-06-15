package model

import "github.com/HeRaNO/cdoj-vjudge/config"

type Submission struct {
	ProblemID int64           `json:"problem_id"`
	Language  config.Language `json:"lang"`
	Code      string          `json:"code"`
}
