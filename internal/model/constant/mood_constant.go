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
	"健康", "健身", "自我照顾", "爱好",
	"身份", "心灵",
	"社群", "家人", "朋友", "伴侣", "约会",
	"家务", "工作", "教育", "旅行",
	"天气", "时事", "金钱",
}

// CommonMoodTags 常见情绪标签
var CommonMoodTags = []string{
	"平静", "开心", "兴奋", "满足", "放松",
	"烦躁", "焦虑", "沮丧", "愤怒", "疲惫",
	"无聊", "紧张", "困惑", "失望", "孤独",
}

// ValidInfluencesMap 有效影响因素映射（用于快速验证）
var ValidInfluencesMap = func() map[string]bool {
	m := make(map[string]bool)
	for _, influence := range Influences {
		m[influence] = true
	}
	return m
}()
