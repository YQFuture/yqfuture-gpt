package consts

// 平台类型
const (
	// Jd 京东
	Jd = 1
	// Pdd 拼多多
	Pdd = 2
	// Qn 千牛
	Qn = 3
)

// 训练状态
const (
	// Undefined 未定义
	Undefined = 0
	// Presetting 预设中
	Presetting = 1
	// PresettingComplete 预设完成
	PresettingComplete = 2
	// Training 训练中
	Training = 11
	// TrainingComplete 训练完成
	TrainingComplete = 12
)

// 训练结果
const (
	// TrainingSuccess 训练成功
	TrainingSuccess = 1
	// TrainingFail 训练失败
	TrainingFail = 2
	// TrainingSuccessPart 部分训练成功
	TrainingSuccessPart = 3
)
