import asyncio
import json
import logging

from aio_pika import DeliveryMode, ExchangeType, Message, connect


class RmqClient():
    _conn = "amqp://guest:guest@localhost:5672/"

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

    async def publish(self, body, routing_key):
        connection = await self._establish_connection()
        channel = await connection.channel()
        exchange = await channel.declare_exchange(
            name="test",
            type=ExchangeType.DIRECT,
            durable=True,
        )
        queue = await channel.declare_queue(name="test", durable=True)
        await queue.bind(exchange=exchange)
        message = Message(
            body=json.dumps(obj=body).encode(),
            content_type="application/json",
            delivery_mode=DeliveryMode.PERSISTENT,
        )
        await exchange.publish(
            message=message,
            routing_key=routing_key,
        )
        logging.debug(f"sent: {message!r}")
        await connection.close()
        return {"status": "accepted"}
