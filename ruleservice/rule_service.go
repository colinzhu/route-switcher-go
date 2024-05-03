package ruleservice

// RuleService is an interface that defines rule management operations.
type RuleService interface {
	RetrieveRules() []Rule
	AddOrUpdate(rule Rule) error
	Delete(rule Rule) error
	FindRule(uriPrefix, fromIP string) Rule
}
