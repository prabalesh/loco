#!/usr/bin/env python3
"""
Problem Generator Bot for Loco Platform

Usage:
    python problem_generator.py --count 10 --difficulty easy --topic arrays
    python problem_generator.py --file problems.json
"""

import argparse
import json
import requests
import os
from typing import List, Dict

LOCO_API_URL = os.getenv("LOCO_API_URL", "http://localhost:8080")
ADMIN_TOKEN = os.getenv("ADMIN_TOKEN", "")

def generate_problem_with_llm(topic: str, difficulty: str) -> Dict:
    """
    Generate a problem using LLM (OpenAI, Claude, etc.)
    
    This is a placeholder - integrate with your LLM service
    """
    # Example problem structure
    problem = {
        "title": f"Sample {topic.title()} Problem",
        "description": f"A {difficulty} problem about {topic}. Given an array of integers, find the sum of all elements.",
        "difficulty": difficulty,
        "category_ids": [1],  # Arrays category
        "tag_ids": [1],       # Sample tag
        "function_name": "solve",
        "return_type": "int",
        "parameters": [
            {"name": "nums", "type": "int[]", "is_custom": False}
        ],
        "validation_type": "EXACT",
        "expected_time_complexity": "O(n)",
        "expected_space_complexity": "O(1)",
        "test_cases": [
            {
                "input": [[1, 2, 3, 4, 5]],
                "expected_output": 15,
                "is_sample": True
            },
            {
                "input": [[10, 20, 30]],
                "expected_output": 60,
                "is_sample": False
            }
        ],
        "reference_solution": {
            "language_slug": "python",
            "code": "def solve(nums):\n    return sum(nums)"
        }
    }
    
    return problem

def bulk_import_problems(problems: List[Dict], validate_references: bool = True):
    """
    Import problems to Loco via bulk API
    """
    url = f"{LOCO_API_URL}/api/v2/admin/problems/bulk"
    
    payload = {
        "problems": problems,
        "options": {
            "validate_references": validate_references,
            "skip_duplicates": True,
            "stop_on_error": False
        }
    }
    
    headers = {
        "Authorization": f"Bearer {ADMIN_TOKEN}",
        "Content-Type": "application/json"
    }
    
    print(f"Importing {len(problems)} problems...")
    
    try:
        response = requests.post(url, json=payload, headers=headers)
        
        if response.status_code in [200, 206]:
            result = response.json()
            print(f"✅ Import completed!")
            print(f"   Created: {result['total_created']}")
            print(f"   Failed: {result['total_failed']}")
            print(f"   Time: {result['processing_time_ms']}ms")
            
            if result['failed_problems']:
                print(f"\n❌ Failed problems:")
                for failure in result['failed_problems']:
                    print(f"   [{failure['index']}] {failure['title']}: {failure['error_message']}")
                    if 'errors' in failure and failure['errors']:
                        for err in failure['errors']:
                            print(f"      - {err}")
            
            if result['created_problems']:
                print(f"\n✅ Created problems:")
                for success in result['created_problems']:
                    print(f"   [{success['index']}] {success['title']} → {success['slug']} (status: {success['validation_status']})")
            
            return result
        else:
            print(f"❌ Import failed: {response.status_code}")
            print(response.text)
            return None
    except Exception as e:
        print(f"❌ Request failed: {e}")
        return None

def main():
    parser = argparse.ArgumentParser(description="Generate and import problems to Loco")
    parser.add_argument("--count", type=int, default=3, help="Number of problems to generate")
    parser.add_argument("--difficulty", choices=["easy", "medium", "hard"], default="easy")
    parser.add_argument("--topic", default="arrays", help="Problem topic")
    parser.add_argument("--file", help="Import from JSON file")
    parser.add_argument("--no-validate", action="store_true", help="Skip reference solution validation")
    
    args = parser.parse_args()
    
    if not ADMIN_TOKEN:
        print("❌ ADMIN_TOKEN environment variable not set")
        # In a real scenario, you'd exit here, but for demonstration we'll continue
        # return
    
    if args.file:
        # Import from file
        try:
            with open(args.file, 'r') as f:
                data = json.load(f)
                if isinstance(data, dict) and "problems" in data:
                    problems = data["problems"]
                else:
                    problems = data
        except Exception as e:
            print(f"❌ Failed to read file: {e}")
            return
    else:
        # Generate problems
        print(f"Generating {args.count} {args.difficulty} problems about {args.topic}...")
        problems = []
        for i in range(args.count):
            problem = generate_problem_with_llm(args.topic, args.difficulty)
            problem["title"] = f"{problem['title']} {i+1}"  # Make unique
            problems.append(problem)
    
    # Import
    bulk_import_problems(problems, validate_references=not args.no_validate)

if __name__ == "__main__":
    main()
