package format

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Table(data [][]string) string {
	if len(data) == 0 {
		return ""
	}

	// make all cells in a column the same width
	max := make([]int, len(data[0]))
	for _, row := range data {
		for c, cell := range row {
			if len(cell) > max[c] {
				max[c] = len(cell)
			}
		}
	}
	for r, row := range data {
		for c, cell := range row {
			for s := len(cell); s < max[c]; s++ {
				data[r][c] += " "
			}
		}
	}

	// make it smaller if the result is too wide
	// only by making the widest column smaller for now
	maxWidth := findTermWidth()
	if maxWidth != 0 {
		width := len(max) - 1
		for _, m := range max {
			width += m
		}
		shortenWith := width - maxWidth
		widestColNo, widestColLen := 0, 0
		for i, m := range max {
			if m > widestColLen {
				widestColNo, widestColLen = i, m
			}
		}
		newTaskColWidth := max[widestColNo] - shortenWith
		if newTaskColWidth < 0 {
			return "table is too wide to display\n"
		}
		if newTaskColWidth < max[widestColNo] {
			for r, row := range data {
				data[r][widestColNo] = row[widestColNo][:newTaskColWidth]
			}
		}
	}

	// print the rows
	var output string
	for r, row := range data {
		if r%3 == 0 {
			output += fmt.Sprintf("%s", "\x1b[48;5;237m")
		}
		for c, col := range row {
			output += col
			if c != len(row)-1 {
				output += " "
			}
		}
		if r%3 == 0 {
			output += fmt.Sprintf("%s", "\x1b[49m")
		}
		output += "\n"
	}

	return output
}

func findTermWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0
	}

	s := string(out)
	s = strings.TrimSpace(s)
	sArr := strings.Split(s, " ")

	width, err := strconv.Atoi(sArr[1])
	if err != nil {
		return 0
	}
	return width
}
