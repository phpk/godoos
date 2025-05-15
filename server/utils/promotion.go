package utils

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// ConvertIntToBase62 将整数转换为 62 进制的字符串（全部大写）
func ConvertIntToBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var result strings.Builder
	for num > 0 {
		remainder := num % 62
		char := base62Chars[remainder]
		if char >= 'a' && char <= 'z' {
			char -= 32 // 转换为大写字母
		}
		result.WriteString(string(char))
		num /= 62
	}

	// 反转结果字符串
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// ConvertBase62ToInt 将 62 进制字符串转换为整数（接受大写输入）
func ConvertBase62ToInt(code string) (int64, error) {
	var num int64
	for _, r := range code {
		index := strings.IndexRune(base62Chars, r)
		if index == -1 {
			return 0, fmt.Errorf("invalid character '%c' in code", r)
		}
		num = num*62 + int64(index)
	}
	return num, nil
}

// HashInt32 对 int32 类型的数字进行简单的哈希变换
func HashInt32(num int32) int64 {
	return int64(num*1933 + 271) // 简单的线性变换
}

// EncodePromotionID 将 int32 整数通过哈希变换后转换为 62 进制的字符串（全部大写）
func EncodePromotionID(num int32) string {
	hashedNum := HashInt32(num)
	return ConvertIntToBase62(hashedNum)
}

// DecodePromotionID 将 62 进制字符串通过哈希逆变换转换为 int32 整数
func DecodePromotionID(code string) (int32, error) {
	int64Num, err := ConvertBase62ToInt(code)
	if err != nil {
		return 0, err
	}

	// 逆向哈希变换
	originalNum := int32((int64Num - 271) / 1933)

	if int64(originalNum) > math.MaxInt32 {
		return 0, fmt.Errorf("integer overflow: %d is too large for int32", int64Num)
	}
	return originalNum, nil
}

var reNonLettersAndDigits = regexp.MustCompile("[^a-zA-Z0-9]")

// 将中文字符转换为拼音 并且去掉非字母 保留数字和字母
func convertChineseToPinyin(str string) string {
	var pinyinStr strings.Builder // 使用 strings.Builder 提升性能
	for _, char := range str {
		if char >= '\u4e00' && char <= '\u9fa5' { // 判断是否是中文字符
			pinyin := pinyin.LazyPinyin(string(char), pinyin.Args{Style: pinyin.Normal})
			for _, p := range pinyin {
				pinyinStr.WriteString(p) // 合并拼音，并确保没有空值
			}
		} else {
			pinyinStr.WriteRune(char) // 如果是英文或其他字符，直接保留
		}
	}
	return pinyinStr.String()
}

// 清理字符串，去掉非字母字符，并转成小写
func CleanAndConvertToLetters(str string) string {
	if str == "" {
		return "" // 防止空字符串
	}

	// 将中文转换为拼音
	str = convertChineseToPinyin(str)

	// 去掉非字母和数字字符
	str = reNonLettersAndDigits.ReplaceAllString(str, "")

	// 转换成小写字母
	return strings.ToLower(str)
}
