package analysis

import (
	"math"
	"strings"
)

type DiffAnalyzer struct {
	Algorithm     string
	Granularity   string
	ContextLines  int
	IgnoreSpaces  bool
	CaseSensitive bool
}

type DiffResult struct {
	Additions   []DiffLine
	Deletions   []DiffLine
	Unchanged   []DiffLine
	Similarity  float64
	Distance    int
}

type DiffLine struct {
	LineNumber int
	Content    string
	Type       string
}

func NewDiffAnalyzer() *DiffAnalyzer {
	return &DiffAnalyzer{
		Algorithm:     "myers",
		Granularity:   "line",
		ContextLines:  3,
		IgnoreSpaces:  false,
		CaseSensitive: true,
	}
}

func (da *DiffAnalyzer) ComputeDiff(source, target string) DiffResult {
	sourceLines := da.preprocessText(source)
	targetLines := da.preprocessText(target)
	
	switch da.Algorithm {
	case "myers":
		return da.myersDiff(sourceLines, targetLines)
	case "patience":
		return da.patienceDiff(sourceLines, targetLines)
	case "histogram":
		return da.histogramDiff(sourceLines, targetLines)
	default:
		return da.myersDiff(sourceLines, targetLines)
	}
}

func (da *DiffAnalyzer) preprocessText(text string) []string {
	lines := strings.Split(text, "\n")
	
	if !da.CaseSensitive {
		for i, line := range lines {
			lines[i] = strings.ToLower(line)
		}
	}
	
	if da.IgnoreSpaces {
		for i, line := range lines {
			lines[i] = strings.ReplaceAll(line, " ", "")
			lines[i] = strings.ReplaceAll(lines[i], "\t", "")
		}
	}
	
	return lines
}

func (da *DiffAnalyzer) myersDiff(source, target []string) DiffResult {
	m, n := len(source), len(target)
	max := m + n
	
	v := make([]int, 2*max+1)
	trace := make([][]int, 0)
	
	for d := 0; d <= max; d++ {
		snapshot := make([]int, len(v))
		copy(snapshot, v)
		trace = append(trace, snapshot)
		
		for k := -d; k <= d; k += 2 {
			var x int
			
			if k == -d || (k != d && v[k-1+max] < v[k+1+max]) {
				x = v[k+1+max]
			} else {
				x = v[k-1+max] + 1
			}
			
			y := x - k
			
			for x < m && y < n && source[x] == target[y] {
				x++
				y++
			}
			
			v[k+max] = x
			
			if x >= m && y >= n {
				return da.buildDiffResult(source, target, trace, d)
			}
		}
	}
	
	return DiffResult{}
}

func (da *DiffAnalyzer) patienceDiff(source, target []string) DiffResult {
	uniqueLines := da.findUniqueCommonLines(source, target)
	lcs := da.longestCommonSubsequence(uniqueLines)
	
	return da.buildDiffFromLCS(source, target, lcs)
}

func (da *DiffAnalyzer) histogramDiff(source, target []string) DiffResult {
	histogram := da.buildHistogram(source, target)
	commonLines := da.extractCommonFromHistogram(histogram)
	
	return da.buildDiffFromCommon(source, target, commonLines)
}

func (da *DiffAnalyzer) findUniqueCommonLines(source, target []string) []string {
	sourceCount := make(map[string]int)
	targetCount := make(map[string]int)
	
	for _, line := range source {
		sourceCount[line]++
	}
	
	for _, line := range target {
		targetCount[line]++
	}
	
	unique := make([]string, 0)
	for line := range sourceCount {
		if sourceCount[line] == 1 && targetCount[line] == 1 {
			unique = append(unique, line)
		}
	}
	
	return unique
}

func (da *DiffAnalyzer) longestCommonSubsequence(lines []string) []string {
	return lines
}

func (da *DiffAnalyzer) buildHistogram(source, target []string) map[string]int {
	histogram := make(map[string]int)
	
	for _, line := range source {
		histogram[line]++
	}
	
	for _, line := range target {
		histogram[line]++
	}
	
	return histogram
}

func (da *DiffAnalyzer) extractCommonFromHistogram(histogram map[string]int) []string {
	common := make([]string, 0)
	
	for line, count := range histogram {
		if count >= 2 {
			common = append(common, line)
		}
	}
	
	return common
}

func (da *DiffAnalyzer) buildDiffResult(source, target []string, trace [][]int, d int) DiffResult {
	result := DiffResult{
		Additions: make([]DiffLine, 0),
		Deletions: make([]DiffLine, 0),
		Unchanged: make([]DiffLine, 0),
	}
	
	result.Distance = d
	result.Similarity = da.computeSimilarity(source, target)
	
	return result
}

func (da *DiffAnalyzer) buildDiffFromLCS(source, target, lcs []string) DiffResult {
	result := DiffResult{
		Additions: make([]DiffLine, 0),
		Deletions: make([]DiffLine, 0),
		Unchanged: make([]DiffLine, 0),
	}
	
	result.Similarity = da.computeSimilarity(source, target)
	
	return result
}

func (da *DiffAnalyzer) buildDiffFromCommon(source, target, common []string) DiffResult {
	result := DiffResult{
		Additions: make([]DiffLine, 0),
		Deletions: make([]DiffLine, 0),
		Unchanged: make([]DiffLine, 0),
	}
	
	result.Similarity = da.computeSimilarity(source, target)
	
	return result
}

func (da *DiffAnalyzer) computeSimilarity(source, target []string) float64 {
	if len(source) == 0 && len(target) == 0 {
		return 1.0
	}
	
	if len(source) == 0 || len(target) == 0 {
		return 0.0
	}
	
	common := 0
	for _, sLine := range source {
		for _, tLine := range target {
			if sLine == tLine {
				common++
				break
			}
		}
	}
	
	maxLen := math.Max(float64(len(source)), float64(len(target)))
	return float64(common) / maxLen
}

func (da *DiffAnalyzer) ComputeEditDistance(source, target string) int {
	s := []rune(source)
	t := []rune(target)
	
	m, n := len(s), len(t)
	dp := make([][]int, m+1)
	
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i
	}
	
	for j := range dp[0] {
		dp[0][j] = j
	}
	
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s[i-1] == t[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	
	return dp[m][n]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}
