package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func combineAllCSVInFolderToXLSX(folderPath, xlsxOutputPath string) error {
	// Define required keywords.
	checkList := []string{"per", "stockdata", "monthlyrevenue", "cashflow"}
	// Read folder files.
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	// Ensure each required keyword is found in at least one file name.
	for _, check := range checkList {
		found := false
		for _, file := range files {
			if !file.IsDir() &&
				strings.Contains(strings.ToLower(file.Name()), strings.ToLower(check)) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing file containing: %s", check)
		}
	}

	// Combine the CSV data from the two groups.
	perAndStockData := combineTwoCSV(
		filepath.Join(folderPath, "per.csv"),
		filepath.Join(folderPath, "stockdata.csv"),
	)
	monthlyRevenueAndCashflow := combineTwoCSV(
		filepath.Join(folderPath, "monthlyrevenue.csv"),
		filepath.Join(folderPath, "cashflow.csv"),
	)

	// Create a new XLSX file.
	f := excelize.NewFile()
	// Sheet1: write perAndStockData.
	sheet1 := f.GetSheetName(f.GetActiveSheetIndex())
	if err := addPERTitle(f, sheet1); err != nil {
		log.Fatalf("Error adding PER title: %v", err)
	}
	if err := addPERHeader(f, sheet1); err != nil {
		log.Fatalf("Error adding PER header: %v", err)
	}
	if err := addStockTitle(f, sheet1); err != nil {
		log.Fatalf("Error adding stock title: %v", err)
	}
	if err := addStockHeader(f, sheet1); err != nil {
		log.Fatalf("Error adding stock header: %v", err)
	}

	dataStartRow := 4
	for i, row := range perAndStockData {
		for j, cellValue := range row {
			cellName, err := excelize.CoordinatesToCellName(j+1, i+dataStartRow)
			if err != nil {
				return fmt.Errorf("failed to get cell name: %v", err)
			}
			if err := f.SetCellValue(sheet1, cellName, cellValue); err != nil {
				return fmt.Errorf("failed to set cell value: %v", err)
			}
		}
	}

	// Sheet2: create and write monthlyRevenueAndCashflow.
	sheet2Name := "Sheet2"

	_, err = f.NewSheet(sheet2Name)
	if err != nil {
		return fmt.Errorf("failed to create new sheet: %v", err)
	}

	// --- Titles ---
	// Set Revenue title (in row 1, columns A:Q)
	if err := addRevenueTitle(f, sheet2Name); err != nil {
		log.Fatalf("Error adding revenue title: %v", err)
	}
	// Set Quarterly title (in row 1, columns S:AK)
	if err := addCashflowTitle(f, sheet2Name); err != nil {
		log.Fatalf("Error adding quarterly title: %v", err)
	}

	// --- Headers ---
	// Revenue header: rows 2–4, columns A–Q.
	if err := addRevenueHeader(f, sheet2Name); err != nil {
		log.Fatalf("Error adding revenue header: %v", err)
	}
	// Quarterly header: rows 2–3, columns S–AK.
	if err := addCashflowHeader(f, sheet2Name); err != nil {
		log.Fatalf("Error adding quarterly header: %v", err)
	}
	dataStartRow = 5

	for i, row := range monthlyRevenueAndCashflow {
		for j, cellValue := range row {
			cellName, err := excelize.CoordinatesToCellName(j+1, i+dataStartRow)
			if err != nil {
				return fmt.Errorf("failed to get cell name: %v", err)
			}
			if err := f.SetCellValue(sheet2Name, cellName, cellValue); err != nil {
				return fmt.Errorf("failed to set cell value: %v", err)
			}
		}
	}

	// Set active sheet to Sheet1.
	idx, err := f.GetSheetIndex(sheet1)
	if err != nil {
		return fmt.Errorf("failed to set active sheet: %v", err)
	}
	f.SetActiveSheet(idx)

	// Save the Excel file.
	if err := f.SaveAs(xlsxOutputPath); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	fmt.Println("Combined data written to XLSX file:", xlsxOutputPath)
	return nil
}

func CombineSuccessfulStocks(stocks []string, downloadDir, finalOutputDir string) {
	for _, stock := range stocks {
		stockDir := filepath.Join(downloadDir, stock)
		finalOutput := filepath.Join(finalOutputDir, stock+".xlsx")
		if err := combineAllCSVInFolderToXLSX(stockDir, finalOutput); err != nil {
			log.Printf("Error combining stock %s: %v", stock, err)
			continue
		}
		log.Printf("Successfully combined data for stock %s", stock)
	}
}

func addPERTitle(f *excelize.File, sheet string) error {
	if err := f.MergeCell(sheet, "A1", "L1"); err != nil {
		return err
	}
	return f.SetCellValue(sheet, "A1", "PER")
}

