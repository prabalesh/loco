# Bulk Import API Documentation

The Bulk Import API allows administrators to create multiple problems at once using a single JSON request. This is useful for automated problem generation from LLMs or importing datasets from external sources.

## Endpoints

### Synchronous Import
`POST /api/v2/admin/problems/bulk`

- **Description**: Processes imports immediately and returns results for each problem.
- **Batch Limit**: 100 problems per request.
- **Authentication**: Admin JWT token required.

### Asynchronous Import
`POST /api/v2/admin/problems/bulk-async`

- **Description**: Queues the import task and returns a job ID. Useful for larger batches.
- **Batch Limit**: 1000 problems per request.
- **Authentication**: Admin JWT token required.

## Request Format

```json
{
  "problems": [
    {
      "title": "Problem Title",
      "description": "At least 20 characters",
      "difficulty": "easy | medium | hard",
      "category_ids": [1, 2],
      "tag_ids": [3, 4],
      "function_name": "solve",
      "return_type": "int",
      "parameters": [
        { "name": "nums", "type": "int[]", "is_custom": false }
      ],
      "validation_type": "EXACT | UNORDERED | SUBSET | ANY_MATCH",
      "expected_time_complexity": "O(n)",
      "expected_space_complexity": "O(1)",
      "test_cases": [
        {
          "input": [[1, 2, 3]],
          "expected_output": 6,
          "is_sample": true
        }
      ],
      "reference_solution": {
        "language_slug": "python",
        "code": "def solve(nums): return sum(nums)"
      }
    }
  ],
  "options": {
    "validate_references": true,
    "skip_duplicates": true,
    "stop_on_error": false
  }
}
```

## Response Format (Sync)

```json
{
  "total_submitted": 10,
  "total_created": 8,
  "total_failed": 2,
  "created_problems": [
    {
      "index": 0,
      "title": "Title",
      "slug": "title",
      "problem_id": 123,
      "validation_status": "validated"
    }
  ],
  "failed_problems": [
    {
      "index": 5,
      "title": "Failed Title",
      "errors": ["title must be 5-200 characters"],
      "error_message": "title must be 5-200 characters"
    }
  ],
  "processing_time_ms": 1250
}
```

## Rate Limiting

- **Limit**: 10 requests per hour per user.
- **Header**: `X-RateLimit-Limit`, `X-RateLimit-Remaining`.

## Bot Script

A sample Python bot script is available at `bots/problem_generator.py`. 
To run:
```bash
export ADMIN_TOKEN="your_token"
python bots/problem_generator.py --count 10 --topic strings
```
