package ruleservice

// RuleService is an interface that defines rule management operations.
type RuleService interface {
	GetRules() map[Rule]bool
	AddOrUpdate(rule Rule)
	Delete(rule Rule)
	FindRule(uriPrefix, fromIP string) Rule
}
