package ruleservice

import (
	"encoding/json"
	"os"
	"sync"
)

type ruleServiceImpl struct {
	mu        sync.Mutex
	rules     map[string]Rule
	rulesFile string
}

func NewRuleService(rulesFileName string) (*ruleServiceImpl, error) {
	service := &ruleServiceImpl{
		rules:     make(map[string]Rule),
		rulesFile: rulesFileName,
	}
	if err := service.loadRules(); err != nil {
		return nil, err
	}
	return service, nil
}

// loadRules reads the JSON file and populates the rules map.
func (r *ruleServiceImpl) loadRules() error {
	data, err := os.ReadFile(r.rulesFile)
	if err != nil {
		return err
	}

	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return err
	}

	// Convert slice to map for easier lookup
	for _, rule := range rules {
		r.rules[rule.ID()] = rule
	}

	return nil
}

func (r *ruleServiceImpl) saveRules() error {
	// Convert the map back to a slice for JSON marshaling
	var rules []Rule
	for _, rule := range r.rules {
		rules = append(rules, rule)
	}

	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.rulesFile, data, 0644)
}

func (r *ruleServiceImpl) GetRules() []Rule {
	var ruleList []Rule
	for _, v := range r.rules {
		ruleList = append(ruleList, v)
	}
	return ruleList
}

func (r *ruleServiceImpl) AddOrUpdate(rule Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID()] = rule
	return r.saveRules()
}

func (r *ruleServiceImpl) Delete(rule Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, rule.ID())
	return r.saveRules()
}

func (r *ruleServiceImpl) FindRule(uriPrefix, fromIP string) Rule {
	return r.rules[uriPrefix+fromIP]
}
