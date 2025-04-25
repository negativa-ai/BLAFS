import os
import pathlib
import time
import unittest

IMAGE = "set_by_debloat_py"  # set by debloat.py

class TestPyTorchInference(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cwd = pathlib.Path(__file__).parent.resolve()
        docker_cmd = f"docker run --rm -d  --name tmp  -v {cwd}:/app  {IMAGE} tail -f /dev/null"
        os.system(docker_cmd)
        time.sleep(5)

    @classmethod
    def tearDownClass(cls):
        os.system("docker stop tmp")

    def test_inference(self):
        cwd = pathlib.Path(__file__).parent.resolve()
        cmd = f"docker exec -i tmp python /app/inference.py"
        ret = os.system(cmd)
        self.assertEqual(ret, 0, "Inference failed")



        
