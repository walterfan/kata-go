package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	gomail "gopkg.in/gomail.v2"
)

var (
	emailTo      string
	emailFrom    string
	emailSubject string
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPass     string

	serviceName string
	jobName     string
	scrapePort  int
	logger      *zap.SugaredLogger
	previousPID uint32

	cpuGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backend_cpu_percent",
		Help: "CPU usage percentage of backend process",
	})
	memGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backend_memory_bytes",
		Help: "Memory usage (RSS) of backend process in bytes",
	})
	restartCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "backend_restart_count",
		Help: "Number of detected restarts of the backend service",
	})
)

var rootCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Systemd service monitor with metrics and alerts",
	Run: func(cmd *cobra.Command, args []string) {
		initLogger()
		logger.Infow("Starting service monitor",
			"service", serviceName,
			"job", jobName,
			"port", scrapePort,
		)

		// Prometheus server
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			addr := fmt.Sprintf(":%d", scrapePort)
			logger.Infof("Serving Prometheus metrics at %s/metrics", addr)
			http.ListenAndServe(addr, nil)
		}()

		ensureCSV()

		c := cron.New()
		c.AddFunc("@every 1m", monitorService)
		c.Start()
		select {}
	},
}

func initLogger() {
	zapLogger, _ := zap.NewProduction()
	logger = zapLogger.Sugar()
}

/*
go run main.go -s tbot.service -j tbot_monitor -p 9000 \
  --email-to you@example.com \
  --email-from monitor@example.com \
  --email-subject "Alert: tbot.service restarted" \
  --smtp-host smtp.gmail.com \
  --smtp-port 587 \
  --smtp-user youruser@gmail.com \
  --smtp-pass yourpassword
*/

func main() {
	rootCmd.Flags().StringVarP(&serviceName, "service", "s", "tbot.service", "Systemd service name to monitor")
	rootCmd.Flags().StringVarP(&jobName, "job-name", "j", "tbot_monitor", "Prometheus job name")
	rootCmd.Flags().IntVarP(&scrapePort, "port", "p", 2112, "Prometheus scrape port")

	//optional email options for alert
	rootCmd.Flags().StringVar(&emailTo, "email-to", "", "Email address to send alerts to")
	rootCmd.Flags().StringVar(&emailFrom, "email-from", "", "Email address to send alerts from")
	rootCmd.Flags().StringVar(&emailSubject, "email-subject", "Service Restarted", "Subject of the email alert")
	rootCmd.Flags().StringVar(&smtpHost, "smtp-host", "smtp.example.com", "SMTP server hostname")
	rootCmd.Flags().IntVar(&smtpPort, "smtp-port", 587, "SMTP server port")
	rootCmd.Flags().StringVar(&smtpUser, "smtp-user", "", "SMTP username")
	rootCmd.Flags().StringVar(&smtpPass, "smtp-pass", "", "SMTP password")

	prometheus.MustRegister(cpuGauge, memGauge, restartCounter)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func ensureCSV() {
	const csvFile = "metrics.csv"
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		f, _ := os.Create(csvFile)
		defer f.Close()
		writer := csv.NewWriter(f)
		writer.Write([]string{"timestamp", "pid", "cpu_percent", "memory_rss_bytes", "uptime_seconds"})
		writer.Flush()
	}
}

func monitorService() {
	ctx := context.Background()
	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		logger.Errorw("Failed to connect to systemd", "error", err)
		return
	}
	defer conn.Close()

	prop, err := conn.GetUnitTypePropertyContext(context.Background(), serviceName, "Service", "MainPID")
	if err != nil {
		logger.Errorw("Failed to get PID", "error", err)
		return
	}

	pidRaw := prop.Value.Value().(uint32)
	pid := int32(pidRaw)

	if pid == 0 {
		logger.Warn("Service not running")
		return
	}

	if previousPID != 0 && pidRaw != previousPID {
		logger.Infow("Service restarted", "oldPID", previousPID, "newPID", pidRaw)
		restartCounter.Inc()
		if emailTo != "" {
			sendEmailAlert(previousPID, pidRaw)
		}

	}
	previousPID = pidRaw

	proc, err := process.NewProcess(pid)
	if err != nil {
		logger.Errorw("Failed to find process", "pid", pid, "error", err)
		return
	}

	cpu, _ := proc.CPUPercent()
	mem, _ := proc.MemoryInfo()
	uptime, _ := proc.CreateTime()
	elapsed := time.Since(time.UnixMilli(uptime))

	logger.Infow("Resource usage",
		"pid", pid,
		"cpu", cpu,
		"memory", mem.RSS,
		"uptime", elapsed.Truncate(time.Second).String(),
	)

	cpuGauge.Set(cpu)
	memGauge.Set(float64(mem.RSS))

	// Append to CSV
	f, _ := os.OpenFile("metrics.csv", os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	writer := csv.NewWriter(f)
	writer.Write([]string{
		time.Now().Format(time.RFC3339),
		strconv.Itoa(int(pid)),
		fmt.Sprintf("%.2f", cpu),
		strconv.FormatUint(mem.RSS, 10),
		fmt.Sprintf("%.0f", elapsed.Seconds()),
	})
	writer.Flush()
}

func sendEmailAlert(oldPID, newPID uint32) {
	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/plain", fmt.Sprintf("Service %s restarted\nOld PID: %d\nNew PID: %d\nTime: %s",
		serviceName, oldPID, newPID, time.Now().Format(time.RFC1123)))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	if err := d.DialAndSend(m); err != nil {
		logger.Errorw("Failed to send email", "error", err)
	} else {
		logger.Info("Sent restart email alert")
	}
}
