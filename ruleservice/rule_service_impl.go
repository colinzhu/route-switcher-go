package ruleservice

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
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

	// check if rules.json exists, if not, then create one
	if _, err := os.Stat(rulesFileName); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(rulesFileName, []byte("[]"), 0644); err != nil {
			return nil, err
		}
		log.Printf("%s not found, created one", rulesFileName)
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

func (r *ruleServiceImpl) FindRule(inUrlPath, fromIP string) Rule {
	var matchedRule *Rule
	for _, oneRule := range r.rulesMap { // match by URIprefix and FromIP
		if strings.HasPrefix(inUrlPath, oneRule.URIprefix) && strings.Contains(oneRule.FromIP, fromIP) {
			//log.Printf("Found one rule by inUrlPath: %s, fromIP: %s", inUrlPath, fromIP)
			matchedRule = &oneRule
			break
		}
	}
	if matchedRule == nil {
		for _, oneRule := range r.rulesMap { // match by URIpfix only
			if strings.HasPrefix(inUrlPath, oneRule.URIprefix) {
				//log.Printf("Found one rule by inUrlPath: %s", inUrlPath)
				matchedRule = &oneRule
				break
			}
		}
	}
	if matchedRule == nil {
		//log.Printf("No rules found by inUrlPath: %s, fromIP: %s", inUrlPath, fromIP)
		return Rule{}
	} else {
		return *matchedRule
	}
}
