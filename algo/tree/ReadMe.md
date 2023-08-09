# Tree data structure  

A tree is a nonlinear hierarchical data structure that consists of nodes connected by edges.  

## Why Tree Data Structure?
Other data structure such as arrays, linked list, stack and queue are linear data structures that score data 
sequentially. in order to perform any operation in a linear data structure, the time complexity increases with 
the increase in the data size. But, it is not acceptable in today's computational world.
  
Different tree data structures allow quicker and easier access to the data as it is a non-linear data structure.

## Tree Terminologies

### Node
A node is an entity that contains a key or value and pointers to ts child nodes.

The last nodes of each path are called leaf nodes or external nodes that do not contain a link/pointer to child nodes.

The node having at least a child node is called an internal node.

### Edge
It is the link between any two nodes.

### Root
It is the topmost node of a tree.

### Height of a Node
The height of a node is the number of edges from the node to the deepest leaf (i.e. the longest path from the node to a leaf node).

### Depth of a Node
The depth of a node is the number of edges from the root to the node.

### Height of a Tree
The height of a Tree is the height of the root node or the depth of the deepest node.

### Degree of a Node
The degree of a node is the total number of branches of that node.

### Forest
A collection of disjoint trees is called a forest.

# Tree Traversal
Traversing a tree means visiting every node in the tree. You might, for instance, want to add all the values in the tree or find the largest one. For all these operations, you will need to visit each node of the tree.

Every tree is a combination of:
1.A node carrying data
2.Two subtrees (left subtree, right subtree)

## Inorder
1.First, visit all the nodes in the left subtree
2.Then the root node
3.Visit all the nodes in the right subtree

```
inorder(root.left)
display(root)
inorder(root.right)
```

## Preorder
1.Visit root node
2.Visit all the nodes in the left subtree
3.Visit all the nodes in the right subtree

```
display(root)
preorder(root.left)
preorder(root.right)
```

## Postorder
1.Visit all the nodes in the left subtree
2.Visit all the nodes in the right subtree
3.Visit the root node

```
postorder(root.left)
postorder(root.right)
display(root)
``` 

# Binary tree
- data item
- address of left child
- address of right child

## full binary tree (02 tree)
every parent node/internal node has either two or no children


### full binary tree theorems
let
i = the number of internal nodes
n = be the total number of nodes
l = number of leaves
lambda = number of levels

1.the number of leaves i + 1
2.the total number of nodes is n = 2i + 1
3.the number of internal nodes is  i = (n - 1) / 2
4.the number of leaves is l = (n + 1) / 2
5.the total number of nodes is 2l - 1
6.the number of internal nodes is i = l - 1
7.the number of leaves is at most 2^(lambda - 1)

## perfect binary tree(2 tree)
every internal node has exactly two child nodes and all the leaf nodes are at the same level

all the internal nodes have degree of 2

1.if a single node has no children, it is a perfect binary tree of height h = 0
2.if a node has h > 0, it is a perfect binary tree if both of its subtrees are of height h - 1 and are non-overlapping

### perfect binary tree theorems
1.a perfect binary tree of height h has 2^(h + 1) - 1 node
2.a perfect binary tree with n nodes has height log(n + 1) - 1 = theta(ln(n))
3.a perfect binary tree of height h has 2^h leaf nodes
4.the average depth of a node in a perfect binary tree is theta(ln(n))

