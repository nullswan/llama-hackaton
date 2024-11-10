package code

import (
	"bufio"
	"strings"
)

func ParseCodeBlocks(input string) []Block {
	var blocks []Block
	scanner := bufio.NewScanner(strings.NewReader(input))
	var currentBlock Block
	inCodeBlock := false
	var selectedLang string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				if currentBlock.Language == selectedLang {
					blocks = append(blocks, currentBlock)
				}
				currentBlock = Block{}
				inCodeBlock = false
			} else {
				currentBlock.Language = strings.TrimSpace(strings.TrimPrefix(line, "```"))
				if selectedLang == "" {
					selectedLang = currentBlock.Language
				}
				inCodeBlock = true
			}
		} else if inCodeBlock {
			if currentBlock.Code == "" {
				currentBlock.Code = line
			} else {
				currentBlock.Code += "\n" + line
			}
		}
	}

	if inCodeBlock {
		if currentBlock.Language == selectedLang {
			blocks = append(blocks, currentBlock)
		}
	}

	return blocks
}
