import argparse
import os
import sys
import importlib.util
import time
import unittest

import yaml


def _pull_image(image):
    os.system(f"docker pull {image}")
    time.sleep(2)


def _load_module_from_path(module_path):
    """
    Dynamically load a Python module from the given file path.
    """
    # Use the filename (without extension) as the module name.
    module_name = module_path.split("/")[-1].rstrip(".py")
    spec = importlib.util.spec_from_file_location(module_name, module_path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def profile(baffs_path, workload_dir, image):
    os.system(f"{baffs_path} shadow --images={image}")
    time.sleep(3)
    test_module = _load_module_from_path(
        os.path.join(workload_dir, "workloads.py"))
    test_module.IMAGE = image
    suite = unittest.defaultTestLoader.loadTestsFromModule(test_module)
    unittest.TextTestRunner(verbosity=2).run(suite)


def debloat(baffs_path, image):
    os.system(f"{baffs_path} debloat --images={image}")
    time.sleep(3)


def validate(workload_dir, image):
    test_module = _load_module_from_path(
        os.path.join(workload_dir, "workloads.py"))
    test_module.IMAGE = image+"-baffs"
    suite = unittest.defaultTestLoader.loadTestsFromModule(test_module)
    res = unittest.TextTestRunner(verbosity=2).run(suite)
    if res.wasSuccessful():
        print("Validation successful")
    else:
        print("Validation failed")
        sys.exit(1)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog='debloater',
        description='Debloat containers')
    parser.add_argument('action', choices=[
                        'profile', 'debloat', 'validate'], help='Action to perform')
    parser.add_argument(
        'dir', type=str, help='Directory containing the workloads script')
    parser.add_argument(
        'baffs', type=str, help='Path to the baffs executable')

    args = parser.parse_args()
    action = args.action
    workloads_dir = args.dir

    with open(f"{args.dir}/config.yml", "r") as f:
        config = yaml.safe_load(f)
    image = config["image"]

    if action == 'profile':
        _pull_image(image)
        profile(args.baffs, workloads_dir, image)
    elif action == 'debloat':
        debloat(args.baffs, image)
    elif action == 'validate':
        validate(workloads_dir, image)
