package service

// CronJob represents a cron job registration entry
// Handler harus berupa fungsi tanpa parameter

type CronJob struct {
	Spec    string
	Handler func()
}

// GetCronJobs returns all cron jobs with injected dependencies
func GetCronJobs(dailyReportSvc *DailyReportService) []CronJob {
	return []CronJob{
		{
			Spec:    "*/5 * * * * *",
			Handler: dailyReportSvc.GenerateAndSendDailyReport,
		},
		// Tambahkan job lain di sini
	}
}
