package ruleservice

type Rule struct {
	URIprefix     string  `json:"uriPrefix"`
	FromIP        string  `json:"fromIP"`
	TargetOptions string  `json:"targetOptions"`
	Target        string  `json:"target"`
	UpdateBy      string  `json:"updateBy"`
	UpdateTime    int64   `json:"updateTime"`
	Remark        *string `json:"remark,omitempty"` // Pointer for optional field, omitempty to ignore null value
}

func NewRule(uriPrefix, fromIP, targetOptions, target, updateBy string, updateTime int64, remark *string) *Rule {
	if remark == nil {
		remark = new(string)
	}
	return &Rule{
		URIprefix:     uriPrefix,
		FromIP:        fromIP,
		TargetOptions: targetOptions,
		Target:        target,
		UpdateBy:      updateBy,
		UpdateTime:    updateTime,
		Remark:        remark,
	}
}
func (r *Rule) ID() string {
	return r.URIprefix + r.FromIP
}

func (r Rule) String() string {
	return r.URIprefix + "|" + r.FromIP
}

// func (r *Rule) Equals(other *Rule) bool {
// 	if other == nil {
// 		return false
// 	}
// 	return r.URIprefix == other.URIprefix && r.FromIP == other.FromIP
// }

// func (r *Rule) HashCode() int {
// 	return hashcode(r.URIprefix + r.FromIP)
// }

// func hashcode(s string) int {
// 	h := 0
// 	for _, c := range s {
// 		h = 31*h + int(c)
// 	}
// 	return h
// }
