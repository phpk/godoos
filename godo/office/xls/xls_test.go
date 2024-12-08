package xls

import (
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	if xlFile, err := Open("t1.xls", "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			fmt.Println("Total Lines ", sheet1.MaxRow, sheet1.Name)
			for i := 265; i <= 267; i++ {
				fmt.Printf("row %v point %v \n", i, sheet1.Row(i))
				if sheet1.Row(i) == nil {
					continue
				}
				row := sheet1.Row(i)
				for index := row.FirstCol(); index < row.LastCol(); index++ {
					fmt.Println(index, "==>", row.Col(index), " ")
					fmt.Printf("%T\n", row.cols[uint16(index)])
				}
				// col1 := .Cols[0]
				// col2 := sheet1.Row(uint16(i)].Cols[1]
				// fmt.Printf("\ncol1 %v \nCol2 %v \n", col1.String(xlFile), col2.String(xlFile))
			}
		}
	}
}
