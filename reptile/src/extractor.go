package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/axgle/mahonia"
)

// Extractor struct
type Extractor struct {
	html        string // html文档
	blockSize   int
	text        string // 除去标签后的文本
	isCharsetGB bool   // 是不是 GB编码
	textLines   []string
	blocksLen   []int
}

// handleEnconding 获取网页编码格式
func (ex *Extractor) handleEnconding() {
	ex.isCharsetGB = true

	match := `(?i:charset\s*=\s*"?([\w\d-]*)"?)`
	regb := regexp.MustCompile(match)
	charset := regb.FindAllStringSubmatch(ex.html, -1)
	if charset != nil {
		str := strings.ToLower(charset[0][1])

		if !strings.HasPrefix(str, "gb") {
			ex.isCharsetGB = false
		}
	}
}

// getDoc 获取纯文本页面
func (ex Extractor) getDoc(doc string) string {
	var (
		reDATA    = regexp.MustCompile(`(?is:<!DOCTYPE.*?>)`)                             // DOCTYPE re
		reCommed  = regexp.MustCompile(`<!--[\s\S]*?-->`)                                 // 注释 re
		reScript  = regexp.MustCompile(`(?i:<\s*script[^>]*>[\w\W]*?<\s*/\s*script\s*>)`) // script re
		reStyle   = regexp.MustCompile(`(?i:<\s*style[^>]*>[^<]*<\s*/\s*style\s*>)`)      // style re
		reTag     = regexp.MustCompile(`<[\s\S]*?>`)                                      // HTML Tag re
		reSpecial = regexp.MustCompile(`&.{1,5};|&#.{1,5};`)                              // Special charcaters re
		reWrap    = regexp.MustCompile(`\r\n|\r`)                                         // Word wrap transform
		reRedun   = regexp.MustCompile(fmt.Sprintf("\n{%s,}", ex.blockSize+1))
	)

	doc = reDATA.ReplaceAllString(doc, "")
	doc = reCommed.ReplaceAllString(doc, "")
	doc = reScript.ReplaceAllString(doc, "")
	doc = reStyle.ReplaceAllString(doc, "")
	doc = reTag.ReplaceAllString(doc, "")
	doc = reSpecial.ReplaceAllString(doc, "")
	// doc = reSpace.ReplaceAllString(doc, "")
	doc = reWrap.ReplaceAllString(doc, "\n")
	doc = reRedun.ReplaceAllString(doc, strings.Repeat("\n", ex.blockSize+1))

	return doc
}

// Split the preprocessed text into lines by '\n'
func (ex *Extractor) getTextLines(text string) {
	var reSpace = regexp.MustCompile(`\s+`) // Spaces re

	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if line != "" {
			line = reSpace.ReplaceAllString(line, "")
			ex.textLines = append(ex.textLines, line)
		}
	}
}

// calcBlockLens 计算每一块的字数
func (ex *Extractor) calcBlockLens() {
	textLineCount := len(ex.textLines)
	blockLen := 0
	blockSize := ex.blockSize
	if textLineCount < ex.blockSize {
		blockSize = textLineCount
	}

	for i := 0; i < blockSize; i++ {
		blockLen += len(ex.textLines[i])
	}

	ex.blocksLen = append(ex.blocksLen, blockLen)

	if ex.blockSize != blockSize {
		return
	}

	for i := 1; i < textLineCount-ex.blockSize; i++ {
		blockLen = ex.blocksLen[i-1] + len(ex.textLines[i-1+ex.blockSize]) - len(ex.textLines[i-1])
		ex.blocksLen = append(ex.blocksLen, blockLen)
	}

}

// GetPlainText 获取文本内容
func (ex *Extractor) GetPlainText() string {
	ex.handleEnconding()

	doc := ex.getDoc(ex.html)
	ex.getTextLines(doc)
	ex.calcBlockLens()

	var (
		i           = 0
		maxTextLen  = 0
		curTextLen  = 0
		blocksCount = len(ex.blocksLen)
		part        = ""
	)

	for i < blocksCount {
		if ex.blocksLen[i] > 0 {
			if ex.textLines[i] != "" {
				part = fmt.Sprintf("%s%s\n", part, ex.textLines[i])
				curTextLen += len(ex.textLines[i])
			}
		} else {
			curTextLen = 0
			part = ""
		}

		if curTextLen > maxTextLen {
			ex.text = part
			maxTextLen = curTextLen
		}
		i++
	}

	if ex.isCharsetGB {
		enc := mahonia.NewEncoder("GBK")
		ex.text = enc.ConvertString(ex.text)
	}

	return ex.text
}
