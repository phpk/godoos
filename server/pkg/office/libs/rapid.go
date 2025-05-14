/*
 * GodoAI - A software focused on localizing AI applications
 * Copyright (C) 2024 https://godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package libs

import (
	"bytes"
	_ "embed" // Needed for go:embed
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

/*
*
./RapidOcrOnnx --models models \
--det ch_PP-OCRv4_det_infer-v7.onnx \
--rec ch_PP-OCRv4_rec_infer-v7.onnx \
--cls ch_ppocr_mobile_v2.0_cls_infer.onnx \
--keys ppocr_keys_v1.txt \
--image $TARGET_IMG \
--numThread $NUM_THREADS \
--padding 50 \
--maxSideLen 1024 \
--boxScoreThresh 0.5 \
--boxThresh 0.3 \
--unClipRatio 1.6 \
--doAngle 1 \
--mostAngle 1 \
--GPU $GPU_INDEX
*/
func RunRapid(imagePaths []string) (string, error) {

	results := make([]string, 0, len(imagePaths))
	//log.Printf("the image paths are- %v\n", imagePaths)
	runFile, err := getRapidDir()
	if err != nil {
		return "", err
	}

	modelDir, err := getRapidModelDir()
	if err != nil {
		return "", err
	}
	for _, imagePath := range imagePaths {
		//log.Printf("the image path is- %v\n", imagePath)
		res, err := ConvertImage(runFile, modelDir, imagePath)
		if err != nil {
			log.Printf("- %v\n", err)
			//return "", err
		} else {
			results = append(results, res)
		}
	}
	//res,err := ConvertImage(tmpfile, imagePath)

	finalResult := strings.Join(results, "\n")
	return finalResult, err
}

func GetImageContent(imagePath string) (string, error) {
	runFile, err := getRapidDir()
	if err != nil {
		return "", err
	}

	modelDir, err := getRapidModelDir()
	if err != nil {
		return "", err
	}
	return ConvertImage(runFile, modelDir, imagePath)
}
func ConvertImage(runFile string, modelDir string, imagePath string) (string, error) {

	// 构建命令
	cmdArgs := []string{
		runFile,
		"--models", modelDir,
		"--det", "ch_PP-OCRv4_det_infer-v7.onnx",
		"--rec", "ch_PP-OCRv4_rec_infer-v7.onnx",
		"--cls", "ch_ppocr_mobile_v2.0_cls_infer.onnx",
		"--keys", "ppocr_keys_v1.txt",
		"--image", imagePath,
		"--numThread", fmt.Sprintf("%d", runtime.NumCPU()),
		"--padding", "50",
		"--maxSideLen", "1024",
		"--boxScoreThresh", "0.5",
		"--boxThresh", "0.3",
		"--unClipRatio", "1.6",
		"--doAngle", "1",
		"--mostAngle", "1",
		"--GPU", "-1",
	}
	// 打印将要执行的命令行
	cmdStr := strings.Join(cmdArgs, " ")
	fmt.Printf("Executing command: %s\n", cmdStr)
	// 使用Command构造命令
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out    // 捕获标准输出
	cmd.Stderr = &stderr // 捕获标准错误

	// 执行命令
	err := cmd.Run()
	if err != nil {
		// 打印错误信息
		log.Printf("执行命令时出错: %v, stderr: %s", err, stderr.String())
		return "", err
	}
	// 输出命令结果
	outputStr := out.String()
	//CloseDll(tmpfile)
	resText, err := ExtractText(outputStr)
	if err != nil {
		log.Printf("提取文本时出错: %v", err)
		return "", err
	}
	return resText, err
}

func ExtractText(output string) (string, error) {
	// 查找 "=====End detect=====" 的位置
	endDetectIndex := strings.Index(output, "=====End detect=====")
	if endDetectIndex == -1 {
		return "", fmt.Errorf("expected '=====End detect=====' not found in output")
	}

	// 从 "=====End detect=====" 后面开始提取文本内容
	contentStartIndex := endDetectIndex + len("=====End detect=====\n")
	if contentStartIndex >= len(output) {
		return "", fmt.Errorf("unexpected end of output after '=====End detect====='")
	}

	// 提取从 "=====End detect=====" 到末尾的字符串，然后去除末尾的花括号
	tempContent := output[contentStartIndex:]

	// 去除开头的数字和空格，以及 "FullDetectTime(...)" 部分
	cleanedContent := strings.TrimSpace(strings.SplitN(tempContent, "\n", 2)[1])

	// 确保去除了所有不需要的内容
	//cleanedOutput := strings.TrimSuffix(cleanedContent, "}")
	// 使用正则表达式去除连续的空行

	// 去除单独的 ?、: B、</>，以及它们前后的空白字符
	re := regexp.MustCompile(`(?m)^\s*(?:\?|\s*B|</>|:)\s*$`) // (?m) 使 ^ 和 $ 匹配每一行的开始和结束
	cleanedOutput := re.ReplaceAllString(cleanedContent, "")  // 删除这些行
	// 这里的正则表达式匹配一个或多个连续的换行符
	re = regexp.MustCompile(`\n\s*\n`)
	cleanedOutput = re.ReplaceAllString(cleanedOutput, "\n") // 将连续的空行替换为单个换行符
	// 返回提取的文本内容
	return cleanedOutput, nil
}
