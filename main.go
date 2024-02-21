package main

import (
	"os"
	"sync"
	"time"

	"github.com/kaiquegarcia/gostudy/gostudy"
	"github.com/kaiquegarcia/gostudy/utils"
)

func main() {
	// dependencies
	logger := utils.NewLogger(
		utils.DefaultPrinter,
		utils.LevelDebug,
	)

	defer utils.PanicHandler(logger)

	reqFilenames := utils.RequiredFilenames{
		HourGrade:       "hour_grade.csv",
		DisciplinesList: "disciplines.csv",
		PlannerOutput:   "planner.csv",
	}

	// run
	var startDate time.Time
	if len(os.Args) > 1 {
		logger.Debug("prepare to parse date from os.Args[1]")
		d, err := time.Parse(gostudy.LayoutDateOnly, os.Args[1])
		if err != nil {
			logger.Error(err, "could not parse os.Args[1]")
			return
		}

		logger.Debug("date parsed successfully")
		startDate = d
	} else {
		logger.Debug("using now + 6days as startDate")
		startDate = time.Now().AddDate(0, 0, 6)
	}

	logger.Debug("reading '%s'", reqFilenames.HourGrade)
	hourGradeRecords, err := utils.ReadCSV(reqFilenames.HourGrade)
	if err != nil {
		logger.Error(err, "could not read '%s'", reqFilenames.HourGrade)
		return
	}

	logger.Debug("'%s' readed successfuly, preparing to extract information from records", reqFilenames.HourGrade)
	hourGrade, err := gostudy.ExtractHourGradeFromTableRecords(hourGradeRecords)
	if err != nil {
		logger.Error(err, "could not extract hour grade from table records")
		return
	}

	logger.Debug("hour grade extracted successfully")
	logger.Debug("reading '%s'", reqFilenames.DisciplinesList)
	disciplineRecords, err := utils.ReadCSV(reqFilenames.DisciplinesList)
	if err != nil {
		logger.Error(err, "could not read '%s'", reqFilenames.DisciplinesList)
		return
	}

	logger.Debug("'%s' readed successfuly, preparing to extract information from records", reqFilenames.DisciplinesList)
	disciplines, err := gostudy.ExtractDisciplineListFromTableRecords(disciplineRecords)
	if err != nil {
		logger.Error(err, "could not extract disciplines list from table records")
		return
	}

	logger.Debug("disciplines list data extracted successfuly, preparing to read disciplines content files (%d disciplines)", len(disciplines))
	wg := &sync.WaitGroup{}
	wg.Add(len(disciplines))
	hadErrors := false
	for _, d := range disciplines {
		go func(discipline *gostudy.Discipline) {
			logger.Debug("reading '%s'", discipline.Filename)
			contentRecords, err := utils.ReadCSV(discipline.Filename)
			if err != nil {
				logger.Error(err, "could not read '%s'", discipline.Filename)
				hadErrors = true
				wg.Done()
				return
			}

			// if any other goroutine had errors before this
			if hadErrors {
				wg.Done()
				return
			}

			logger.Debug("'%s' readed successfuly, preparing to extract information from records", discipline.Filename)
			err = gostudy.ExtractDisciplineContentFromTableRecords(discipline, contentRecords)
			if err != nil {
				logger.Error(err, "could not extract '%s's content from table records", discipline.Name)
				hadErrors = true
				wg.Done()
				return
			}

			logger.Debug("'%s's content extracted from table records successfully", discipline.Name)
			wg.Done()
		}(d)
	}
	wg.Wait()
	if hadErrors {
		logger.Warn("one or more discipline's goroutines had errors, check the unordered logs to discover what happened")
		return
	}

	logger.Debug("all disciplines contents extracted successfully, initializing planner")
	planner := gostudy.NewPlanner(logger, hourGrade, disciplines, startDate)

	logger.Debug("preparing to mount planner")
	err = planner.Mount()
	if err != nil {
		logger.Error(err, "could not mount planner")
		return
	}

	logger.Debug("planner mounted successfully, retrieving result records")
	plannerRecords := planner.ResultRecords()

	logger.Debug("result records retrieved successfully, writing '%s'", reqFilenames.PlannerOutput)
	err = utils.WriteCSV(reqFilenames.PlannerOutput, plannerRecords)
	if err != nil {
		logger.Error(err, "could not write '%s'", reqFilenames.PlannerOutput)
		return
	}

	logger.Debug("'%s' written successfully", reqFilenames.PlannerOutput)
	logger.Info("procedure finished successfully")
}
