package analytics

import (
	"encoding/csv"
	"os"
	"sales_tracker/internal/entity"
	"strconv"

	"github.com/wb-go/wbf/zlog"
)

// SaveAnalyticsToCSV - сохраняем аналитику в csv файл
func SaveAnalyticsToCSV(filename string, result entity.AnalyticsResult) error {
	// создаем файл
	file, err := os.Create(filename)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("filename", filename).Msg("SaveAnalyticsToCSV: Failed to create file")
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("SaveAnalyticsToCSV: Failed to close file")
		}
	}(file)

	// записываем данные в формате csv
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	headers := []string{"TotalCount", "TotalSum", "AvgAmount", "Median", "P90"}
	err = writer.Write(headers)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("SaveAnalyticsToCSV: Failed to write headers to file")
		return err
	}

	record := []string{
		strconv.Itoa(int(result.TotalCount)),
		strconv.FormatFloat(result.TotalSum, 'f', 2, 64),
		strconv.FormatFloat(result.AvgAmount, 'f', 2, 64),
		strconv.FormatFloat(result.Median, 'f', 2, 64),
		strconv.FormatFloat(result.P90, 'f', 2, 64),
	}

	return writer.Write(record)
}
