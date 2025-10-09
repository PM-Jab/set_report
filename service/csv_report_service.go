package service

import (
	"fmt"
	"os"
	"set-report/entity"
	"time"
)

func GenerateTargetCSVreport(listReportTarget []entity.TargetReportData) string {
	csvContent := "Symbol,Low,LowDate,High,HighDate,Target,LatestTargetDate,CurrentPrice,CurrentDate\n"
	for _, report := range listReportTarget {
		csvContent += fmt.Sprintf("%s,%.2f,%s,%.2f,%s,%.2f,%s,%.2f,%s\n",
			report.Symbol,
			report.Low,
			report.LowDate,
			report.High,
			report.HighDate,
			report.Target,
			report.LatestTargetDate,
			report.CurrentPrice,
			report.CurrentDate,
		)
	}

	// create file name with current timestamp
	filename := fmt.Sprintf("target_report_%s.csv", time.Now().Format("20060102_150405"))
	// write to file
	err := os.WriteFile(filename, []byte(csvContent), 0644)
	if err != nil {
		fmt.Println("Error writing CSV file:", err)
		return ""
	}
	fmt.Println("CSV report generated:", filename)
	return filename
}

func GenerateMonthlyAverageCSVreport(listReport []entity.MonthlyAverageStockPriceBySymbolData) string {
	if len(listReport) == 0 {
		return ""
	}
	monthContent := ""
	for _, report := range listReport[0].MonthlyStockPrice {
		monthContent += fmt.Sprintf(",%s", report.Date)
	}
	csvContent := fmt.Sprintf("Symbol,TotalPositive%s\n", monthContent)

	for _, report := range listReport {
		row := fmt.Sprintf("%s,%d", report.Symbol, report.TotalPositive)
		for _, month := range report.MonthlyStockPrice {
			row += fmt.Sprintf(",%.2f", month.Change)
		}
		csvContent += fmt.Sprintf("%s\n", row)
	}

	// create file name with current timestamp
	filename := fmt.Sprintf("monthly_average_report_%s.csv", time.Now().Format("20060102_150405"))
	// write to file
	err := os.WriteFile(filename, []byte(csvContent), 0644)
	if err != nil {
		fmt.Println("Error writing CSV file:", err)
		return ""
	}
	fmt.Println("CSV report generated:", filename)
	return filename
}
