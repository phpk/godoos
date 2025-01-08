package vector

import (
	"errors"
	"fmt"
	"math"
)

const isNormalizedPrecisionTolerance = 1e-6

// cosineSimilarity 计算两个向量的余弦相似度。
// 向量在计算前会被归一化。
// 结果值表示相似度，值越高表示向量越相似。
func cosineSimilarity(a, b []float32) (float32, error) {
	// 向量必须具有相同的长度
	if len(a) != len(b) {
		return 0, errors.New("向量必须具有相同的长度")
	}

	// 归一化向量
	aNorm := normalizeVector(a)
	bNorm := normalizeVector(b)

	// 计算点积
	dotProduct, err := dotProduct(aNorm, bNorm)
	if err != nil {
		return 0, fmt.Errorf("无法计算点积: %w", err)
	}

	return dotProduct, nil
}

// dotProduct 计算两个向量的点积。
// 对于归一化的向量，点积等同于余弦相似度。
// 结果值表示相似度，值越高表示向量越相似。
func dotProduct(a, b []float32) (float32, error) {
	// 向量必须具有相同的长度
	if len(a) != len(b) {
		return 0, errors.New("向量必须具有相同的长度")
	}

	var dotProduct float32
	for i := range a {
		dotProduct += a[i] * b[i]
	}

	return dotProduct, nil
}

// normalizeVector 归一化一个浮点数向量。
// 归一化是指将向量的每个分量除以向量的模（长度），使得归一化后的向量长度为 1。
func normalizeVector(v []float32) []float32 {
	var norm float64
	for _, val := range v {
		norm += float64(val * val)
	}
	if norm == 0 {
		return v // 避免除以零的情况
	}
	norm = math.Sqrt(norm)

	res := make([]float32, len(v))
	for i, val := range v {
		res[i] = float32(float64(val) / norm)
	}

	return res
}

// isNormalized 检查向量是否已经归一化。
// 如果向量的模接近 1，则认为它是归一化的。
func isNormalized(v []float32) bool {
	var sqSum float64
	for _, val := range v {
		sqSum += float64(val) * float64(val)
	}
	magnitude := math.Sqrt(sqSum)
	return math.Abs(magnitude-1) < isNormalizedPrecisionTolerance
}
