package util

import (
	"errors"
	"fmt"

	workermodel "github.com/HeRaNO/cdoj-execution-worker/model"
	"github.com/HeRaNO/cdoj-vjudge/config"
	"github.com/HeRaNO/cdoj-vjudge/model"
)

func MakeSubmission(req model.Submission, checker string, limit *model.Limitation) (workermodel.ExecRequest, error) {
	problemID := fmt.Sprint(req.ProblemID)

	codeName := ""
	execName := ""
	runName := ""
	compileArgs := make([]string, 0)
	runArgs := make([]string, 0)

	switch req.Language {
	case config.C:
		codeName, execName, runName = "main.c", "gcc", "main"
		compileArgs = append(compileArgs, "gcc", "main.c", "-o", "main", "-std=c11", "-O2", "-DONLINE_JUDGE", "-lm")
		runArgs = append(runArgs, "./main")
	case config.CPP:
		codeName, execName, runName = "main.cpp", "g++", "main"
		compileArgs = append(compileArgs, "g++", "main.cpp", "-o", "main", "-std=c++17", "-O2", "-DONLINE_JUDGE")
		runArgs = append(runArgs, "./main")
	case config.Java:
		codeName, execName, runName = "Main.java", "javac", "java"
		compileArgs = append(compileArgs, "javac", "-encoding", "UTF-8", "-sourcepath", ".", "-d", ".", "Main.java")
		javaMemLim := fmt.Sprintf("-Xss%dm -Xms%dm -Xmx%dm", limit.MemoryLimit, limit.MemoryLimit, limit.MemoryLimit)
		runArgs = append(runArgs, "java", "-Dfile.encoding=UTF-8", "-XX:+UseSerialGC", javaMemLim, "Main")
	case config.Python3:
		codeName, execName, runName = "main.py", "python3", "python3"
		compileArgs = append(compileArgs, "python3", "-m", "py_compile", "main.py")
		runArgs = append(runArgs, "python3", "__pycache__/main.cpython-37.pyc")
	default:
		return workermodel.ExecRequest{}, errors.New("unknown language code")
	}

	return workermodel.ExecRequest{
		CompilePhases: workermodel.CompilePhase{
			Compile: workermodel.Phase{
				Exec:    execName,
				RunArgs: compileArgs,
				Limits: workermodel.Limitation{
					Time:   5000,
					Memory: 1024 << 20,
				},
			},
			SourceCode: workermodel.SourceCodeDescriptor{
				Name:    codeName,
				Content: req.Code,
			},
			ExecName: execName,
		},
		RunPhases: workermodel.RunPhase{
			Run: workermodel.Phase{
				Exec:    runName,
				RunArgs: runArgs,
				Limits: workermodel.Limitation{
					Time:   int32(limit.TimeLimit),
					Memory: limit.MemoryLimit,
				},
			},
			ProblemID: problemID,
		},
		CheckPhase: checker,
	}, nil
}

func MakeCEResult(msg *workermodel.OmitString) model.Result {
	realMsg := msg.S
	if msg.OmitSize != 0 {
		realMsg = fmt.Sprintf("%s (%d byte omitted)", msg.S, msg.OmitSize)
	}
	return model.Result{
		Verdict: "Compile Error",
		Message: realMsg,
	}
}

func MakeREResult(errMsg string, resp workermodel.ExecResult) model.Result {
	verdict := fmt.Sprintf("Runtime Error on test case %d", resp.Case)
	return model.Result{
		Verdict:    verdict,
		Message:    errMsg,
		TimeUsed:   resp.UserTimeUsed,
		MemoryUsed: resp.MemoryUsed,
	}
}

func MakeIEResult() model.Result {
	return model.Result{
		Verdict: "Internal Error",
		Message: "Something wrong with cdoj-vjudge T_T",
	}
}

func MakeOKResult(errMsg string, rep interface{}) model.Result {
	if errMsg == "success" {
		resp := rep.(workermodel.ExecResult)
		return model.Result{
			Verdict:    "Accepted",
			TimeUsed:   resp.UserTimeUsed,
			MemoryUsed: resp.MemoryUsed,
		}
	} else if errMsg == "running" {
		resp := rep.(int)
		verdict := fmt.Sprintf("Running on test case %d", resp)
		return model.Result{
			Verdict: verdict,
		}
	} else { // "wrong answer"
		resp := rep.(workermodel.ExecResult)
		verdict := fmt.Sprintf("Wrong on test case %d", resp.Case)
		return model.Result{
			Verdict:    verdict,
			Message:    resp.CheckerResult.S,
			TimeUsed:   resp.UserTimeUsed,
			MemoryUsed: resp.MemoryUsed,
		}
	}
}

func MakeUnknownResult(errMsg string) model.Result {
	return model.Result{
		Verdict: "Internal Error",
		Message: errMsg,
	}
}
