import logging

import aiomisc
from client import RmqClient


def run():
    with aiomisc.entrypoint() as loop:
        rmq_client = RmqClient()
        logging.info('rabbitmq client initialized')
        loop.run_until_complete(
            rmq_client.publish(
                body={'python': 'test'},
                routing_key='test',
            ),
        )


def main():
    run()


if __name__ == '__main__':
    main()
