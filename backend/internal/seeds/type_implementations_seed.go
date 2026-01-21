package seeds

import (
	"log"

	"github.com/prabalesh/loco/backend/internal/domain"
	"gorm.io/gorm"
)

// SeedTypeImplementations seeds type implementations for all languages
func SeedTypeImplementations(db *gorm.DB) error {
	// Get custom types
	var treeNode, listNode domain.CustomType
	if err := db.Where("name = ?", "TreeNode").First(&treeNode).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "ListNode").First(&listNode).Error; err != nil {
		return err
	}

	// Get languages
	var python, javascript domain.Language
	if err := db.Where("slug = ?", "python").First(&python).Error; err != nil {
		return err
	}
	if err := db.Where("slug = ?", "javascript").First(&javascript).Error; err != nil {
		return err
	}

	implementations := []domain.TypeImplementation{
		// Python - TreeNode
		{
			CustomTypeID: treeNode.ID,
			LanguageID:   python.ID,
			ClassDefinition: `
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right
`,
			DeserializerCode: `
def deserialize_treenode(data):
    if not data or data[0] is None:
        return None
    
    root = TreeNode(data[0])
    queue = [root]
    i = 1
    
    while queue and i < len(data):
        node = queue.pop(0)
        
        # Left child
        if i < len(data) and data[i] is not None:
            node.left = TreeNode(data[i])
            queue.append(node.left)
        i += 1
        
        # Right child
        if i < len(data) and data[i] is not None:
            node.right = TreeNode(data[i])
            queue.append(node.right)
        i += 1
    
    return root
`,
			SerializerCode: `
def serialize_treenode(root):
    if not root:
        return []
    
    result = []
    queue = [root]
    
    while queue:
        node = queue.pop(0)
        if node:
            result.append(node.val)
            queue.append(node.left)
            queue.append(node.right)
        else:
            result.append(None)
    
    # Remove trailing None values
    while result and result[-1] is None:
        result.pop()
    
    return result
`,
		},

		// Python - ListNode
		{
			CustomTypeID: listNode.ID,
			LanguageID:   python.ID,
			ClassDefinition: `
class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next
`,
			DeserializerCode: `
def deserialize_listnode(data):
    if not data:
        return None
    
    dummy = ListNode(0)
    current = dummy
    
    for val in data:
        current.next = ListNode(val)
        current = current.next
    
    return dummy.next
`,
			SerializerCode: `
def serialize_listnode(head):
    result = []
    current = head
    
    while current:
        result.append(current.val)
        current = current.next
    
    return result
`,
		},

		// JavaScript - TreeNode
		{
			CustomTypeID: treeNode.ID,
			LanguageID:   javascript.ID,
			ClassDefinition: `
class TreeNode {
    constructor(val, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}
`,
			DeserializerCode: `
function deserializeTreeNode(data) {
    if (!data || data.length === 0 || data[0] === null) {
        return null;
    }
    
    const root = new TreeNode(data[0]);
    const queue = [root];
    let i = 1;
    
    while (queue.length > 0 && i < data.length) {
        const node = queue.shift();
        
        // Left child
        if (i < data.length && data[i] !== null) {
            node.left = new TreeNode(data[i]);
            queue.push(node.left);
        }
        i++;
        
        // Right child
        if (i < data.length && data[i] !== null) {
            node.right = new TreeNode(data[i]);
            queue.push(node.right);
        }
        i++;
    }
    
    return root;
}
`,
			SerializerCode: `
function serializeTreeNode(root) {
    if (!root) return [];
    
    const result = [];
    const queue = [root];
    
    while (queue.length > 0) {
        const node = queue.shift();
        if (node) {
            result.push(node.val);
            queue.push(node.left);
            queue.push(node.right);
        } else {
            result.push(null);
        }
    }
    
    // Remove trailing nulls
    while (result.length > 0 && result[result.length - 1] === null) {
        result.pop();
    }
    
    return result;
}
`,
		},

		// JavaScript - ListNode
		{
			CustomTypeID: listNode.ID,
			LanguageID:   javascript.ID,
			ClassDefinition: `
class ListNode {
    constructor(val, next = null) {
        this.val = val;
        this.next = next;
    }
}
`,
			DeserializerCode: `
function deserializeListNode(data) {
    if (!data || data.length === 0) {
        return null;
    }
    
    const dummy = new ListNode(0);
    let current = dummy;
    
    for (const val of data) {
        current.next = new ListNode(val);
        current = current.next;
    }
    
    return dummy.next;
}
`,
			SerializerCode: `
function serializeListNode(head) {
    const result = [];
    let current = head;
    
    while (current) {
        result.push(current.val);
        current = current.next;
    }
    
    return result;
}
`,
		},
	}

	for _, impl := range implementations {
		// Check if exists
		var existing domain.TypeImplementation
		if err := db.Where("custom_type_id = ? AND language_id = ?", impl.CustomTypeID, impl.LanguageID).
			First(&existing).Error; err == nil {
			log.Printf("TypeImplementation for Type %d / Lang %d already exists", impl.CustomTypeID, impl.LanguageID)
			continue // Already exists
		}

		// Create
		if err := db.Create(&impl).Error; err != nil {
			return err
		}
		log.Printf("Created TypeImplementation for Type %d / Lang %d", impl.CustomTypeID, impl.LanguageID)
	}

	return nil
}
