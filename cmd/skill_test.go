package cmd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var markdownCodeRefPattern = regexp.MustCompile("`((?:\\.\\./)+references/[^`]+)`")

func TestSkillReferencePathsExist(t *testing.T) {
	withRepoRoot(t)

	for _, skillPath := range []string{
		filepath.Join(".agents", "skills", "yixiaoer", "SKILL.md"),
	} {
		raw, err := os.ReadFile(skillPath)
		if err != nil {
			t.Fatal(err)
		}
		matches := markdownCodeRefPattern.FindAllStringSubmatch(string(raw), -1)
		if len(matches) == 0 {
			t.Fatalf("expected %s to include references/* links", skillPath)
		}
		for _, match := range matches {
			ref := strings.TrimSpace(match[1])
			target := filepath.Clean(filepath.Join(filepath.Dir(skillPath), filepath.FromSlash(ref)))
			if _, err := os.Stat(target); err != nil {
				t.Fatalf("expected referenced path to exist: skill=%s ref=%s target=%s err=%v", skillPath, ref, target, err)
			}
		}
	}
}