## complete binary tree(021 tree)
just like a full binary tree, but with two major differences
1.every level must be completely filled
2.all the leaf elements must lean towards the left
3.the last leaf element might not have a right sbiling (i.e. a complete binary tree doesn't have to be a full binary tree)

## degenerate or pathological tree
the tree having a single child either left or right

## skewed binary tree
a pathological/degenerate tree in which the tree is either dominated by the left nodes or thr right nodes
left-skewed binary tree
right-skewed binary tree

## balanced binary tree(-1 tree)
the difference between the height of the left and the right subtree for each node is either 0 or 1

df = abs(height of left child - height of right child)

### the conditions for height-balance binary tree
1.difference between the left and the right subtree for any node is not more than one
2.the left subtree is balanced
3.the right subtree is balanced

### applications
AVL tree
balanced binary search tree

## AVL tree
it is a self-balancing binary search tree in which each node maintains extra information called a balance factor whose value is either -1, 0, or +1

Inventor: Georgy Adelson-Velsky and Landis

### balance factor
balance factor of a node in avl tree is the difference between the height of the left subtree and that of the right subtree of that node

bf = height of left subtree - height of right subtree
or
bf = height of right subtree - height of left subtree

### rotating the subtrees in an avl tree
- left rotate
the arrangement of the nodes on the right is transformed into the arrangements on the left node.
  
initial
root-> x node
x -> left: alpha node
x -> right: y node
y -> left: beta node
y -> right: gamma node

after left-rotate
root -> y node
y -> left: x node
y -> right: gamma node
x -> left: alpha node
x -> right: beta node

- right rotate
  the arrangement of the nodes on the left is transformed into the arrangements on the right node.
  
- left-right and right-left rotate
in left-right rotation, the arrangements are first shifted to the left and then to the right
  
initial
p -> z node
z -> left: x node
z -> right: delta node
x -> left: alpha node
x -> right: y node
y -> left: beta node
y -> right: gamma node

1.do left rotate on x-y
p -> z node
z -> left: y node
z -> right: delta node
y -> left: x node
y -> right: gamma node
x -> left: alpha node
x -> right: beta node
(left-skewed tree)

2.do right rotation on y-z
p -> y node
y -> left: x node
y -> right: z node
x -> left: alpha node
x -> right: beta node
z -> left: gamma node
z -> right: delta node

in right-left rotation, the arrangements are first shifted to the right and then to the left

insert a new node
always inserted as a leaf node with balance factor equal to 0

1.go to the appropriate leaf node to insert a new node using the following recursive steps. Compare new key with root key of the current tree.
a.if new key  < root key, call insertion algorithm on the left subtree of the current node until the leaf node is reached.
b.else is new key > root key, call insertion algorithm on the right subtree of current node until the leaf node is reached.
c.else, return leaf node
2.compare leaf key obtained from the above steps with new key
a.if new key < leaf key, make new node as the left child of leaf node
b.else, make new node as right child of leaf node
3.update balance factor of the nodes
4.if the nodes are unbalanced, then rebalance the node
a.if balance factor > 1, it means the height of the left subtree is greater than that of the right subtree. So, do a right rotation or left-right rotation
- if new node key < left child key do right rotation
- else, do left-right rotation
b.if balance factor < -1, it means the height of the right subtree is greater than that of the left subtree. So, do right rotation or right-left rotation
  - if new node key > right child key do left rotation
    - else, do right-left rotation
    
delete a node
a node is always deleted as leaf node. After deleting a node, the balance factors of the nodes get changes. In order to rebalance the balance factor, suitable rotations are performed.

1.locate node to be deleted (recursion is used to find node to be deleted)
2.there are three cases for deleting a node
a.if node to be deleted is the leaf node, then remove node to be deleted
b.if node to be deleted has one child, then substitute the contents of node to be deleted with that of the child. remove the child
c.if node to be deleted has two children, find the inorder successor w of node to be deleted (i.e. node with a minimum value of key in the right subtree)
3.update balance factor of the nodes
4.rebalance the tree if the balance factor of any of the nodes is not equal to -1, 0 or 1
a.if balance factor of current node > 1
- if balance factor of left child >= 0, do right rotation
- else do left-right rotation
b.if balance factor of current node < -1
- if balance factor of right child <= 0, do left rotation
- else do right-left rotation

## B-tree
self-balancing search tree in which each node can contain more than one key and can have more than two children. It is a generalized form of the binary search tree.

it is also known as a height-balanced m-way tree.

### why b-tree?
the need for b-tree arose with the rise in the need for lesser time in accessing the physical storage media like a hard disk. The secondary storage devices are slower with a large capacity. There was a need for such types of data structures that minimize the disk accesses.

Other data structures such as a binary search tree, avl tree, red-black tree, etc can storage only one key in one node. if you have to store a large number of keys, then the height of such trees becomes very large and the access time increases.

however, b-tree can store many keys in a single node and can have multuple child nodes. This decreases the height significantly allowing faster disk accesses.

### b-tree properties
1.for each node x, the keys are stored in increasing order
2.in each node, there is a boolean value x.leaf which is true if x is a leaf
3.if n is the order of the tree, each internal node can contain at most n - 1 keys along with a pointer to each child
4.each node except root can have at most n children and at least n / 2 children
5.all leaves have the same depth (i.e. height-h of the tree)
6.the root has at least 2 children and contains a minimum of 1 key
7.if n >= 1, then for any n-key b-tree of height h and minimum degree t >= 2, h >= log{t}{(n + 1) / 2}

### searching
1.starting from the root node, compare k with the first keyof the node. if k = the first key of the node, return the node and the index.
2.if k.leaf = true, return null (not found)
3.if k < the first key of the root node, search the left child of this key recursively.
4.if there is more than one key in the current node and k > the first key, compare k with the next key in the node.
if k < next key, search the left child of this key (i.e. k lies in between the first and the second keys).
else, search the right child of the key
5.repeat steps 1 to 4 until the leaf us reached

### insertion 
Inserting an element on a b-tree consists of two events: searching the appropriate node to insert the element and splitting the node if required. Insertion operation always takes place in the bottom-up approach.

1.if the tree is empty, allocate a root node and insert the key.
2.update the allowed number of keys in the node.
3.search the appropriate node for insertion.
4.if the node is full, follow the steps below.
5.insert the elements in increasing order.
6.now, there are elements greater than its limit. So, split at the median.
7.push the median key upwards and make the left keys as a left child and the right keys as right child.
8.if the node is not full, follow the steps below.
9.insert the node in increasing order.

### deletion 
deleting an element on a b-tree consists of 3 main events:
1.searching the node where the key to be deleted exists
2.deleting the key
3.balancing the tree if required

while deleting a tree, a condition called underflow may occur. underflow occurs when a node contains less than the minimum number of keys is should be hold.

the terms to be understood before studying deletion operation are:
1.inorder predecessor
the largest key on the left child of a node is called its inorder predecessor
2.inorder successor
the smallest key on the right child of a node is called its inorder successor

#### deletion operation
before going through the steps below, one must know these facts about a b-tree of degree m
1.a node can have a maximum of m children
(i.e. 3)
2.a node can contain a maximum of m - 1 keys
(i.e. 2)
3.a node should have a minimum of m/2 children
(i.e. 2)
4.a node (except root node) should contain a minimum of m/2 - 1 keys
(i.e. 1)
向上取整

case 1
the key to be deleted lies in the leaf. there are 2 cases for it.
1.the deletion of the key does not violate the property of the minimum number of keys a node should hold.
2.the deletion of the key violates the property of the minimum number of keys a node should hold, In this case, we borrow a key from its immediate neighboring sibling node in the order of left to right.

first, visit the immediate left sibling. if the left sibling node has more than a minimum number of keys, then borrow a key from this node.
else, check to borrow from the immediate right sibling node.

if both the immediate sibling nodes already have a minimum number of keys, then merge the node with either the left sibling node or the right sibling node. The merging is done through the parent node.

case 2
if the key to be deleted lies in the internal node, the following cases occurs.
1.the internal node, which is deleted, is replaced by an inorder predecessor if the left child has more than the minimum number of keys.
2.the internal node, which is deleted, is replaced by an inorder successor if the right child has more than the minimum number of keys.
3.if either child has exactly a minimum number of keys then, merge the left and right children
after merging if the parent node has less than the minimum number of keys then, look for the sibling as in case 1

case 3
in this case, the height of the tree shrinks. if the target key lies in an internal node, and the deletion of the key leads to a fewer number of keys in the node (i.e. less than the minimum required), then look for the inorder predecessor and the inorder successor. if both the children contain a minimum number of keys then, borrowing cannot take place. this leads to case 2. i.e. merging the children.

again, look for the sibling to borrow a key. but, if the sibling also has only a minimum number of keys then, merge the node with the sibling along with the parent. Arrange the children accordingly (increasing order)

## red black tree
2-3 tree

it is a self-balancing binary search tree in which each node contains an extra bit for denoting the color of the node, either red or black.

a red-black tree satisfies the following properties:
1.red/black properties: every node is colored, either red or black
2.root property: the root is black
3.leaf property: every leaf(nil) is black
4.red property: if a red node has children then, the children are always black
5.depth property: for each node, any simple path from this node to any of its descendant leaf has the same black-depth
(the number of black nodes)

attributes
color, key, leftChild, rightChild, parent(except root node)

### how the red-black tree maintains the property of self-balancing?
the red-black color is meant for balancing the tree.

the limitations put on the node ensure that any simple path from the root to a leaf is not more than twice as long as any other such path. it helps in maintaining the self-balancing property og the red-black tree.

### operations on a red-black tree
rotating the subtrees in a red-black tree
in rotation operation, the positions of the nodes of a subtree are interchanged.
rotation operation is used for maintaining the properties of a red-black tree when they are violated by other operations such as insertion and deletion.

#### left rotate
the arrangement of the nodes on the right is transformed into the arrangements on the left node

#### right rotate
the arrangement of the nodes on the left is transformed into the arrangements on the right node.

#### left-right rotate
the arrangements are first shifted to the left and then to the right

#### right-left rotate
the arrangements are first shifted to the right and then to the left

#### inserting an element into a red-black tree
while inserting a new node, the new node is always inserted as red node. after insertion of a new node, if the tree is violating the properties of the red-black tree then, we do the following operations
1.recolor
2.rotation

#### insert algorithm
1.let y be the leaf (i.e. nil) and x be the root of the tree
2.check if the tree is empty (i.e. whether x is nil). if yes, insert newNode as a root node and color it black.
3.else, repeat steps following steps until leaf (nil) is readched.
a)compare newKey with rootKey
b)if newKey is greater than rootKey, traverse through the right subtree.
c)else traverse through the left subtree
4.assign the parent of the leaf as a parent of newNode
5.if leafKey is greater than newKey, make newNode as rightChild
6.else, make newNode as leftChild
7.assign null to the left and rightChild of newNode.
8.assigned red color to newNode
9.call insertFix-algorithm to maintain the property of red-black tree if violated

