package db

import (
	"tg-keyword-reply-bot/common"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Initialize gorm to use sqlite
)

var db *gorm.DB

type setting struct {
	gorm.Model
	Key   string `gorm:"unique;not null"`
	Value string
}

type rule struct {
	gorm.Model
	GroupId  int64 `gorm:"unique;not null"`
	RuleJson string
}

// Init Database initialization, including creating a new database (if it has not been created), reading and writing basic data
func Init(newToken string) (token string) {
	dbtmp, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic("failed to connect database")
	}
	db = dbtmp
	db.AutoMigrate(&setting{}, &rule{})
	var tokenSetting setting
	db.Find(&tokenSetting, "Key=?", "token")
	token = tokenSetting.Value
	if newToken != "" {
		token = newToken
		if tokenSetting.ID > 0 {
			tokenSetting.Value = newToken
			db.Model(&tokenSetting).Update(tokenSetting)
		} else {
			db.Create(&setting{
				Key:   "token",
				Value: newToken,
			})
		}
	}
	readAllGroupRules()
	return
}

// AddNewGroup Add a record to the database to record the rule for the new group
func AddNewGroup(groupId int64) {
	db.Create(&rule{
		GroupId:  groupId,
		RuleJson: "",
	})
}

// UpdateGroupRule Update group rules
func UpdateGroupRule(groupId int64, ruleJson string) {
	db.Model(&rule{}).Where("group_id=?", groupId).Update("rule_json", ruleJson)
}

func readAllGroupRules() {
	var allGroupRules []rule
	db.Find(&allGroupRules)
	for _, rule := range allGroupRules {
		ruleStruct := common.Json2kvs(rule.RuleJson)
		common.AllGroupRules[rule.GroupId] = ruleStruct
		common.AllGroupId = append(common.AllGroupId, rule.GroupId)
	}
}
