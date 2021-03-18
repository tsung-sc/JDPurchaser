package models

type UserInfo struct {
	HouseholdAppliance int    `json:"householdAppliance"`
	ImgURL             string `json:"imgUrl"`
	LastLoginTime      string `json:"lastLoginTime"`
	NickName           string `json:"nickName"`
	PlusStatus         string `json:"plusStatus"`
	RealName           string `json:"realName"`
	UserLevel          int    `json:"userLevel"`
	UserScoreVO        struct {
		AccountScore     int    `json:"accountScore"`
		ActivityScore    int    `json:"activityScore"`
		ConsumptionScore int    `json:"consumptionScore"`
		Default          bool   `json:"default"`
		FinanceScore     int    `json:"financeScore"`
		Pin              string `json:"pin"`
		RiskScore        int    `json:"riskScore"`
		TotalScore       int    `json:"totalScore"`
	} `json:"userScoreVO"`
}
