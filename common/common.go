package common

import "encoding/json"

// RuleMap Dictionary of group rules, matching rules => reply content
type RuleMap map[string]string

var (
	// AllGroupRules Dictionary of rules for all groups
	AllGroupRules = make(map[int64]RuleMap)
	// AllGroupId The ids of all groups currently serving
	AllGroupId []int64
)

func (rm RuleMap) String() string {
	s, err := json.Marshal(rm)
	if err != nil {
		return ""
	}
	return string(s)
}

// Json2kvs Convert json string to rule dictionary
func Json2kvs(rulesJson string) RuleMap {
	tkvs := make(RuleMap)
	_ = json.Unmarshal([]byte(rulesJson), &tkvs)
	return tkvs
}

// AddNewGroup Add an entry for the new group in memory
func AddNewGroup(gid int64) {
	AllGroupId = append(AllGroupId, gid)
	AllGroupRules[gid] = make(RuleMap)
}
