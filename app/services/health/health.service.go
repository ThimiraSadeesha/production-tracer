package health

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/thimira/production-tracer/internal/db"
)

// startTime marks process boot, used to report uptime.
var startTime = time.Now()

// HealthCheck reports overall service health: database connectivity plus
// system memory (RAM) and disk (ROM) usage. It returns 503 when the database
// is unreachable.
func HealthCheck(c *gin.Context) {
	status := "ok"

	// --- Database ---
	dbHealth := gin.H{"status": "up"}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		status = "degraded"
		dbHealth["status"] = "down"
		dbHealth["error"] = err.Error()
	}

	// --- Memory (RAM), system-wide ---
	ram := gin.H{}
	if vm, err := mem.VirtualMemory(); err == nil {
		ram = gin.H{
			"total":        humanizeBytes(vm.Total),
			"used":         humanizeBytes(vm.Used),
			"available":    humanizeBytes(vm.Available),
			"used_percent": round2(vm.UsedPercent),
		}
	}

	// --- Disk (ROM) usage of the root filesystem ---
	rom := gin.H{}
	if du, err := disk.Usage("/"); err == nil {
		rom = gin.H{
			"total":        humanizeBytes(du.Total),
			"used":         humanizeBytes(du.Used),
			"free":         humanizeBytes(du.Free),
			"used_percent": round2(du.UsedPercent),
		}
	}

	httpStatus := http.StatusOK
	if status != "ok" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status":    status,
		"service":   "production-tracer",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"uptime":    time.Since(startTime).Round(time.Second).String(),
		"database":  dbHealth,
		"memory":    ram,
		"disk":      rom,
	})
}

// humanizeBytes renders a byte count as a human-readable string (KiB, MiB, …).
func humanizeBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// round2 rounds to two decimal places.
func round2(f float64) float64 {
	return math.Round(f*100) / 100
}
