package store

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/process"
)

func KillByPid(pid int) error {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return fmt.Errorf("failed to create process object: %w", err)
	}

	// 杀死进程及其子进程
	if err := p.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}
	return nil
}
func StopCmd(name string) error {
	cmd, ok := processes[name]
	if !ok {
		return fmt.Errorf("Process not found")
	}
	processes[name].Running = false
	err := KillByPid(cmd.Cmd.Process.Pid)
	if err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}
	return nil
}
func StopProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := StopCmd(name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Process %s stopped.", name)
}
func StopAllHandler() error {
	processesMu.Lock()
	defer processesMu.Unlock()

	for name, cmd := range processes {
		if err := KillByPid(cmd.Cmd.Process.Pid); err != nil {
			return fmt.Errorf("failed to stop process %s: %v", name, err)
		}
		processes[name].Running = false
	}
	return nil
}
func StopAll(w http.ResponseWriter, r *http.Request) {
	err := StopAllHandler()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "All processes stopped.")
}

// Restart a specific process by stopping and then starting it.
func ReStartProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Stop the process first
	if cmd, ok := processes[name]; ok {
		if err := KillByPid(cmd.Cmd.Process.Pid); err != nil {
			// respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to stop process %s before restart: %v", name, err))
			// return
			log.Printf("Failed to stop process %s before restart: %v", name, err)
		}
		// processesMu.Lock()
		// defer processesMu.Unlock()
		processes[name].Running = false
	} else {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Process %s not found to restart", name))
		return
	}
	time.Sleep(time.Second * 2)
	// Start the process again
	err := ExecuteScript(name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to restart process %s: %v", name, err))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Process %s restarted.", name)
}
