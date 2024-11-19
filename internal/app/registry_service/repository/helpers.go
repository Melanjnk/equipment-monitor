package repository

import (
	"database/sql"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
)

func figure[T model.EquipmentKind|model.OperationalStatus](f T) byte {
	return byte(f + '0')
} 

func joinIntegralArray[T model.EquipmentKind|model.OperationalStatus](intArray []T) string {
	const (
		zero byte	= '0'
		separator	= ','
	)
	l := len(intArray)
	buffer := make([]byte, l + l - 1)
	buffer[0] = figure(intArray[0])
	for i := 1; i < l; i++ {
		buffer[i + i - 1] = separator
		buffer[i + i] = figure(intArray[i])
	}
	return string(buffer)
}

func checkAffect(result sql.Result, err error) (bool, error) {
	if err == nil {
		var count int64
		if count, err = result.RowsAffected(); err == nil {
			return count > 0, nil
		}
	}
	return false, err
}
