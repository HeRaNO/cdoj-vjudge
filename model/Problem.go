package model

const (
	TableProblem           = "problem_problem"
	TableProblemSample     = "problem_problemsample"
	TableProblemLimitation = "problem_limitation"
)

type Problem struct {
	ID           int64  `gorm:"column:id" json:"id"`
	Title        string `gorm:"column:title" json:"title"`
	Content      string `gorm:"column:content" json:"content"`
	Resources    string `gorm:"column:resources" json:"resources"`
	Constraints  string `gorm:"column:constraints" json:"constraints"`
	Input        string `gorm:"column:standard_input" json:"input"`
	Output       string `gorm:"column:standard_output" json:"output"`
	Note         string `gorm:"column:note" json:"note"`
	IsDisable    int    `gorm:"column:disable" json:"disable"`
	Submit       int32  `gorm:"column:submit" json:"submit"`
	Accept       int32  `gorm:"column:accept" json:"accept"`
	Checker      string `gorm:"column:_checker" json:"checker"`
	LimitationID int64  `gorm:"column:limitation_id" json:"limitation_id"`
}

type Limitation struct {
	ID          int64 `gorm:"column:id" json:"id"`
	TimeLimit   int64 `gorm:"column:time_limit" json:"time_limit"`
	MemoryLimit int64 `gorm:"column:memory_limit" json:"memory_limit"`
	OutputLimit int64 `gorm:"column:output_limit" json:"output_limit"`
	CpuLimit    int64 `gorm:"column:cpu_limit" json:"cpu_limit"`
}

type ProblemSample struct {
	ID           int64  `gorm:"column:id" json:"id"`
	SampleInput  string `gorm:"column:input_content" json:"sample_input"`
	SampleOutput string `gorm:"column:output_content" json:"sample_output"`
	ProblemID    int64  `gorm:"column:problem_id" json:"problem_id"`
}

type Sample struct {
	SampleInput  string `json:"sample_input"`
	SampleOutput string `json:"sample_output"`
}

type Statement struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Resources   string   `json:"resources"`
	Constraints string   `json:"constraints"`
	Input       string   `json:"input"`
	Output      string   `json:"output"`
	Note        string   `json:"note"`
	Checker     string   `json:"checker"`
	Samples     []Sample `json:"samples"`
	TimeLimit   int64    `json:"time_limit"`
	MemoryLimit int64    `json:"memory_limit"`
	Submit      int32    `json:"submit"`
	Accept      int32    `json:"accept"`
}

type ProblemInfo struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

func MakeStatement(problem Problem, limitation *Limitation, problemSamples []ProblemSample) *Statement {
	samples := make([]Sample, len(problemSamples))
	for i, sample := range problemSamples {
		samples[i] = Sample{
			SampleInput:  sample.SampleInput,
			SampleOutput: sample.SampleOutput,
		}
	}
	return &Statement{
		ID:          problem.ID,
		Title:       problem.Title,
		Content:     problem.Content,
		Resources:   problem.Resources,
		Constraints: problem.Constraints,
		Input:       problem.Input,
		Output:      problem.Output,
		Note:        problem.Note,
		Checker:     problem.Checker,
		Samples:     samples,
		Submit:      problem.Submit,
		Accept:      problem.Accept,
		TimeLimit:   limitation.TimeLimit,
		MemoryLimit: limitation.MemoryLimit,
	}
}
