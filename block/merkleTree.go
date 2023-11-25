package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

// Row represents a tier in the Merkle Tree
type Row []string

// MerkleTree represents the Merkle Tree
type MerkleTree struct {
	rows []Row
}

// MerklePathElement represents a hash value and its order in the Merkle Path
type MerklePathElement struct {
	hash                    string
	useFirstInConcatenation bool
}

// NewMerkleTree constructs a Merkle Tree from a list of transactions
func NewMerkleTree(transactions TransactionList) MerkleTree {
	var row Row
	for _, tx := range transactions.Transactions {
		row = append(row, tx.Hash)
	}

	tree := MerkleTree{
		rows: []Row{row},
	}

	for {
		rowAbove := tree.makeRowAbove(tree.rows[len(tree.rows)-1])
		tree.rows = append(tree.rows, rowAbove)
		tree.adjustRows() // Add the adjustment here
		if tree.isComplete() {
			break
		}
	}

	return tree
}

// MerkleRoot returns the Merkle Root of the tree
func (tree MerkleTree) MerkleRoot() string {
	return tree.rows[len(tree.rows)-1][0]
}

// MerklePathForLeaf returns the Merkle Path for a specific leaf in the tree
func (tree MerkleTree) MerklePathForLeaf(leafIndex int) MerklePath {
	var merklePath MerklePath
	i := leafIndex
	for _, row := range tree.rows[:len(tree.rows)-1] {
		sibling, useFirstInConcatenation := tree.evaluateSibling(row, i)
		merklePathElement := MerklePathElement{
			hash:                    row[sibling],
			useFirstInConcatenation: useFirstInConcatenation,
		}
		merklePath = append(merklePath, merklePathElement)
		i = i / 2
	}
	return merklePath
}

func (tree MerkleTree) evaluateSibling(row Row, myIndex int) (siblingIndex int, useFirstInConcatenation bool) {
	if myIndex%2 == 1 {
		siblingIndex = myIndex - 1
		useFirstInConcatenation = true
	} else if myIndex+1 <= len(row)-1 {
		siblingIndex = myIndex + 1
		useFirstInConcatenation = false
	} else {
		siblingIndex = myIndex
		useFirstInConcatenation = true
	}
	return
}

func (tree MerkleTree) adjustRows() {
	for level := len(tree.rows) - 1; level >= 0; level-- {
		row := tree.rows[level]
		//fmt.Println("Level", level, "has", len(row), "nodes")
		if len(row)%2 != 0 && len(row) > 1 {
			//fmt.Println("Adjusted row length to ensure even nodes in level", level)
			lastNode := row[len(row)-1]
			//fmt.Println("Duplicating node:", lastNode)
			tree.rows[level] = append(tree.rows[level], lastNode)
		}
	}
}
func (tree MerkleTree) makeRowAbove(below Row) Row {
	size := int(math.Ceil(float64(len(below)) / 2.0))
	row := make(Row, size)
	for i := range row {
		leftChild := i * 2
		rightChild := leftChild + 1
		if rightChild <= len(below)-1 {
			row[i] = tree.joinAndHash(below[leftChild], below[rightChild])
		} else {
			row[i] = tree.joinAndHash(below[leftChild], below[leftChild])
		}
	}
	return row
}

func (tree MerkleTree) joinAndHash(a, b string) string {
	data := strings.Join([]string{a, b}, "")
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (tree MerkleTree) isComplete() bool {
	return len(tree.rows[len(tree.rows)-1]) == 1
}

// MerklePath represents the path in the Merkle Tree
type MerklePath []MerklePathElement

func (tree MerkleTree) Display() {
	fmt.Println("Merkle Tree:")
	for i, row := range tree.rows {
		fmt.Printf("Level %d: %d nodes\n", i, len(row))
		fmt.Println(strings.Join(row, ", "))
	}
}
