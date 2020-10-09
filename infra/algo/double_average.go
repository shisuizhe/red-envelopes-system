package algo

import (
	"math/rand"
	"time"
)

// 二倍均值算法
func DoubleAverage(count, amount int64) int64 {
	// count 剩余数量 amount 剩余金额(单位/分)
	const MIN = int64(1) 		// 最小金额 1 分钱
	if count <= 0 {
		return 0
	}
	if count == 1 {
		return amount
	}
	max := amount - MIN * count	// 计算出最大可用金额
	avg := max / count			// 计算最大可用平均值
	avg2 := avg * 2 + MIN		// 二倍均值基础再加上最小金额，防止出现0值
	// 随机红包金额序列元素，把二倍均值作为随机的最大数
	rand.Seed(time.Now().UnixNano())
	x := rand.Int63n(avg2) + MIN
	return x
}
