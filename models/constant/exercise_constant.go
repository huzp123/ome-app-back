package constant

// ExerciseTypes 运动类型选项
var ExerciseTypes = []string{
	"跑步", "走路", "骑行", "游泳", "瑜伽",
	"健身", "篮球", "足球", "网球", "羽毛球",
	"乒乓球", "爬山", "跳舞", "滑雪", "拳击",
}

// ValidExerciseTypesMap 有效运动类型映射（用于快速验证）
var ValidExerciseTypesMap = func() map[string]bool {
	m := make(map[string]bool)
	for _, exerciseType := range ExerciseTypes {
		m[exerciseType] = true
	}
	return m
}()