func addPERHeader(f *excelize.File, sheet string) error {
	// Row 2: Set single-cell headers (with vertical merge)
	perHeaders := []struct {
		col  string
		text string
	}{
		{"A", "交易週別"},
		{"B", "收盤價"},
		{"C", "漲跌價"},
		{"D", "漲跌幅"},
		{"E", "河流圖 EPS(元)"},
		{"F", "目前 PER (倍)"},
	}
	for _, h := range perHeaders {
		cell := h.col + "2"
		if err := f.SetCellValue(sheet, cell, h.text); err != nil {
			return err
		}
		if err := f.MergeCell(sheet, cell, h.col+"3"); err != nil {
			return err
		}
	}
	return nil
}

func addStockTitle(f *excelize.File, sheet string) error {
	if err := f.MergeCell(sheet, "H1", "AH1"); err != nil {
		return err
	}
	return f.SetCellValue(sheet, "H1", "Stock Data")
}

func addStockHeader(f *excelize.File, sheet string) error {
	// Row 2: Single headers with vertical merge (columns H–P).
	weeklySingle := []struct {
		col  string
		text string
	}{
		{"H", "交易週別"},
		{"I", "交易日數"},
		{"J", "開盤"},
		{"K", "最高"},
		{"L", "最低"},
		{"M", "收盤"},
		{"N", "漲跌"},
		{"O", "漲跌(%)"},
		{"P", "振幅(%)"},
	}
	for _, h := range weeklySingle {
		cell := h.col + "2"
		if err := f.SetCellValue(sheet, cell, h.text); err != nil {
			return err
		}
		if err := f.MergeCell(sheet, cell, h.col+"3"); err != nil {
			return err
		}
	}

	// Row 2: Group header "成交張數" spanning columns Q to R.
	if err := f.SetCellValue(sheet, "Q2", "成交張數"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "Q2", "R2"); err != nil {
		return err
	}

	// Group header "成交金額" spanning columns S to T.
	if err := f.SetCellValue(sheet, "S2", "成交金額"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "S2", "T2"); err != nil {
		return err
	}

	// Group header "法人買賣超(千張)" spanning columns U to X.
	if err := f.SetCellValue(sheet, "U2", "法人買賣超(千張)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "U2", "X2"); err != nil {
		return err
	}

	// Row 2: Single header "外資持股(%)" in column Y (vertical merge).
	if err := f.SetCellValue(sheet, "Y2", "外資持股(%)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "Y2", "Y3"); err != nil {
		return err
	}

	// Group header "融資(千張)" spanning columns Z to AA.
	if err := f.SetCellValue(sheet, "Z2", "融資(千張)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "Z2", "AA2"); err != nil {
		return err
	}

	// Group header "融券(千張)" spanning columns AB to AC.
	if err := f.SetCellValue(sheet, "AB2", "融券(千張)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "AB2", "AC2"); err != nil {
		return err
	}

	// Row 2: Single header "券資比(%)" in column AD (vertical merge).
	if err := f.SetCellValue(sheet, "AD2", "券資比(%)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "AD2", "AD3"); err != nil {
		return err
	}

	// Row 3: Subheaders for the grouped headers.
	// For "成交張數" (columns Q–R).
	weeklySub1 := []struct {
		col  string
		text string
	}{
		{"Q", "千張"},
		{"R", "日均"},
	}
	for _, h := range weeklySub1 {
		if err := f.SetCellValue(sheet, h.col+"3", h.text); err != nil {
			return err
		}
	}

	// For "成交金額" (columns S–T).
	weeklySub2 := []struct {
		col  string
		text string
	}{
		{"S", "億元"},
		{"T", "日均"},
	}
	for _, h := range weeklySub2 {
		if err := f.SetCellValue(sheet, h.col+"3", h.text); err != nil {
			return err
		}
	}

	// For "法人買賣超(千張)" (columns U–X).
	weeklySub3 := []struct {
		col  string
		text string
	}{
		{"U", "外資"},
		{"V", "投信"},
		{"W", "自營"},
		{"X", "合計"},
	}
	for _, h := range weeklySub3 {
		if err := f.SetCellValue(sheet, h.col+"3", h.text); err != nil {
			return err
		}
	}

	// For "融資(千張)" (columns Z and AA).
	if err := f.SetCellValue(sheet, "Z3", "增減"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "AA3", "餘額"); err != nil {
		return err
	}

	// For "融券(千張)" (columns AB and AC).
	if err := f.SetCellValue(sheet, "AB3", "增減"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "AC3", "餘額"); err != nil {
		return err
	}

	return nil
}

func addRevenueTitle(f *excelize.File, sheet string) error {
	if err := f.MergeCell(sheet, "A1", "Q1"); err != nil {
		return err
	}
	return f.SetCellValue(sheet, "A1", "Revenue")
}

