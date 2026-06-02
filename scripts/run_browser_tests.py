#!/usr/bin/env python3
from __future__ import annotations

import asyncio
import functools
import http.server
import socketserver
import threading
from pathlib import Path

from playwright.async_api import async_playwright


class QuietHandler(http.server.SimpleHTTPRequestHandler):
    def log_message(self, format: str, *args) -> None:  # noqa: A003
        return


def start_server(root: Path) -> tuple[socketserver.TCPServer, threading.Thread]:
    handler = functools.partial(QuietHandler, directory=str(root))
    server = socketserver.TCPServer(("127.0.0.1", 0), handler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    return server, thread


async def run_browser_test(base_url: str) -> None:
    async with async_playwright() as p:
      browser = await p.chromium.launch()
      page = await browser.new_page()
      await page.goto(f"{base_url}/tests/harness.html", wait_until="networkidle")
      await page.wait_for_selector("#summary[data-status='pass'], #summary[data-status='fail']", timeout=15000)
      status = await page.get_attribute("#summary", "data-status")
      summary = await page.text_content("#summary")
      if status != "pass":
          details = await page.locator("#results").all_inner_texts()
          raise RuntimeError(f"Browser tests failed: {summary}\n" + "\n".join(details))
      print(summary)
      await browser.close()


def main() -> int:
    root = Path(__file__).resolve().parent.parent / "web"
    server, thread = start_server(root)
    base_url = f"http://127.0.0.1:{server.server_address[1]}"

    try:
        asyncio.run(run_browser_test(base_url))
        return 0
    finally:
        server.shutdown()
        server.server_close()
        thread.join(timeout=2)


if __name__ == "__main__":
    raise SystemExit(main())
