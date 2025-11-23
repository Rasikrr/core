// nolint: unused
package config

import (
	"os/exec"
	"strings"
)

// getGitTag возвращает Git tag на текущем коммите
// Возвращает пустую строку если тега нет
func getGitTag() string {
	// git describe --exact-match --tags HEAD
	// Возвращает тег только если он указывает на текущий коммит
	cmd := exec.Command("git", "describe", "--exact-match", "--tags", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getGitCommit возвращает короткий hash текущего коммита (7 символов)
func getGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "--short=7", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getGitCommitFull возвращает полный hash текущего коммита
func getGitCommitFull() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getGitBranch возвращает название текущей ветки
func getGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// buildVersion определяет версию приложения с приоритетом:
// 1. Git tag на текущем коммите (v1.0.0)
// 2. Git commit hash (c13e0f8)
// 3. "unknown" если Git недоступен
func buildVersion() string {
	// Приоритет 1: Проверяем есть ли тег на текущем коммите
	if tag := getGitTag(); tag != "" {
		return tag
	}

	// Приоритет 2: Используем короткий commit hash
	if commit := getGitCommit(); commit != "" {
		return commit
	}

	// Fallback: Git недоступен
	return "unknown"
}
