import http.server
import socketserver
import json
import threading
from urllib.parse import urlparse, parse_qs


class EchoRequestHandler(http.server.BaseHTTPRequestHandler):
    def _set_response_headers(self, content_type="application/json"):
        self.send_response(200)
        self.send_header("Content-Type", content_type)
        # Cast server_address to tuple to fix type error
        port = (
            self.server.server_address[1]
            if isinstance(self.server.server_address, tuple)
            else 0
        )
        self.send_header("Server-Port", str(port))
        self.end_headers()

    def _handle_request(self):
        # Get request information
        # Cast server_address to tuple to fix type error
        port = (
            self.server.server_address[1]
            if isinstance(self.server.server_address, tuple)
            else 0
        )
        request_info = {
            "method": self.command,
            "path": self.path,
            "headers": dict(self.headers.items()),
            "query_params": parse_qs(urlparse(self.path).query),
            "server_port": port,
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


def run_server(port):
    server_address = ("", port)
    httpd = socketserver.TCPServer(server_address, EchoRequestHandler)
    print(f"Server running on port {port}")
    httpd.serve_forever()


def run_multiple_servers(ports):
    threads = []
    for port in ports:
        thread = threading.Thread(target=run_server, args=(port,))
        thread.daemon = True
        threads.append(thread)
        thread.start()

    # Keep the main thread running
    try:
        # Wait for keyboard interrupt
        while True:
            for thread in threads:
                if not thread.is_alive():
                    print(f"A server thread has died. Exiting.")
                    return
            threading.Event().wait(1)
    except KeyboardInterrupt:
        print("Shutting down servers...")


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(description="Run multiple HTTP echo servers")
    parser.add_argument(
        "-p",
        "--ports",
        type=int,
        nargs="+",
        default=[8081, 8082, 8083],
        help="Ports to run the servers on (default: 8081, 8082, 8083)",
    )
    args = parser.parse_args()

    run_multiple_servers(args.ports)
