import asyncio
import json
import logging
import os

import aiomisc
from aio_pika import DeliveryMode, IncomingMessage, Message, connect


class PubSub():

    def __init__(self, config):
        self._conn = "amqp://{user}:{password}@{host}:{port}/".format(**config)

    async def _establish_connection(self):
        return await connect(
            url=self._conn,
            loop=asyncio.get_running_loop(),
            timeout=10,
        )

    async def subscribe(self, queue_name, on_message):
        connection = await self._establish_connection()
        channel = await connection.channel()
        await channel.set_qos(prefetch_count=8)
        queue = await channel.declare_queue(queue_name, durable=True)
        logging.info(f"{queue_name} starting to consume messages")
        await queue.consume(on_message)

    async def publish(self, body, route_key):
        connection = await self._establish_connection()
        channel = await connection.channel()
        message = Message(
            body=json.dumps(obj=body).encode(),
            content_type="application/json",
            delivery_mode=DeliveryMode.PERSISTENT,
        )
        await channel.default_exchange.publish(message, routing_key=route_key)
        logging.debug(f"sent: {message!r}")
        await connection.close()
        return {"status": "accepted"}


async def listen_for_messages(rmq_client: PubSub):
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
        rmq_client = PubSub(config=config)
        logging.info("rabbitmq client initialized")
        loop.create_task(listen_for_messages(rmq_client=rmq_client))
        loop.run_forever()


def main():
    run()


if __name__ == "__main__":
    main()