func addRevenueHeader(f *excelize.File, sheet string) error {
	// Row 2:
	// A2: "月別" with rowspan=3 (merge A2:A4)
	if err := f.SetCellValue(sheet, "A2", "月別"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "A2", "A4"); err != nil {
		return err
	}
	// B2:G2: "當月股價" (colspan=6)
	if err := f.SetCellValue(sheet, "B2", "當月股價"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "B2", "G2"); err != nil {
		return err
	}
	// H2:L2: "營業收入" (colspan=5)
	if err := f.SetCellValue(sheet, "H2", "營業收入"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "H2", "L2"); err != nil {
		return err
	}
	// M2:Q2: "合併營業收入" (colspan=5)
	if err := f.SetCellValue(sheet, "M2", "合併營業收入"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "M2", "Q2"); err != nil {
		return err
	}

	// Row 3: For "當月股價" details (each with vertical merge across rows 3–4).
	priceHeaders := []struct {
		col   string
		value string
	}{
		{"B", "開盤"},
		{"C", "收盤"},
		{"D", "最高"},
		{"E", "最低"},
		{"F", "漲跌(元)"},
		{"G", "漲跌(%)"},
	}
	for _, h := range priceHeaders {
		cell := fmt.Sprintf("%s3", h.col)
		if err := f.SetCellValue(sheet, cell, h.value); err != nil {
			return err
		}
		if err := f.MergeCell(sheet, cell, fmt.Sprintf("%s4", h.col)); err != nil {
			return err
		}
	}

	// Row 3: Under "營業收入":
	// H3:J3: "單月"
	if err := f.SetCellValue(sheet, "H3", "單月"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "H3", "J3"); err != nil {
		return err
	}
	// K3:L3: "累計"
	if err := f.SetCellValue(sheet, "K3", "累計"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "K3", "L3"); err != nil {
		return err
	}

	// Row 4: Detailed headers for revenue groups.
	// Under "營業收入":
	if err := f.SetCellValue(sheet, "H4", "營收(億)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "I4", "月增(%)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "J4", "年增(%)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "K4", "營收(億)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "L4", "年增(%)"); err != nil {
		return err
	}
	// Under "合併營業收入":
	if err := f.SetCellValue(sheet, "M4", "營收(億)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "N4", "月增(%)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "O4", "年增(%)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "P4", "營收(億)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "Q4", "年增(%)"); err != nil {
		return err
	}

	return nil
}

func addCashflowTitle(f *excelize.File, sheet string) error {
	if err := f.MergeCell(sheet, "S1", "AK1"); err != nil {
		return err
	}
	return f.SetCellValue(sheet, "S1", "Cash Flow")
}

func addCashflowHeader(f *excelize.File, sheet string) error {
	// Row 2:
	// S2: "季度" with rowspan=2 (merge S2:S3)
	if err := f.SetCellValue(sheet, "S2", "季度"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "S2", "S3"); err != nil {
		return err
	}
	// T2: "平均股本(億)" with rowspan=2 (merge T2:T3)
	if err := f.SetCellValue(sheet, "T2", "平均股本(億)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "T2", "T3"); err != nil {
		return err
	}
	// U2: "財報評分" with rowspan=2 (merge U2:U3)
	if err := f.SetCellValue(sheet, "U2", "財報評分"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "U2", "U3"); err != nil {
		return err
	}
	// V2:Y2: "季度股價" (colspan=4)
	if err := f.SetCellValue(sheet, "V2", "季度股價"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "V2", "Y2"); err != nil {
		return err
	}
	// Z2:AA2: "獲利(億)" (colspan=2)
	if err := f.SetCellValue(sheet, "Z2", "獲利(億)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "Z2", "AA2"); err != nil {
		return err
	}
	// AB2:AG2: "現金流量(億)" (colspan=6)
	if err := f.SetCellValue(sheet, "AB2", "現金流量(億)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "AB2", "AG2"); err != nil {
		return err
	}
	// AH2:AI2: "現金餘額(億)" (colspan=2)
	if err := f.SetCellValue(sheet, "AH2", "現金餘額(億)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "AH2", "AI2"); err != nil {
		return err
	}
	// AJ2: "現金流量(%)" with rowspan=2 (merge AJ2:AJ3)
	if err := f.SetCellValue(sheet, "AJ2", "現金流量(%)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "AJ2", "AJ3"); err != nil {
		return err
	}
	// AK2: "稅後EPS(元)" with rowspan=2 (merge AK2:AK3)
	if err := f.SetCellValue(sheet, "AK2", "稅後EPS(元)"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "AK2", "AK3"); err != nil {
		return err
	}

	// Row 3: Subheaders for the grouped columns.
	// Under "季度股價" (columns V–Y):
	if err := f.SetCellValue(sheet, "V3", "上期收盤"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "W3", "本期收盤"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "X3", "漲跌(元)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "Y3", "漲跌(%)"); err != nil {
		return err
	}
	// Under "獲利(億)" (columns Z–AA):
	if err := f.SetCellValue(sheet, "Z3", "稅前淨利"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "AA3", "稅後淨利"); err != nil {
		return err
	}
	// Under "現金流量(億)" (columns AB–AG):
	if err := f.SetCellValue(sheet, "AB3", "營業活動"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "AC3", "投資活動"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "AD3", "融資活動"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "AE3", "其他活動"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "AF3", "淨現金流"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "AG3", "自由金流"); err != nil {
		return err
	}
	// Under "現金餘額(億)" (columns AH–AI):
	if err := f.SetCellValue(sheet, "AH3", "期初餘額"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "AI3", "期末餘額"); err != nil {
		return err
	}

	return nil
}
