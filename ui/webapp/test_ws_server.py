#!/usr/bin/env python3
import asyncio
from faker import Faker
import json
import random
import time
import uuid
import websockets


fake = Faker()


async def hello(websocket, path):
    while True:
        msg = json.dumps({
            "status": {
                "active": [
                    {"name": fake.name(), "company": fake.company(), "ssn": fake.ssn(), "uuid": uuid.uuid4().hex}
                    for _ in range(20)
                ],
                "historical": [
                    {"name": fake.name(), "company": fake.company(), "ssn": fake.ssn(), "uuid": uuid.uuid4().hex}
                    for _ in range(20)
                ],
            },
            "tasks": [
                {"name": fake.name(), "company": fake.company(), "ssn": fake.ssn(), "uuid": uuid.uuid4().hex}
                for _ in range(20)
            ]
        })
        await websocket.send(msg)
        print(f"> {msg}")
        time.sleep(10)


start_server = websockets.serve(hello, "localhost", 8080)

asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()