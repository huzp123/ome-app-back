package constant

// TimeContext 时间上下文常量
const (
	TimeContextNow   = "now"
	TimeContextToday = "today"
)

// MoodLevel 情绪等级常量
const (
	MoodLevelMin = 1
	MoodLevelMax = 7
)

// TimeContexts 所有时间上下文选项
var TimeContexts = []string{
	TimeContextNow,
	TimeContextToday,
}

// MoodLevelDescriptions 情绪等级描述
var MoodLevelDescriptions = map[string]string{
	"1": "非常愉快",
	"2": "愉快",
	"3": "有点愉快",
	"4": "不悲不喜",
	"5": "有点不愉快",
	"6": "不愉快",
	"7": "非常不愉快",
}

// Influences 影响因素选项
var Influences = []string{
	// ——个人状态——
	"健康", "健身", "饮食", "睡眠", "自我照顾", "生理周期", "身体状况",
	// ——兴趣与成长——
	"爱好", "娱乐", "音乐", "阅读", "艺术", "学习", "成长计划",
	// ——身份与内在——
	"身份", "心灵", "信仰", "价值观", "未来计划",
	// ——关系网络——
	"社群", "家人", "朋友", "伴侣", "约会", "宠物",
	// ——生活事务——
	"家务", "工作", "教育", "财务", "金钱", "交通",
	// ——外部环境——
	"旅行", "天气", "季节", "节日", "时事", "社会事件", "文化", "语言", "噪音", "环境",
	// ——媒体信息——
	"社交媒体", "新闻", "体育赛事",
}

// CommonMoodTags 常见情绪标签
var CommonMoodTags = []string{
	// ——积极——
	"平静", "开心", "兴奋", "满足", "放松", "惊喜", "感激", "希望", "欣慰", "骄傲", "心安", "被理解", "庆祝",
	// ——中性或复杂——
	"无聊", "思念", "怀旧", "释然", "困惑", "紧张",
	// ——消极——
	"烦躁", "焦虑", "沮丧", "愤怒", "疲惫", "失望", "孤独", "压抑", "恐惧", "悲伤", "羞愧", "尴尬", "被忽视",
}


// ValidInfluencesMap 有效影响因素映射（用于快速验证）
var ValidInfluencesMap = func() map[string]bool {
	m := make(map[string]bool)
	for _, influence := range Influences {
		m[influence] = true
	}
	return m
}()
