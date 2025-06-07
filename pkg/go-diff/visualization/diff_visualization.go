package visualization

import (
	"fmt"
	"strings"
)

type DiffVisualizer struct {
	ColorOutput     bool
	ShowLineNumbers bool
	ContextLines    int
	UnifiedFormat   bool
}

type VisualDiff struct {
	Header   string
	Hunks    []DiffHunk
	Summary  DiffSummary
}

type DiffHunk struct {
	SourceStart int
	SourceCount int
	TargetStart int
	TargetCount int
	Lines       []VisualLine
}

type VisualLine struct {
	Type       string
	Content    string
	LineNumber int
	Color      string
}

type DiffSummary struct {
	Additions int
	Deletions int
	Total     int
}

func NewDiffVisualizer() *DiffVisualizer {
	return &DiffVisualizer{
		ColorOutput:     true,
		ShowLineNumbers: true,
		ContextLines:    3,
		UnifiedFormat:   true,
	}
}

func (dv *DiffVisualizer) Visualize(sourceFile, targetFile, sourceContent, targetContent string) VisualDiff {
	if dv.UnifiedFormat {
		return dv.generateUnifiedDiff(sourceFile, targetFile, sourceContent, targetContent)
	}
	return dv.generateSideBySideDiff(sourceFile, targetFile, sourceContent, targetContent)
}

func (dv *DiffVisualizer) generateUnifiedDiff(sourceFile, targetFile, sourceContent, targetContent string) VisualDiff {
	diff := VisualDiff{
		Header: dv.generateHeader(sourceFile, targetFile),
		Hunks:  make([]DiffHunk, 0),
	}
	
	sourceLines := strings.Split(sourceContent, "\n")
	targetLines := strings.Split(targetContent, "\n")
	
	hunk := dv.createHunk(sourceLines, targetLines)
	diff.Hunks = append(diff.Hunks, hunk)
	diff.Summary = dv.computeSummary(diff.Hunks)
	
	return diff
}

func (dv *DiffVisualizer) generateSideBySideDiff(sourceFile, targetFile, sourceContent, targetContent string) VisualDiff {
	diff := VisualDiff{
		Header: dv.generateHeader(sourceFile, targetFile),
		Hunks:  make([]DiffHunk, 0),
	}
	
	sourceLines := strings.Split(sourceContent, "\n")
	targetLines := strings.Split(targetContent, "\n")
	
	hunk := dv.createSideBySideHunk(sourceLines, targetLines)
	diff.Hunks = append(diff.Hunks, hunk)
	diff.Summary = dv.computeSummary(diff.Hunks)
	
	return diff
}

func (dv *DiffVisualizer) generateHeader(sourceFile, targetFile string) string {
	return fmt.Sprintf("--- %s\n+++ %s", sourceFile, targetFile)
}

func (dv *DiffVisualizer) createHunk(sourceLines, targetLines []string) DiffHunk {
	hunk := DiffHunk{
		SourceStart: 1,
		SourceCount: len(sourceLines),
		TargetStart: 1,
		TargetCount: len(targetLines),
		Lines:       make([]VisualLine, 0),
	}
	
	maxLen := len(sourceLines)
	if len(targetLines) > maxLen {
		maxLen = len(targetLines)
	}
	
	for i := 0; i < maxLen; i++ {
		if i < len(sourceLines) && i < len(targetLines) {
			if sourceLines[i] == targetLines[i] {
				line := VisualLine{
					Type:       "unchanged",
					Content:    " " + sourceLines[i],
					LineNumber: i + 1,
					Color:      dv.getColor("unchanged"),
				}
				hunk.Lines = append(hunk.Lines, line)
			} else {
				deleteLine := VisualLine{
					Type:       "deletion",
					Content:    "-" + sourceLines[i],
					LineNumber: i + 1,
					Color:      dv.getColor("deletion"),
				}
				addLine := VisualLine{
					Type:       "addition",
					Content:    "+" + targetLines[i],
					LineNumber: i + 1,
					Color:      dv.getColor("addition"),
				}
				hunk.Lines = append(hunk.Lines, deleteLine, addLine)
			}
		} else if i < len(sourceLines) {
			line := VisualLine{
				Type:       "deletion",
				Content:    "-" + sourceLines[i],
				LineNumber: i + 1,
				Color:      dv.getColor("deletion"),
			}
			hunk.Lines = append(hunk.Lines, line)
		} else if i < len(targetLines) {
			line := VisualLine{
				Type:       "addition",
				Content:    "+" + targetLines[i],
				LineNumber: i + 1,
				Color:      dv.getColor("addition"),
			}
			hunk.Lines = append(hunk.Lines, line)
		}
	}
	
	return hunk
}

