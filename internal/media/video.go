package media

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type VideoMetadata struct {
	Width    int
	Height   int
	Duration float64
	Format   string
}

func ProbeVideo(path string) (VideoMetadata, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return VideoMetadata{}, err
	}
	ffprobePath, err := resolveBinaryPath("ffprobe")
	if err != nil {
		return VideoMetadata{}, yxerrors.Remote("failed to locate ffprobe", err.Error()).
			WithHint("请确认 ffprobe 已安装并可在 PowerShell 中通过 Get-Command ffprobe 找到。")
	}
	cmd := exec.Command(ffprobePath, "-v", "error", "-print_format", "json", "-show_streams", "-show_format", abs)
	output, err := cmd.Output()
	if err != nil {
		return VideoMetadata{}, yxerrors.Remote("failed to probe video metadata", err.Error()).
			WithHint("请确认 ffprobe 可用，且视频文件路径正确。")
	}

	var payload struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
			Duration  string `json:"duration"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
			Format   string `json:"format_name"`
		} `json:"format"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		return VideoMetadata{}, err
	}

	meta := VideoMetadata{
		Format: strings.TrimSpace(payload.Format.Format),
	}
	for _, stream := range payload.Streams {
		if stream.CodecType != "video" {
			continue
		}
		meta.Width = stream.Width
		meta.Height = stream.Height
		if meta.Duration == 0 {
			meta.Duration = parseDurationSeconds(stream.Duration)
		}
		break
	}
	if meta.Duration == 0 {
		meta.Duration = parseDurationSeconds(payload.Format.Duration)
	}
	if meta.Width <= 0 || meta.Height <= 0 || meta.Duration <= 0 {
		return VideoMetadata{}, yxerrors.Usage("video metadata is incomplete", map[string]interface{}{
			"width":    meta.Width,
			"height":   meta.Height,
			"duration": meta.Duration,
		}).WithHint("请确认视频文件可被 ffprobe 正常读取，且文件未损坏。")
	}
	return meta, nil
}

func parseDurationSeconds(value string) float64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	var seconds float64
	_, _ = fmt.Sscanf(value, "%f", &seconds)
	return seconds
}

func resolveBinaryPath(name string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "(Get-Command "+name+" | Select-Object -ExpandProperty Source)")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
