import logging

import aiomisc

from . client import RmqClient


def run():
    with aiomisc.entrypoint() as loop:
        rmq_client = RmqClient()
        logging.info("rabbitmq client initialized")
        loop.create_task(rmq_client.publish(
            body={"foo": "bar"},
            routing_key="test",
            ),
        )
        loop.run_forever()


def main():
    run()


if __name__ == "__main__":
    main()
