package dal

import (
	"context"
	"errors"
	"sync"

	"github.com/HeRaNO/cdoj-vjudge/config"
	"github.com/HeRaNO/cdoj-vjudge/model"
	"github.com/HeRaNO/cdoj-vjudge/util"
)

func GetSamplesByProblemID(ctx context.Context, problemID int64) ([]model.ProblemSample, error) {
	rdb := config.RDB
	samples := make([]model.ProblemSample, 0)
	result := rdb.WithContext(ctx).Model(&model.Sample{}).Table(model.TableProblemSample).Where("problem_id = ?", problemID).Find(&samples)
	if result.Error != nil {
		return nil, result.Error
	}
	return samples, nil
}

func GetLimitationByLimitationID(ctx context.Context, limitationID int64) (*model.Limitation, error) {
	rdb := config.RDB
	limitation := model.Limitation{}
	result := rdb.WithContext(ctx).Model(&model.Limitation{}).Table(model.TableProblemLimitation).Where("id = ?", limitationID).Find(&limitation)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no problem limitation record")
	}
	if result.RowsAffected > 1 {
		return nil, errors.New("duplicate limitation_id but why???")
	}

	return &limitation, nil
}

func GetStatementByProblemID(ctx context.Context, problemID int64) (*model.Statement, error) {
	rdb := config.RDB
	problem := model.Problem{}
	result := rdb.WithContext(ctx).Model(&model.Problem{}).Table(model.TableProblem).Where("id = ?", problemID).Find(&problem)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no problem record")
	}
	if result.RowsAffected > 1 {
		return nil, errors.New("duplicate problem_id but why???")
	}

	if problem.IsDisable == 1 {
		return nil, errors.New("the problem is not published")
	}

	wg := sync.WaitGroup{}
	oneErr := util.OneError{}
	samples := make([]model.ProblemSample, 0)
	limitation := &model.Limitation{}

	wg.Add(1)
	go func(wg *sync.WaitGroup, oneErr *util.OneError) {
		defer wg.Done()
		limit, err := GetLimitationByLimitationID(ctx, problem.LimitationID)
		if err != nil {
			oneErr.Add(err)
		}
		limitation = limit
	}(&wg, &oneErr)

	wg.Add(1)
	go func(wg *sync.WaitGroup, oneErr *util.OneError) {
		defer wg.Done()
		sample, err := GetSamplesByProblemID(ctx, problemID)
		if err != nil {
			oneErr.Add(err)
		}
		samples = sample
	}(&wg, &oneErr)

	wg.Wait()

	if oneErr.Err != nil {
		return nil, oneErr.Err
	}
	return model.MakeStatement(problem, limitation, samples), nil
}

func GetProblemInfo(ctx context.Context, problemName *string, offset int64, limit int64) ([]model.ProblemInfo, error) {
	rdb := config.RDB
	problemInfo := make([]model.ProblemInfo, 0)
	result := rdb.WithContext(ctx).Model(&model.Problem{}).Table(model.TableProblem).Select([]string{"id", "title"})
	if *problemName != "" {
		result = result.Where("title like ?", "%"+*problemName+"%")
	}
	result = result.Where("disable = 0").Order("id").Offset(int(offset)).Limit(int(limit)).Scan(&problemInfo)
	if result.Error != nil {
		return nil, result.Error
	}

	return problemInfo, nil
}

func GetProblemCheckerInfo(ctx context.Context, problemID int64) (string, *model.Limitation, error) {
	rdb := config.RDB
	problem := model.Problem{}
	result := rdb.WithContext(ctx).Model(&model.Problem{}).Table(model.TableProblem).Where("id = ?", problemID).Find(&problem)
	if result.Error != nil {
		return "", nil, result.Error
	}
	if result.RowsAffected == 0 {
		return "", nil, errors.New("no problem record")
	}
	if result.RowsAffected > 1 {
		return "", nil, errors.New("duplicate problem_id but why???")
	}

	if problem.IsDisable == 1 {
		return "", nil, errors.New("the problem is not published")
	}

	limit, err := GetLimitationByLimitationID(ctx, problem.LimitationID)
	if err != nil {
		return "", nil, err
	}

	return problem.Checker, limit, nil
}