why newly inserted nodes are always red in a red-black tree?
this is because inserting a red node does not violate the depth property of a red-black tree.
if you attach a red node to a red node, then the rule is violated but it is easier to fix this problem than the problem than introduced by violating the depth property.

algorithm to maintain red-black property after insertion
this algorithm is used for maintaining the property of a red-black tree if insertion ofa newNode violates this property.
1.do the following while the parent of newNode p is RED
2.if p is the left child of grandParent gP of z, do the following.
case1:
a)if the color of the right child of gP of z is RED, set the color of both the children of gP as BLACK and the color of gP as RED.
b)assign gP to newNode.
case2:
a)else if newNode is the right child of p then, assign p to newNode.
b)left-rotate newNode
case3:
a)set colorof p as BLACK and color of gP as RED
b)right-rotate gP
3.else, do the following
a)if the color of the left child of gP of z is RED, set the color of both the children of gP as BLACK and the colorof gP as RED
b)assign gP to newNode
c)else if newNode is the left of p then, assign p to newNode and right-rotate newNode
d)set colorofp as BLACK and color of gP as RED
e)left-rotate gP
4.set the root of the tree as black

algorithm to delete a node
1.save the color of nodeToBeDeleted in originalColor
2.if the left child of nodeToBeDeleted is null
a)assigne the right child of nodeToBeDeleted to x
b)transplant nodeToBeDeleted with x
3.else if the right child of nodeToBeDeleted is null
a)assign the left child of nodeToBeDeleted is null
b)transplant nodeToBeDeleted with x
4.else
a)assign the minimum of right subtree of nodeToBeDeleted into y.
b)save the colorof y in originalColor
c)assign the rightChild of y into x
d)if y is a child of nodeToBeDeleted, then set the parent of x as y.
e)else, transplant y with rightChild of y
f)transplate nodeToBeDeleted with y
g)set the color of y with originalColor
5.if the originalColor is BLACK, call DeleteFix(x)

