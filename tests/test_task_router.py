import unittest
import os
import sys

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'scripts')))
import task_router

class TestTaskRouter(unittest.TestCase):

    def test_route_to_gpt_for_code(self):
        task = "Write a Python script to calculate fibonacci and add unit tests."
        model, reason = task_router.route_task(task)
        self.assertEqual(model, "GPT")

    def test_route_to_claude_for_architecture(self):
        task = "Refactor the current directory structure and write new documentation."
        model, reason = task_router.route_task(task)
        self.assertEqual(model, "Claude")

    def test_route_to_gemini_for_scanning(self):
        task = "Scan all repos and perform a fast performance script update."
        model, reason = task_router.route_task(task)
        self.assertEqual(model, "Gemini")

    def test_ambiguous_task_defaults_to_claude(self):
        task = "Look at this issue."
        model, reason = task_router.route_task(task)
        self.assertEqual(model, "Claude")
        self.assertTrue("Default" in reason)

if __name__ == '__main__':
    unittest.main()
