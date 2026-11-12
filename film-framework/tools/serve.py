#!/usr/bin/env python3
"""静态文件服务，带 Range 支持 —— 审片模式（?audit=1）需要它。

python3 -m http.server 不响应 Range 请求：浏览器拿不到可跳转范围，
对媒体元素 currentTime 的赋值会被静默忽略，进度条跳转就失效。
这个服务按 RFC 7233 回 206 分片，观看、审片、录制、抓帧都用它。

    python3 tools/serve.py <目录> [端口]     # 端口默认 8099
"""
import os
import re
import sys
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer


class RangeHandler(SimpleHTTPRequestHandler):
    range_re = re.compile(r"bytes=(\d*)-(\d*)$")

    def send_head(self):
        rng = self.headers.get("Range")
        if rng is None:
            return super().send_head()
        f, start, end, size = self._range_target(rng)
        if f is None:
            return None
        return self._send_partial(f, start, end, size)

    def _range_target(self, rng):
        """打开文件并算出分片边界；范围非法时回 416，其余交回原逻辑。"""
        path = self.translate_path(self.path)
        if os.path.isdir(path):
            super().send_head()  # 目录走原逻辑（列表/重定向），已写完响应
            return None, 0, 0, 0
        try:
            f = open(path, "rb")
        except OSError:
            self.send_error(404, "File not found")
            return None, 0, 0, 0
        try:
            size = os.fstat(f.fileno()).st_size
            m = self.range_re.match(rng.strip())
            if not m or size == 0:
                f.close()
                self.send_error(416, "Requested Range Not Satisfiable")
                return None, 0, 0, size
            lo_s, hi_s = m.groups()
            if lo_s == "":
                # 后缀式 "bytes=-N"：取末尾 N 字节
                lo, hi = max(size - int(hi_s), 0), size - 1
            else:
                # "bytes=N-" 到尾，"bytes=N-M" 到 M
                lo = int(lo_s)
                hi = size - 1 if hi_s == "" else min(int(hi_s), size - 1)
            if lo > hi or lo >= size:
                f.close()
                self.send_error(416, "Requested Range Not Satisfiable")
                return None, 0, 0, size
            f.seek(lo)
            return f, lo, hi, size
        except Exception:
            f.close()
            raise

    def _send_partial(self, f, lo, hi, size):
        try:
            self.send_response(206)
            self.send_header("Content-Type", self.guess_type(self.path))
            # 总长必须真实：Chrome 靠它建立可跳转范围，给 * 会被判为不可 seek
            self.send_header("Content-Range", f"bytes {lo}-{hi}/{size}")
            self.send_header("Content-Length", str(hi - lo + 1))
            self.send_header("Accept-Ranges", "bytes")
            self.end_headers()
            remaining = hi - lo + 1
            while remaining > 0:
                chunk = f.read(min(64 * 1024, remaining))
                if not chunk:
                    break
                self.wfile.write(chunk)
                remaining -= len(chunk)
        except (BrokenPipeError, ConnectionResetError):
            pass  # 播放器掐掉连接是常态
        finally:
            f.close()
        return None

    def end_headers(self):
        self.send_header("Accept-Ranges", "bytes")
        super().end_headers()


def main():
    directory = sys.argv[1] if len(sys.argv) > 1 else "."
    port = int(sys.argv[2]) if len(sys.argv) > 2 else 8099
    handler = lambda *a, **kw: RangeHandler(*a, directory=directory, **kw)
    server = ThreadingHTTPServer(("127.0.0.1", port), handler)
    print(f"serving {directory} at http://127.0.0.1:{port} (Range OK)")
    server.serve_forever()


if __name__ == "__main__":
    main()
