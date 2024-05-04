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

func (r *Rule) ID() string {
	return r.URIprefix + r.FromIP
}

func (r Rule) String() string {
	return r.URIprefix + "," + r.FromIP
}
