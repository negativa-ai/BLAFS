import os
import pathlib
import redis
import time
import unittest

import yaml

IMAGE = "set_by_debloat_py"  # set by debloat.py


class TestRedis(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        docker_cmd = f"docker run --rm -d  --name tmp --network host {IMAGE}"
        os.system(docker_cmd)
        time.sleep(5)
        cls.redis_client = redis.StrictRedis(
            host="localhost", port=6379, decode_responses=True
        )

    @classmethod
    def tearDownClass(cls):
        os.system("docker stop tmp")

    def test_basic_key_value_operations(self):
        # Set and get a value
        self.redis_client.set("key1", "value1")

        self.assertEqual(self.redis_client.get("key1"), "value1")

        # Expire key after 5 seconds
        self.redis_client.set("key2", "value2", ex=5)
        # Wait 6 seconds and check the key
        time.sleep(6)
        self.assertIsNone(self.redis_client.get("key2"))

    def test_data_structure(self):
        self.redis_client.rpush("mylist", "item1", "item2", "item3")
        self.assertListEqual(
            self.redis_client.lrange(
                "mylist", 0, -1), ["item1", "item2", "item3"]
        )
        # Hash operations
        self.redis_client.hset(
            "user:1000", mapping={"name": "Alice", "age": 30, "location": "NY"}
        )
        result = self.redis_client.hgetall("user:1000")
        expected = {
            "name": "Alice",
            "age": "30",
            "location": "NY",
        }  # Redis stores values as strings
        self.assertEqual(result, expected,
                         "Hash contents do not match expected values")

        # Set operations
        self.redis_client.sadd("myset", "apple", "banana", "cherry")
        result = self.redis_client.smembers("myset")
        expected = {"apple", "banana", "cherry"}
        self.assertEqual(result, expected,
                         "Set contents do not match expected values")

        # Sorted Set operations
        self.redis_client.zadd(
            "scores", {"Alice": 100, "Bob": 200, "Charlie": 150})
        result = self.redis_client.zrange("scores", 0, -1, withscores=True)
        expected = [
            ("Alice", 100.0),
            ("Charlie", 150.0),
            ("Bob", 200.0),
        ]  # Sorted by score
        self.assertEqual(
            result, expected, "Sorted set contents do not match expected values"
        )

    def test_subscriptions(self):
        def subscribe_to_channel(channel_name):
            pubsub = self.redis_client.pubsub()
            pubsub.subscribe(channel_name)

            print(f"Subscribed to channel: {channel_name}")
            for message in pubsub.listen():
                if message["type"] == "message":
                    self.assertEqual(
                        message["data"], "Breaking News: Redis is awesome!"
                    )
                    break

        # Publisher function
        def publish_to_channel(channel_name, message):
            self.redis_client.publish(channel_name, message)

        # Run Publisher and Subscriber
        import threading

        # Start a subscriber in a separate thread
        threading.Thread(
            target=subscribe_to_channel, args=("news",), daemon=True
        ).start()

        # Publish messages
        time.sleep(1)  # Wait for subscriber to be ready
        publish_to_channel("news", "Breaking News: Redis is awesome!")

        time.sleep(1)  # Allow time for messages to be processed

    def test_caching(self):
        def slow_function(x):
            time.sleep(3)
            return x * x

        def cached_function(x):
            cache_key = f"slow_function:{x}"
            if self.redis_client.exists(cache_key):
                print(f"Cache hit for {x}")
                return int(self.redis_client.get(cache_key)), "hit"
            else:
                print(f"Cache miss for {x}")
                result = slow_function(x)
                # Cache for 10 seconds
                self.redis_client.set(cache_key, result, ex=10)
                return result, "miss"

        miss = cached_function(4)
        self.assertEqual(miss[0], 16)
        self.assertEqual(miss[1], "miss")

        hit = cached_function(4)
        self.assertEqual(hit[0], 16)
        self.assertEqual(hit[1], "hit")
