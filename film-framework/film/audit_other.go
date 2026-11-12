//go:build !js || !wasm

package film

import "os"

// auditMode 在 native 端用环境变量开：NSGO_FILM_AUDIT=1。
// 审片控制条在 macOS 上同样成立（AVPlayer 支持精确 seek），
// 只是入口从 URL 参数换成了环境变量。
func auditMode() bool { return os.Getenv("NSGO_FILM_AUDIT") == "1" }

// timelineMode 在 native 端：NSGO_FILM_TIMELINE=1。
func timelineMode() bool { return os.Getenv("NSGO_FILM_TIMELINE") == "1" }
