package office

import (
	"regexp"
	"strings"
)

var minChunkSize = 100

func SplitText(text string, ChunkSize int) []string {
	splits := make([]string, 0)
	Texts := splitText(text, ChunkSize)
	splits = append(splits, Texts...)
	return splits

}

// SplitText 是一个将给定文本根据指定的最大长度分割成多个字符串片段的函数。
// text 是需要分割的原始文本。
// maxLength 是可选参数，指定每个分割片段的最大长度。如果未提供，将使用默认值 50。
// 返回值是分割后的字符串片段数组。
func splitText(text string, maxLength ...int) []string {
	defaultMaxLength := 256 // 默认的最大长度值为 256

	// 检查是否提供了 maxLength 参数，若未提供，则使用默认值
	if len(maxLength) == 0 {
		maxLength = append(maxLength, defaultMaxLength)
	}

	// 调用内部函数进行实际的文本分割操作，传入指定的最大长度值
	return splitTextInternal(text, maxLength[0])
}

// splitTextInternal 将给定的文本根据指定的最大长度拆分成多个字符串。
// 文本会被处理，以便在拆分时尽可能保持句子的完整性和自然性。
//
// 参数:
//
//	text string - 需要拆分的原始文本。
//	maxLength int - 每个拆分后字符串的最大长度。
//
// 返回值:
//
//	[]string - 拆分后的字符串数组。
func splitTextInternal(text string, maxLength int) []string {

	// 处理文本，替换多个换行符，压缩空格，并移除多余的换行符
	if strings.Contains(text, "\n") {
		text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n")
		text = regexp.MustCompile(`\s`).ReplaceAllString(text, " ")
		text = strings.ReplaceAll(text, "\n\n", "")
	}

	// 为标点符号添加换行符，以改善文本的拆分效果
	text = addNewlinesForPunctuation(text)
	text = addNewlinesForEllipsis(text)
	text = addNewlinesForQuestionMarksAndPeriods(text)
	text = strings.TrimSuffix(text, "\n")

	// 将处理后的文本按空格拆分成句子
	sentences := strings.Fields(text)
	// 用于存储最终的句子数组
	finalSentences := make([]string, 0)

	for i, s := range sentences {
		// 如果句子长度超过最大长度，则进行进一步的拆分
		if len(s) > maxLength {
			// 首先按标点符号拆分句子
			punctuatedSentences := splitByPunctuation(s)
			for _, p := range punctuatedSentences {
				// 如果拆分后的部分仍然超过最大长度，则按多个空格拆分
				if len(p) > maxLength {
					parts := splitByMultipleSpaces(p)
					// 对于每个仍超过最大长度的部分，进一步按引号拆分
					for _, part := range parts {
						if len(part) > maxLength {
							quotedParts := splitByQuotes(part)
							// 将拆分得到的部分插入到原始句子列表的适当位置
							sentences = appendBefore(i, quotedParts, sentences)
							break
						}
					}
					// 将按标点符号拆分的部分插入到原始句子列表的适当位置
					sentences = appendBefore(i, punctuatedSentences, sentences)
				}
			}
		} else {
			// 如果句子长度小于10个字符，尝试与前一个句子合并
			if len(s) < minChunkSize {
				// 检查前一个句子，如果它们的长度之和小于 maxLength，则合并它们
				if i > 0 && len(finalSentences) > 0 {
					prevSentence := finalSentences[len(finalSentences)-1]
					if len(prevSentence)+len(s) <= maxLength {
						finalSentences[len(finalSentences)-1] += " " + s
						continue
					}
				}
				// 如果无法合并，直接添加到结果中
				finalSentences = append(finalSentences, s)
			}
		}
	}
	return finalSentences
}

// addNewlinesForPunctuation 在文本的标点符号后添加换行符。
// 该函数特别针对中文文本设计，适用于分隔句子以改善可读性。
// 参数：
//
//	text string - 需要处理的原始文本。
//
// 返回值：
//
//	string - 经过处理，标点符号后添加了换行符的文本。
func addNewlinesForPunctuation(text string) string {
	// 使用正则表达式匹配句尾标点符号后紧接着的非句尾标点字符，并在标点符号后添加换行符。
	return regexp.MustCompile(`([;；.!?。！？\?])([^”’])`).ReplaceAllString(text, "$1\n$2")
}

