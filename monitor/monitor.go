// Package monitor is used to collect stats of the raspberry
package monitor

import (
	"errors"
	"fmt"
	"os/exec"
	"raspberryConverter/models"
	"raspberryConverter/player"
	"regexp"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// GetStatus returns the current Status of the system
func GetStatus() (models.Status, error) {
	CPU, err := getCPULoad()
	if err != nil {
		return statusError(err)
	}
	RAM, err := getMemoryLoad()
	if err != nil {
		return statusError(err)
	}
	temp, err := getTemperature()
	if err != nil {
		return statusError(err)
	}
	playerConfig, err := player.GetConfig()
	if err != nil {
		return statusError(err)
	}
	video, err := player.GetCurrentOutputModeString()
	if err != nil {
		return statusError(err)
	}
	status := models.Status{
		CPU:         CPU,
		RAM:         RAM,
		Temperature: temp,
		URL:         playerConfig.URL,
		Video:       video,
		Status:      player.GetStatus(),
	}
	return status, nil
}

// getCPULoad returns the usage in % of the busiest core of the CPU.
func getCPULoad() (int, error) {
	t, err := time.ParseDuration("-1s")
	if err != nil {
		return 0, err
	}
	allCoresPercentUsed, err := cpu.Percent(t, true)
	if err != nil {
		return 0, err
	}
	percent := 0.0
	for _, core := range allCoresPercentUsed {
		if percent < core {
			percent = core
		}
	}
	return int(percent), nil
}

// getMemoryLoad return the used ammount of memory in %
func getMemoryLoad() (int, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return int(v.UsedPercent), nil
}

// getTemperature returns the temperature of the system in celsius degrees
func getTemperature() (float64, error) {
	output, err := exec.Command("vcgencmd", "measure_temp").CombinedOutput()
	if err != nil {
		return 0, err
	}
	tempS := regexp.MustCompile("[=']").Split(string(output), -1)[1]
	tempF, err := strconv.ParseFloat(tempS, 64)
	if err != nil {
		return 0, err
	}
	return tempF, nil
}

// statusError is a helper function that log an error an return a better formated error
func statusError(err error) (models.Status, error) {
	fmt.Println(err)
	return models.Status{}, errors.New("Error geting current status information")
}