algorithm to maintain red-black property after deletion
this algorithm is implemented when a black node is deleted because it violates the black depth property of the red-black tree

this violation is corrected by assuming that node x (which is occupying y's original position) has an extra black. this makes node x neither red nor black. it is either doubly black or black-and-red. this violates the red black properties.

however, the color attribute of x is not changed rather the extra black is represented in x's pointing to the node.

the extra black can be removed if
1.it reaches the root node
2.if x points to a red-black node. in this case, x is colored black
3.suitable rotations and recoloring are performed

the following algorithm retains the properties of a red-black tree
1.do the following until the x is not the root of the tree and the color of x is BLACK
2.if x is the left child of its parent then,
a)assign w to the sibling of x
b)if the right child of parent of x is RED
case1:
a>set the color of the right child of the parent of x as BLACK
b>set the color of the parent of x as RED
c>left-rotate the parent of x
d>assign the rightChild of the parent of x tow
c)if the color of both the right and the leftChild of w is BLACK
case2:
a>set the color of w as RED
b>assign the parent of x to x
d)else if the color the rightChild of w is BLACK
case3:
a>set the color of the leftChild of w as BLACK
b>set the color of w as RED
c>right-rotate w
d>assign the rightChild of the parent of x to w
e)if any of the above cases do not occur, then do the following
case4:
a>set the color of w as the color of the parent of x.
b>set the color of the parent of x as BLACK
c>set the color of the right child of w as BLACK
d>left-rotate the parent of x
e>set x as the root of the tree
3.else the same as above with right changed to left and vice versa
4.set the color of x as BLACK