// addNewlinesForEllipsis 是一个函数，用于在文本中的省略号后添加换行符。
// 这个函数特别适用于处理文本，以改善其可读性。
// 参数：
//
//	text string - 需要处理的原始文本。
//
// 返回值：
//
//	string - 经过处理，省略号后添加了换行符的文本。
func addNewlinesForEllipsis(text string) string {
	// 使用正则表达式匹配文本中的省略号，并在之后添加换行符。
	return regexp.MustCompile(`(\.{6})([^"’”」』])`).ReplaceAllString(text, "$1\n$2")
}

// addNewlinesForQuestionMarksAndPeriods 为文本中的问号和句号添加换行符。
// 该函数使用正则表达式寻找以问号、句号或其他标点符号结尾的句子，并在这些句子后添加换行符，
// 使得每个句子都独占一行。这有助于提高文本的可读性。
// 参数：
//
//	text string - 需要处理的原始文本。
//
// 返回值：
//
//	string - 经过处理，问号和句号后添加了换行符的文本。
func addNewlinesForQuestionMarksAndPeriods(text string) string {
	// 使用正则表达式匹配以特定标点符号结尾的句子，并在这些句子后添加换行符。
	return regexp.MustCompile(`([;；!?。！？\?]["’”」』]{0,2})([^;；!?，。！？\?])`).ReplaceAllString(text, "$1\n$2")
}

// splitByPunctuation 函数根据标点符号将句子分割成多个部分。
// 参数：
//
//	sentence - 需要分割的原始句子字符串。
//
// 返回值：
//
//	[]string - 分割后的字符串数组。
func splitByPunctuation(sentence string) []string {
	// 使用正则表达式匹配并替换句子中的特定标点符号，以换行符分隔符号和其后的文字。
	return strings.Fields(regexp.MustCompile(`([,，.]["’”」』]{0,2})([^,，.])`).ReplaceAllString(sentence, "$1\n$2"))
}

// splitByMultipleSpaces 函数根据多个连续空格或换行符分割输入字符串，并返回一个字符串切片。
// 参数：
//
//	part - 需要分割的原始字符串。
//
// 返回值：
//
//	分割后的字符串切片。
func splitByMultipleSpaces(part string) []string {
	// 使用正则表达式匹配并替换多个空格或换行符，同时保留这些分隔符与非空字符之间的边界。
	return strings.Fields(regexp.MustCompile(`([\n]{1,}| {2,}["’”」』]{0,2})([^\s])`).ReplaceAllString(part, "$1\n$2"))
}

// splitByQuotes 根据引号分割字符串。
// 该函数使用正则表达式寻找被引号包围的单词，并将这些单词与其它内容分割成多个字符串。
//
// 参数:
//
//	part string - 需要分割的原始字符串。
//
// 返回值:
//
//	[]string - 分割后的字符串数组。
func splitByQuotes(part string) []string {
	// 使用正则表达式替换引号包围的单词，为其前后添加换行符，然后以换行符为分隔符分割字符串
	return strings.Fields(regexp.MustCompile(`( ["’”」』]{0,2})([^ ])`).ReplaceAllString(part, "$1\n$2"))
}

// appendBefore 在指定索引处插入新的句子数组，并返回更新后的句子数组。
// index: 插入位置的索引。
// newSentences: 要插入的新句子数组。
// sentences: 原始句子数组。
// 返回值: 更新后的句子数组。
func appendBefore(index int, newSentences []string, sentences []string) []string {
	// 将原数组分为两部分：index之前的部分和index及之后的部分。
	// 然后将新句子数组插入到index之前的部分之后，最后合并所有部分。
	return append(sentences[:index], append(newSentences, sentences[index:]...)...)
}

// SplitText2 根据最大块大小和指定的分隔符分割纯文本文件内容
func SplitText2(content string, maxChunkSize int, splitChars ...rune) []string {
	defaultSplitChars := []rune{',', '.', '\n', '！', '。', '；'}

	var chunks []string
	currentChunk := ""

	if len(splitChars) == 0 {
		splitChars = defaultSplitChars
	}

	// 按照最大块大小和分隔符分割文本
	for _, char := range content {
		if len(currentChunk)+1 > maxChunkSize && (char == ',' || char == '.') {
			chunks = append(chunks, currentChunk)
			currentChunk = ""
		} else {
			currentChunk += string(char)
		}

		// 如果字符是分隔符，不管长度，都创建一个新的块
		if contains(splitChars, char) {
			chunks = append(chunks, currentChunk)
			currentChunk = ""
		}
	}

	// 添加最后一个块
	if currentChunk != "" {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

// contains 检查 rune 切片是否包含指定的字符
func contains(rs []rune, r rune) bool {
	for _, rr := range rs {
		if rr == r {
			return true
		}
	}
	return false
}