func (dv *DiffVisualizer) createSideBySideHunk(sourceLines, targetLines []string) DiffHunk {
	hunk := DiffHunk{
		SourceStart: 1,
		SourceCount: len(sourceLines),
		TargetStart: 1,
		TargetCount: len(targetLines),
		Lines:       make([]VisualLine, 0),
	}
	
	maxLen := len(sourceLines)
	if len(targetLines) > maxLen {
		maxLen = len(targetLines)
	}
	
	for i := 0; i < maxLen; i++ {
		sourceContent := ""
		targetContent := ""
		lineType := "unchanged"
		
		if i < len(sourceLines) {
			sourceContent = sourceLines[i]
		}
		
		if i < len(targetLines) {
			targetContent = targetLines[i]
		}
		
		if sourceContent != targetContent {
			if sourceContent != "" && targetContent != "" {
				lineType = "modified"
			} else if sourceContent != "" {
				lineType = "deletion"
			} else {
				lineType = "addition"
			}
		}
		
		line := VisualLine{
			Type:       lineType,
			Content:    fmt.Sprintf("%-40s | %s", sourceContent, targetContent),
			LineNumber: i + 1,
			Color:      dv.getColor(lineType),
		}
		hunk.Lines = append(hunk.Lines, line)
	}
	
	return hunk
}

func (dv *DiffVisualizer) getColor(lineType string) string {
	if !dv.ColorOutput {
		return ""
	}
	
	switch lineType {
	case "addition":
		return "\033[32m"
	case "deletion":
		return "\033[31m"
	case "modified":
		return "\033[33m"
	case "unchanged":
		return "\033[0m"
	default:
		return "\033[0m"
	}
}

func (dv *DiffVisualizer) computeSummary(hunks []DiffHunk) DiffSummary {
	summary := DiffSummary{}
	
	for _, hunk := range hunks {
		for _, line := range hunk.Lines {
			switch line.Type {
			case "addition":
				summary.Additions++
			case "deletion":
				summary.Deletions++
			}
			summary.Total++
		}
	}
	
	return summary
}

func (dv *DiffVisualizer) RenderToString(diff VisualDiff) string {
	var builder strings.Builder
	
	builder.WriteString(diff.Header)
	builder.WriteString("\n")
	
	for _, hunk := range diff.Hunks {
		builder.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n",
			hunk.SourceStart, hunk.SourceCount,
			hunk.TargetStart, hunk.TargetCount))
		
		for _, line := range hunk.Lines {
			if dv.ShowLineNumbers {
				builder.WriteString(fmt.Sprintf("%4d: ", line.LineNumber))
			}
			
			if dv.ColorOutput {
				builder.WriteString(line.Color)
			}
			
			builder.WriteString(line.Content)
			
			if dv.ColorOutput {
				builder.WriteString("\033[0m")
			}
			
			builder.WriteString("\n")
		}
	}
	
	builder.WriteString(fmt.Sprintf("\n%d additions, %d deletions, %d total changes\n",
		diff.Summary.Additions, diff.Summary.Deletions, diff.Summary.Total))
	
	return builder.String()
}
