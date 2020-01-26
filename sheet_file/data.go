package sheet_file

import (
	"github.com/tealeg/xlsx"
)

func CreateFile() *xlsx.File {
	data, _ := Asset("base_sheet.xlsx")
	file, _ := xlsx.OpenBinary(data)
	return file
}
