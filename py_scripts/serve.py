import http.server
import socketserver
import json
from urllib.parse import urlparse, parse_qs


class EchoRequestHandler(http.server.BaseHTTPRequestHandler):
    def _set_response_headers(self, content_type="application/json"):
        self.send_response(200)
        self.send_header("Content-Type", content_type)
        self.end_headers()

    def _handle_request(self):
        # Get request information
        request_info = {
            "method": self.command,
            "path": self.path,
            "headers": dict(self.headers.items()),
            "query_params": parse_qs(urlparse(self.path).query),
        }

        # Get request body if present
        content_length = int(self.headers.get("Content-Length", 0))
        if content_length > 0:
            request_body = self.rfile.read(content_length).decode("utf-8")
            try:
                request_info["body"] = json.loads(request_body)
            except json.JSONDecodeError:
                request_info["body"] = request_body

        # Return the request info as JSON
        self._set_response_headers()
        self.wfile.write(json.dumps(request_info, indent=2).encode("utf-8"))

    def do_GET(self):
        self._handle_request()

    def do_POST(self):
        self._handle_request()

    def do_PUT(self):
        self._handle_request()

    def do_DELETE(self):
        self._handle_request()


def run_server(port=8000):
    server_address = ("", port)
    httpd = socketserver.TCPServer(server_address, EchoRequestHandler)
    print(f"Server running on port {port}")
    httpd.serve_forever()


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(description="Run a simple HTTP echo server")
    parser.add_argument(
        "-p", "--port", type=int, default=8000, help="Port to run the server on"
    )
    args = parser.parse_args()

    run_server(args.port)
