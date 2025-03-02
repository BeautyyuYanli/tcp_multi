#!/usr/bin/env python3

import asyncio
import logging
import websockets

# Configure logging
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger("websocket_server")


async def echo(websocket, port) -> None:
    """
    Echo handler that receives messages and sends them back with server port info.

    Args:
        websocket: The WebSocket connection to the client
        port: The port number this server is running on
    """
    client_id = id(websocket)
    logger.info(f"Client {client_id} connected to server on port {port}")

    try:
        async for message in websocket:
            logger.info(
                f"Received message from client {client_id} on port {port}: {message}"
            )
            response = f"Server on port {port} echoes: {message}"
            await websocket.send(response)
            logger.info(f"Echoed message back to client {client_id} on port {port}")
    except websockets.exceptions.ConnectionClosed as e:
        logger.info(f"Client {client_id} disconnected from server on port {port}: {e}")
    except Exception as e:
        logger.error(f"Error handling client {client_id} on port {port}: {e}")
    finally:
        logger.info(f"Client {client_id} connection closed on port {port}")


async def start_server(host: str, port: int) -> None:
    """
    Start a WebSocket server on the specified host and port.

    Args:
        host: The hostname to bind to
        port: The port to bind to
    """
    logger.info(f"Starting WebSocket server on {host}:{port}")
    # Create a handler that includes the port information
    handler = lambda websocket: echo(websocket, port)
    async with websockets.serve(handler, host, port):
        await asyncio.Future()  # Run forever


async def main() -> None:
    """Start multiple WebSocket servers"""
    host = "localhost"
    ports = [8081, 8082, 8083]

    # Create tasks for each server
    tasks = []
    for port in ports:
        tasks.append(asyncio.create_task(start_server(host, port)))

    # Wait for all tasks to complete (they won't unless there's an error)
    await asyncio.gather(*tasks)


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Servers stopped by user")
    except Exception as e:
        logger.error(f"Server error: {e}")
