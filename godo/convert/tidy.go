/*
 * GodoOS - A lightweight cloud desktop
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
package convert

import (
	"bytes"
	"fmt"
	"io"

	"github.com/beevik/etree"
)

// TidyWithEtree 使用beevik/etree库进行简单的XML清理
func Tidy(r io.Reader) ([]byte, error) {
	// 读取并解析XML
	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("error reading and parsing XML: %w", err)
	}

	// 清理操作：例如，移除空节点
	removeEmptyNodes(doc.Root())

	// 格式化XML
	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("error writing formatted XML: %w", err)
	}

	return buf.Bytes(), nil
}

// removeEmptyNodes 遍历XML树并移除空节点
func removeEmptyNodes(node *etree.Element) {
	for i := len(node.Child) - 1; i >= 0; i-- { // 逆序遍历以安全删除
		token := node.Child[i]
		element, ok := token.(*etree.Element) // 检查是否为etree.Element类型
		if ok {
			text := element.Text() // 获取元素的文本
			if text == "" && len(element.Attr) == 0 && len(element.Child) == 0 {
				node.RemoveChildAt(i)
			} else {
				removeEmptyNodes(element) // 递归处理子节点，传入指针
			}
		}
	}
}
