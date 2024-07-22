package progress

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func StopProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	cmd, ok := processes[name]
	if !ok {
		respondWithError(w, http.StatusNotFound, "Process not found")
		return
	}

	// 停止进程并更新status
	if err := cmd.Cmd.Process.Kill(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//delete(processes, name) // 更新status，表示进程已停止
	cmd.Running = false
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Process %s stopped.", name)
}
func StopAllHandler() error {
	processesMu.Lock()
	defer processesMu.Unlock()

	for name, cmd := range processes {
		if err := cmd.Cmd.Process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("Failed to stop process %s: %v", name, err)
		}
	}
	return nil
}
func StopAll(w http.ResponseWriter, r *http.Request) {
	err := StopAllHandler()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// processesMu.Lock()
	// defer processesMu.Unlock()

	// for name, cmd := range processes {
	// 	if err := cmd.Cmd.Process.Signal(os.Interrupt); err != nil {
	// 		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to stop process %s: %v", name, err))
	// 		return
	// 	}
	// }

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "All processes stopped.")
}

// Restart a specific process by stopping and then starting it.
func ReStartProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Stop the process first
	if cmd, ok := processes[name]; ok {
		if err := cmd.Cmd.Process.Kill(); err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to stop process %s before restart: %v", name, err))
			return
		}
		//delete(processes, name)
		cmd.Running = false
	} else {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Process %s not found to restart", name))
		return
	}
	// Start the process again
	err := ExecuteScript(name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to restart process %s: %v", name, err))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Process %s restarted.", name)
}
