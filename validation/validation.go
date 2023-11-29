package validation

import (
	"encoding/hex"
	"strings"
)

func GenerateTargetHash(difficultyTarget int) (string, []byte) {
	// returns a 256 bit hash
	// difficultyTarget is the number of leading zeros
	// hexadecimal string representation of the hash will contain 64 characters
	// 4 bits per character
	if difficultyTarget < 0 {
		difficultyTarget = 0
	} else if difficultyTarget > 256 {
		difficultyTarget = 256
	}
	totalCharacters := 64
	numZeros := difficultyTarget / 4
	leftover := difficultyTarget % 4

	// create a hex string with the required number of zeros
	// and the required number of leftover characters
	switch leftover {
	case 0:
		targetString := strings.Repeat("0", numZeros) + strings.Repeat("f", totalCharacters-numZeros)
		// convert the hex string to a byte array
		targetBytes, _ := hex.DecodeString(targetString)
		return targetString, targetBytes
	case 1:
		targetString := strings.Repeat("0", numZeros) + "7" + strings.Repeat("f", totalCharacters-numZeros-1)
		// convert the hex string to a byte array
		targetBytes, _ := hex.DecodeString(targetString)
		return targetString, targetBytes
	case 2:
		targetString := strings.Repeat("0", numZeros) + "3" + strings.Repeat("f", totalCharacters-numZeros-1)
		// convert the hex string to a byte array
		targetBytes, _ := hex.DecodeString(targetString)
		return targetString, targetBytes
	case 3:
		targetString := strings.Repeat("0", numZeros) + "1" + strings.Repeat("f", totalCharacters-numZeros-1)
		// convert the hex string to a byte array
		targetBytes, _ := hex.DecodeString(targetString)
		return targetString, targetBytes
	}
	return "", nil
}

// returns true if the hash is less than or equal to the target
// returns false if the hash is greater than the target
func CompareWithTargetHash(hash []byte, target []byte) bool {
	// must be the same length
	if len(hash) != len(target) || len(hash) != 32 {
		// fmt.Println("Hash and target are not the same length")
		return false
	}
	for i := 0; i < len(hash); i++ {
		if hash[i] < target[i] {
			return true
		} else if hash[i] > target[i] {
			return false
		}
	}
	return true
}
