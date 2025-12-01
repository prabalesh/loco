-- Insert default supported languages
INSERT INTO languages (language_id, name, version, extension, default_template, executor_config) VALUES

-- Python 3
('python', 'Python 3', '3.11', '.py', 
'from typing import List

class Solution:
    # Problem-specific function goes here
    pass

if __name__ == "__main__":
    # Problem-specific input parsing + call goes here
    pass',
'{"docker_image": "python:3.11-slim", "timeout": 5000, "memory_limit": 256}'::jsonb),

-- JavaScript (Node.js)
('javascript', 'JavaScript (Node.js)', '20.x', '.js',
'class Solution {
    // Problem-specific function goes here
}

function main() {
    // Problem-specific input parsing + call goes here
}

main();',
'{"docker_image": "node:20-slim", "timeout": 5000, "memory_limit": 256}'::jsonb),

-- C++
('cpp', 'C++', '17', '.cpp',
'#include <bits/stdc++.h>
using namespace std;

class Solution {
public:
    // Problem-specific function goes here
};

int main() {
    // Problem-specific input parsing + call goes here
    return 0;
}',
'{"docker_image": "gcc:11", "timeout": 3000, "memory_limit": 256}'::jsonb),

-- Java
('java', 'Java', '17', '.java',
'import java.util.*;

class Solution {
    // Problem-specific function goes here
}

public class Main {
    public static void main(String[] args) {
        // Problem-specific input parsing + call goes here
    }
}',
'{"docker_image": "openjdk:17-slim", "timeout": 5000, "memory_limit": 512}'::jsonb),

-- Go
('go', 'Go', '1.21', '.go',
'package main

import "fmt"

// Problem-specific function goes here

func main() {
    // Problem-specific input parsing + call goes here
}',
'{"docker_image": "golang:1.21-alpine", "timeout": 5000, "memory_limit": 256}'::jsonb)

ON CONFLICT (language_id) DO NOTHING;

-- Insert common tags
INSERT INTO tags (name, slug) VALUES
('Array', 'array'),
('Hash Table', 'hash-table'),
('String', 'string'),
('Dynamic Programming', 'dynamic-programming'),
('Math', 'math'),
('Sorting', 'sorting'),
('Greedy', 'greedy'),
('Depth-First Search', 'depth-first-search'),
('Breadth-First Search', 'breadth-first-search'),
('Binary Search', 'binary-search'),
('Two Pointers', 'two-pointers'),
('Stack', 'stack'),
('Queue', 'queue'),
('Linked List', 'linked-list'),
('Tree', 'tree'),
('Graph', 'graph'),
('Backtracking', 'backtracking'),
('Bit Manipulation', 'bit-manipulation'),
('Heap', 'heap'),
('Sliding Window', 'sliding-window')

ON CONFLICT (name) DO NOTHING;