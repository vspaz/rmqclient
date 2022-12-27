import json
import logging
import os

import aiomisc
from aio_pika import DeliveryMode, IncomingMessage

from .client import RmqClient

async def listen_for_messages(rmq_client: RmqClient):
    async def on_message_received(message: IncomingMessage):
        async with message.process(requeue=True):
            try:
                resp = json.loads(message.body)
                logging.debug(f"message received: {resp!r}")
            except Exception as err:
                logging.error(err)

    await rmq_client.subscribe(
        queue_name="test_queue",
        on_message=on_message_received)


def run():
    config = dict(
        user=os.getenv("RABBITMQ_USER", "guest"),
        password=os.getenv("RABBITMQ_PASSWORD", "guest"),
        host=os.getenv("RABBITMQ_HOST", "localhost"),
        port=os.getenv("RABBITMQ_PORT", "5672"),
    )
    with aiomisc.entrypoint() as loop:
        rmq_client = RmqClient(config=config)
        logging.info("rabbitmq client initialized")
        loop.create_task(listen_for_messages(rmq_client=rmq_client))
        loop.run_forever()


def main():
    run()


if __name__ == "__main__":
    main()
