-- Remove seeded data
DELETE FROM tags WHERE slug IN (
    'array', 'hash-table', 'string', 'dynamic-programming', 'math',
    'sorting', 'greedy', 'depth-first-search', 'breadth-first-search',
    'binary-search', 'two-pointers', 'stack', 'queue', 'linked-list',
    'tree', 'graph', 'backtracking', 'bit-manipulation', 'heap', 'sliding-window'
);

DELETE FROM languages WHERE language_id IN ('python', 'javascript', 'cpp', 'java', 'go');
