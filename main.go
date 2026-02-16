package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed dist/index.html
var templateHTML string

func main() {
	coveragePath := flag.String("coverage", "", "path to cover.out file (required)")
	codebasePath := flag.String("codebase", "./", "path to the Go source code root")
	routeCoveragePath := flag.String("route-coverage", "", "path to route-coverage.rcov file (optional)")
	outputPath := flag.String("out", "coverage.html", "output HTML file path")
	flag.Parse()

	if *coveragePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	coverageBytes, err := os.ReadFile(*coveragePath)
	if err != nil {
		fatal("reading coverage file: %v", err)
	}
	coverageStr := string(coverageBytes)

	modulePrefix := detectModulePrefix(coverageStr)
	if modulePrefix == "" {
		fmt.Fprintln(os.Stderr, "Module prefix not detected; using auto path resolution")
	} else {
		fmt.Fprintf(os.Stderr, "Detected module prefix: %s\n", modulePrefix)
	}

	codeFilesMap := buildCodeFilesMap(coverageStr, modulePrefix, *codebasePath)
	fmt.Fprintf(os.Stderr, "Collected %d source files\n", len(codeFilesMap))

	coverageGz, err := gzBase64(coverageBytes)
	if err != nil {
		fatal("compressing coverage: %v", err)
	}

	codeMapJSON, err := json.Marshal(codeFilesMap)
	if err != nil {
		fatal("marshaling code files map: %v", err)
	}
	codeMapGz, err := gzBase64(codeMapJSON)
	if err != nil {
		fatal("compressing code files map: %v", err)
	}

	var routeCoverageGz string
	if *routeCoveragePath != "" {
		rcBytes, err := os.ReadFile(*routeCoveragePath)
		if err != nil {
			fatal("reading route coverage: %v", err)
		}
		routeCoverageGz, err = gzBase64(rcBytes)
		if err != nil {
			fatal("compressing route coverage: %v", err)
		}
	}

	configs := buildConfigsJS(coverageGz, codeMapGz, routeCoverageGz)
	output := strings.Replace(templateHTML, "// CONFIGURATIONS_PLACEHOLDER", configs, 1)

	if err := os.WriteFile(*outputPath, []byte(output), 0644); err != nil {
		fatal("writing output: %v", err)
	}

	stat, _ := os.Stat(*outputPath)
	fmt.Fprintf(os.Stderr, "Written %s (%.1f KB)\n", *outputPath, float64(stat.Size())/1024)
}

func detectModulePrefix(coverageContent string) string {
	filePaths := extractCoverageFilePaths(coverageContent)
	if len(filePaths) == 0 {
		return ""
	}
	return detectCommonPrefix(filePaths)
}

func extractCoverageFilePaths(coverageContent string) []string {
	lines := strings.Split(coverageContent, "\n")
	filePaths := make([]string, 0, len(lines))
	seen := make(map[string]bool, len(lines))
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}
		colonIdx := strings.LastIndex(line, ":")
		if colonIdx < 0 {
			continue
		}
		fp := line[:colonIdx]
		if seen[fp] {
			continue
		}
		seen[fp] = true
		filePaths = append(filePaths, fp)
	}
	return filePaths
}

func detectCommonPrefix(filePaths []string) string {
	if len(filePaths) == 0 {
		return ""
	}
	parts := strings.Split(filePaths[0], "/")
	prefixParts := []string{}
	for i := 0; i < len(parts)-1; i++ {
		segment := parts[i]
		allMatch := true
		for _, fp := range filePaths {
			fpParts := strings.Split(fp, "/")
			if i >= len(fpParts) || fpParts[i] != segment {
				allMatch = false
				break
			}
		}
		if !allMatch {
			break
		}
		prefixParts = append(prefixParts, segment)
	}
	return strings.Join(prefixParts, "/")
}

func buildCodeFilesMap(coverageContent, modulePrefix, codebasePath string) map[string]string {
	lines := strings.Split(coverageContent, "\n")
	needed := make(map[string]bool)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}
		colonIdx := strings.LastIndex(line, ":")
		if colonIdx < 0 {
			continue
		}
		fp := line[:colonIdx]
		needed[fp] = true
	}

	result := make(map[string]string)
	codebasePath = strings.TrimSuffix(codebasePath, "/")

	for fullPath := range needed {
		diskPath := resolveCoverageDiskPath(fullPath, modulePrefix, codebasePath)
		if diskPath == "" {
			fmt.Fprintf(os.Stderr, "Warning: could not resolve source path for %s\n", fullPath)
			continue
		}
		data, err := os.ReadFile(diskPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not read %s: %v\n", diskPath, err)
			continue
		}
		result[fullPath] = string(data)
	}
	return result
}

func resolveCoverageDiskPath(fullPath, modulePrefix, codebasePath string) string {
	normalized := filepath.FromSlash(strings.TrimSpace(fullPath))
	if normalized == "" {
		return ""
	}

	if filepath.IsAbs(normalized) && fileExists(normalized) {
		return normalized
	}

	if modulePrefix != "" {
		prefixWithSep := modulePrefix + "/"
		if strings.HasPrefix(fullPath, prefixWithSep) {
			relPath := strings.TrimPrefix(fullPath, prefixWithSep)
			candidate := filepath.Join(codebasePath, filepath.FromSlash(relPath))
			if fileExists(candidate) {
				return candidate
			}
		}
	}

	direct := filepath.Join(codebasePath, normalized)
	if fileExists(direct) {
		return direct
	}

	parts := strings.Split(filepath.ToSlash(normalized), "/")
	for i := 1; i < len(parts)-1; i++ {
		relPath := filepath.FromSlash(strings.Join(parts[i:], "/"))
		candidate := filepath.Join(codebasePath, relPath)
		if fileExists(candidate) {
			return candidate
		}
	}

	return ""
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func gzBase64(data []byte) (string, error) {
	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	if _, err := gz.Write(data); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func buildConfigsJS(coverageGz, codeMapGz, routeCoverageGz string) string {
	var sb strings.Builder
	sb.WriteString("window.configs = {\n")
	sb.WriteString(fmt.Sprintf(`  "SIXTINEAI_COVERAGE_DATA": "%s",`, coverageGz))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "SIXTINEAI_CODE_FILES_MAP": "%s"`, codeMapGz))
	if routeCoverageGz != "" {
		sb.WriteString(",\n")
		sb.WriteString(fmt.Sprintf(`  "SIXTINEAI_ROUTE_COVERAGE_DATA": "%s"`, routeCoverageGz))
	}
	sb.WriteString("\n};")
	return sb.String()
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
