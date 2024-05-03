package ruleservice

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type ruleServiceImpl struct {
	mu        sync.Mutex
	rulesMap  map[string]Rule
	rulesFile string
}

func NewRuleService(rulesFileName string) (RuleService, error) {
	service := &ruleServiceImpl{
		rulesMap:  make(map[string]Rule),
		rulesFile: rulesFileName,
	}
	if err := service.loadRules(); err != nil {
		return nil, err
	}
	return service, nil
}

// load rules.json to rulesMap
func (r *ruleServiceImpl) loadRules() error {
	data, err := os.ReadFile(r.rulesFile)
	if err != nil {
		return err
	}

	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return err
	}

	for _, rule := range rules {
		r.rulesMap[rule.ID()] = rule
	}
	return nil
}

// save rulesMap to rules.json
func (r *ruleServiceImpl) saveRules() error {
	var rules []Rule
	for _, rule := range r.rulesMap {
		rules = append(rules, rule)
	}

	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.rulesFile, data, 0644)
}

func (r *ruleServiceImpl) RetrieveRules() []Rule {
	var ruleList []Rule
	for _, v := range r.rulesMap {
		ruleList = append(ruleList, v)
	}
	return ruleList
}

func (r *ruleServiceImpl) AddOrUpdate(rule Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rulesMap[rule.ID()] = rule
	return r.saveRules()
}

func (r *ruleServiceImpl) Delete(rule Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rulesMap, rule.ID())
	return r.saveRules()
}

func (r *ruleServiceImpl) FindRule(uriPrefix, fromIP string) Rule {
	rule := r.rulesMap[uriPrefix+fromIP]
	if rule.URIprefix == uriPrefix && rule.FromIP == fromIP {
		log.Printf("find rule by uriPrefix: %s, fromIP: %s", uriPrefix, fromIP)
		return rule
	}
	rule = r.rulesMap[uriPrefix]
	if rule.URIprefix == uriPrefix {
		log.Printf("find rule by uriPrefix: %s", uriPrefix)
		return rule
	}
	return Rule{}
}
