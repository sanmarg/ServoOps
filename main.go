package main

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type svc struct {
	PID         int    `json:"pid"`
	Name        string `json:"name"`
	User        string `json:"user"`
	MemoryUsage string `json:"memory_usage"`
	CPUUsage    string `json:"cpu_usage"`
	StartTime   string `json:"start_time"`
}

var svcs []svc

// execute commands on os
func executeCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("failed to execute command")
	}
	return string(output), nil
}

func makestruct() {
	output, err := executeCommand("ps axo pid,comm,user,%mem,%cpu,lstart")
	if err != nil {
		fmt.Println("Failed to execute command:", err)
		return
	}

	lines := strings.Split(output, "\n")
	var newSvcs []svc

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Printf("Failed to parse PID: %s\n", fields[0])
			continue
		}
		name := fields[1]
		user := fields[2]
		memoryUsage := fields[3]
		cpuUsage := fields[4]
		startTime := strings.Join(fields[5:], " ")
		newSvcs = append(newSvcs, svc{PID: pid, Name: name, User: user, MemoryUsage: memoryUsage, CPUUsage: cpuUsage, StartTime: startTime})
	}

	svcs = newSvcs
	fmt.Println("Services updated")
}

func getServices(c *gin.Context) {
	c.JSON(http.StatusOK, svcs)
}

func svckill(id int) {
	fmt.Println(executeCommand("kill -9 " + strconv.Itoa(id)))
}

func killsvcById(c *gin.Context) {
	id := c.Param("id")
	for _, svc := range svcs {
		if strconv.Itoa(svc.PID) == id {
			c.JSON(http.StatusOK, svc)
			svckill(svc.PID)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "service not found"})
}

func startService(c *gin.Context) {
	serviceName := c.Param("name")
	output, err := executeCommand("systemctl start " + serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to start service", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "service started", "output": output})
}

func stopService(c *gin.Context) {
	serviceName := c.Param("name")
	output, err := executeCommand("systemctl stop " + serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to stop service", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "service stopped", "output": output})
}

func authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "c2FubWFyZwo=" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func logMiddleware(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	duration := time.Since(startTime)
	fmt.Printf("%s %s %s %d %s\n", c.Request.Method, c.Request.RequestURI, c.ClientIP(), c.Writer.Status(), duration)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Web API"})
}

func main() {
	fmt.Println("Web API Project")
	makestruct()
	router := gin.Default()
	router.Use(logMiddleware)
	router.Use(authMiddleware)

	router.GET("/", welcome)
	router.GET("/services", getServices)
	router.DELETE("/services/:id", killsvcById)
	router.POST("/services/start/:name", startService)
	router.POST("/services/stop/:name", stopService)
	router.GET("/health", healthCheck)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			makestruct()
		}
	}()

	router.Run(":8082")
}
