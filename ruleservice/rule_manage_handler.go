package ruleservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RuleManageHandler struct {
	ruleService RuleService
}

func NewRuleManageHandler(rs RuleService) *RuleManageHandler {
	return &RuleManageHandler{ruleService: rs}
}

func (h *RuleManageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.retrieveRules(w)
	case http.MethodPost:
		h.addOrUpdateOneRule(w, r)
	case http.MethodDelete:
		h.deleteOneRule(w, r)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func (h *RuleManageHandler) retrieveRules(w http.ResponseWriter) {
	rules := h.ruleService.RetrieveRules()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rules); err != nil {
		handleErr(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *RuleManageHandler) addOrUpdateOneRule(w http.ResponseWriter, r *http.Request) {
	var rule Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		handleErr(w, err, http.StatusBadRequest)
		return
	}

	if err := h.ruleService.AddOrUpdate(rule); err != nil {
		handleErr(w, err, http.StatusInternalServerError)
		return
	}
	retrievedRules := h.ruleService.RetrieveRules()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(retrievedRules); err != nil {
		handleErr(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *RuleManageHandler) deleteOneRule(w http.ResponseWriter, r *http.Request) {
	var rule Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		handleErr(w, err, http.StatusBadRequest)
		return
	}

	if err := h.ruleService.Delete(rule); err != nil {
		handleErr(w, err, http.StatusInternalServerError)
		return
	}

	retrievedRules := h.ruleService.RetrieveRules()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(retrievedRules); err != nil {
		handleErr(w, err, http.StatusInternalServerError)
		return
	}
}

func handleErr(w http.ResponseWriter, err error, status int) {
	log.Printf("Error: %v", err)
	http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), status)
	return
}
