package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/k6mil6/game/internal/service"
	"log"
	"net/http"
	"os"
)

type Decorator func(http.Handler) http.Handler

type LifeStates struct {
	service.LifeService
}

type State struct {
	Fill int `json:"fill"`
}

func New(ctx context.Context,
	lifeService service.LifeService,
) (http.Handler, error) {
	serveMux := http.NewServeMux()

	lifeState := LifeStates{
		LifeService: lifeService,
	}

	serveMux.HandleFunc("/nextstate", lifeState.nextState)
	serveMux.HandleFunc("/setstate", lifeState.setState)
	serveMux.HandleFunc("/reset", lifeState.reset)

	return serveMux, nil
}

func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

func (ls *LifeStates) nextState(w http.ResponseWriter, r *http.Request) {
	worldState := ls.LifeService.NewState()

	err := json.NewEncoder(w).Encode(worldState.Cells)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ls *LifeStates) setState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var state State
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = ls.saveFillPercentage(state.Fill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ls *LifeStates) saveFillPercentage(fill int) error {
	file, err := os.Create("state.cfg")
	if err != nil {
		return err
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = file.WriteString(fmt.Sprintf("%d%%", fill))
	if err != nil {
		return err
	}

	return nil
}

func (ls *LifeStates) reset(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("state.cfg")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fillPercentageStr := string(data)

	var fillPercentage int
	_, err = fmt.Sscanf(fillPercentageStr, "%d", &fillPercentage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{"fill": fillPercentage})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
