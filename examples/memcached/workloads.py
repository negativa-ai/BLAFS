import os
import time
import unittest

from pymemcache.client.base import Client

IMAGE = "set_by_debloat_py"  # set by debloat.py


class TestMemcached(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        docker_cmd = f"docker run --rm -d  --name tmp --network host {IMAGE}"
        os.system(docker_cmd)
        time.sleep(5)
        cls.client = Client(("localhost", 11211))

    @classmethod
    def tearDownClass(cls):
        os.system("docker stop tmp")

    def test_set_and_get(self):
        """Test setting and getting a key-value pair."""
        key, value = "test_key", b"test_value"
        self.client.set(key, value)
        retrieved_value = self.client.get(key)

        self.assertEqual(retrieved_value, value, "Retrieved value does not match stored value")

    def test_delete_key(self):
        """Test deleting a key."""
        key, value = "delete_key", "delete_value"
        self.client.set(key, value)
        self.client.delete(key)

        retrieved_value = self.client.get(key)
        self.assertIsNone(retrieved_value, "Key was not deleted")

    def test_key_expiration(self):
        """Test setting a key with an expiration time."""
        key, value, expiration = "exp_key", "exp_value", 2
        self.client.set(key, value, expire=expiration)

        time.sleep(expiration + 1)  # Wait for the key to expire
        retrieved_value = self.client.get(key)
        self.assertIsNone(retrieved_value, "Key did not expire as expected")

    def test_multiple_keys(self):
        """Test setting and retrieving multiple keys."""
        items = {"key1": b"value1", "key2": b"value2", "key3": b"value3"}
        self.client.set_many(items)

        retrieved_items = self.client.get_many(items.keys())
        self.assertEqual(retrieved_items, items, "Retrieved values do not match stored values")

    def test_increment_and_decrement(self):
        """Test incrementing and decrementing a key."""
        key, initial_value = "counter", 10
        self.client.set(key, initial_value)

        # Increment
        new_value = self.client.incr(key, 5)
        self.assertEqual(new_value, initial_value + 5, "Increment operation failed")

        # Decrement
        new_value = self.client.decr(key, 3)
        self.assertEqual(new_value, initial_value + 2, "Decrement operation failed")

    def test_flush_all(self):
        """Test flushing all keys."""
        items = {"flush_key1": "value1", "flush_key2": "value2"}
        self.client.set_many(items)

        self.client.flush_all()

        retrieved_items = self.client.get_many(items.keys())
        self.assertFalse(retrieved_items, "Keys were not flushed")
