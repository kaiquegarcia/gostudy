package main

import (
	"os"
	"time"

	"github.com/kaiquegarcia/gostudy/v2/logging"
	"github.com/kaiquegarcia/gostudy/v2/planner"
	"github.com/kaiquegarcia/gostudy/v2/utils"
)

func main() {
	// dependencies
	logger := logging.NewLogger(
		logging.DefaultPrinter,
		logging.LevelDebug,
	)

	defer utils.PanicHandler(logger)

	reqFilenames := utils.RequiredFilenames{
		HourGrade:       "hour_grade.csv",
		DisciplinesList: "disciplines.csv",
		Output:          "planner.csv",
	}

	// run
	var startDate time.Time
	if len(os.Args) > 1 {
		logger.Debug("prepare to parse date from os.Args[1]")
		d, err := time.Parse(planner.LayoutDateOnly, os.Args[1])
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
	hourGrade, err := planner.NewHourGradeFromRow(hourGradeRecords)
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
	disciplines, err := planner.NewDisciplineFromRows(disciplineRecords)
	if err != nil {
		logger.Error(err, "could not extract disciplines list from table records")
		return
	}

	makerReady := false
	defer func() {
		if makerReady {
			return
		}
		// fallback if something go wrong before initializing maker
		for _, d := range disciplines {
			d.Close()
		}
	}()

	logger.Debug("disciplines list data extracted successfuly, initializing planner maker")
	maker, err := planner.NewMaker(logger, hourGrade, disciplines, startDate, reqFilenames.Output)
	if err != nil {
		logger.Error(err, "could not initialize planner maker")
		return
	}
	defer maker.Close()
	makerReady = true

	logger.Debug("preparing to mount planner")
	err = maker.Mount()
	if err != nil {
		logger.Error(err, "could not mount planner")
		return
	}

	logger.Debug("planner mounted successfully")
}
