package document

import (
	"encoding/csv"
	"io"
	"strconv"
	"warehouse_control/internal/entity"

	"github.com/wb-go/wbf/zlog"
)

// SaveProductHistoryToCSVWriter - сохранить историю продукта в CSV в указанный writer
func SaveProductHistoryToCSVWriter(w io.Writer, history []entity.ProductLogs) error {
	// BOM для Excel
	if _, err := w.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return err
	}

	writer := csv.NewWriter(w)
	writer.Comma = ';'
	defer writer.Flush()

	headers := []string{"ProductID", "OldName", "NewName", "OldDescription", "NewDescription", "OldQuantity", "NewQuantity", "ChangedAt"}
	if err := writer.Write(headers); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Error writing headers")
		return err
	}

	for _, log := range history {
		row := []string{
			log.ProductID.String(),
			log.OldName,
			log.NewName,
			log.OldDescription,
			log.NewDescription,
			strconv.Itoa(log.OldQuantity),
			strconv.Itoa(log.NewQuantity),
			log.ChangedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			zlog.Logger.Error().Err(err).Msg("Error writing to csv")
			return err
		}
	}

	return nil
}
