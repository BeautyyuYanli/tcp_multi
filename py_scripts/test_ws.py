#!/usr/bin/env python3

import asyncio
import logging
import websockets
from websockets.exceptions import ConnectionClosed

# Configure logging
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger("websocket_client")


async def connect_and_send(client_num: int) -> None:
    """
    Connect to the WebSocket server and send messages.

    Args:
        client_num: The client number (used for identification)
    """
    uri = "ws://localhost:8443"

    try:
        async with websockets.connect(uri) as websocket:
            logger.info(f"Client {client_num} connected")

            # Send messages from 1 to 100 for this client
            for i in range(3):
                message = f"Client {client_num} - Message {i}"
                await websocket.send(message)
                logger.info(f"Client {client_num} sent: {message}")

                # Receive the echo response
                response = await websocket.recv()
                logger.info(f"Client {client_num} received: {response}")

                # Add a small delay between messages
                await asyncio.sleep(0.1)

    except ConnectionClosed as e:
        logger.error(f"Client {client_num} connection closed: {e}")
    except Exception as e:
        logger.error(f"Client {client_num} error: {e}")


async def main() -> None:
    """Create multiple WebSocket connections and send data"""
    # Create tasks for each client (1 to 100)
    tasks = [connect_and_send(i) for i in range(3)]

    # Run all client tasks concurrently
    await asyncio.gather(*tasks)


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Client stopped by user")
    except Exception as e:
        logger.error(f"Client error: {e}")
