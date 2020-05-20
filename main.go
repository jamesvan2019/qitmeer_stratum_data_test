package main

import (
	"encoding/hex"
	"fmt"
	"github.com/Qitmeer/crypto/sha3"
	"golang.org/x/crypto/blake2b"
)

func main() {
	// jobID := 1
	prevHashReversed := "3ebbe1de1524c5d8450759652bfa0b1502f7b4b878b320dfcb56a3325b6307a3"
	cb1 := "0100010001"
	cb2 := "00000000000000000000000000000000000000000000000000000000000000005108"
	cb3 := "162f7777772e6d656572706f6f6c2e636f6d2f32303230"
	cb4 := "ffffffffffffffff01007841cb020000001976a914499896c7814a6f49fa256bc5feaa5882a665339188ac0000000000000000"
	zeroHash := "0000000000000000000000000000000000000000000000000000000000000000"
	merkleBranches := []string{}
	version := "0000000c"
	version = ReverseByWidth(version, 1)
	nbits := "1c1fffff"
	nbits = ReverseByWidth(nbits, 1)
	nTime := "5ec28a82"
	nTime = ReverseByWidth(nTime, 1)
	// mainHeight := 1
	// needCleanJob := true
	ex1 := "02000007"   // give by pool
	ex2 := "d5104dc7"   // random value
	nonce := "00000001" // nonce 1
	nonce = ReverseByWidth(nonce, 1)
	powType := "06" // 06 qitmeer_keccak256

	coinbaseHash := Blake2bd(cb1 + Blake2bd(cb2+ex1+ex2+cb3) + cb4)
	merkleHash := MakeMerkleRoot(coinbaseHash, merkleBranches)
	prevHash := prevHashReversed
	header := version + prevHash + merkleHash + zeroHash + nbits + nTime + nonce + powType
	fmt.Println(header)
	//output 0c0000003ebbe1de1524c5d8450759652bfa0b1502f7b4b878b320dfcb56a3325b6307a354b77ab9a3929854146abe12ba50c7a21022a245de2e8446a3bec5d6228e14de0000000000000000000000000000000000000000000000000000000000000000ffff1f1c828ac25e0100000006
	b, _ := hex.DecodeString(header)
	keccak := sha3.NewQitmeerKeccak256()
	keccak.Write(b)
	r := keccak.Sum(nil)
	// output c6acbc6536ad3b16024c96356ae199ec39c7754a153b7a40c99c2c500a09867a
	h := hex.EncodeToString(r)
	fmt.Println(h)
	if h == "c6acbc6536ad3b16024c96356ae199ec39c7754a153b7a40c99c2c500a09867a" {
		fmt.Println("check success")
		return
	}
	fmt.Println("check failed")
	// submit data
	//[submit]{PoolUser, jobID, ExtraNonce2, timestampStr,nonceStr}
	// params="[work01 1 d5104dc7 5ec28a82 00000001 ]"
}

//reverse LittleEndian bytes
func ReverseByWidth(str string, width int) string {
	s, _ := hex.DecodeString(str)
	newS := make([]byte, len(s))
	for i := 0; i < (len(s) / width); i += 1 {
		j := i * width
		copy(newS[len(s)-j-width:len(s)-j], s[j:j+width])
	}
	return hex.EncodeToString(newS)
}

// double blake2b
func Blake2bd(s string) string {
	b, _ := hex.DecodeString(s)
	h := blake2b.Sum256(b)
	h1 := blake2b.Sum256(h[:])
	return hex.EncodeToString(h1[:])
}

// Make MerkleRoot
func MakeMerkleRoot(coinbaseHash string, merkleBranches []string) string {
	for i := 0; i < len(merkleBranches); i++ {
		coinbaseHash += merkleBranches[i]
		coinbaseHash = Blake2bd(coinbaseHash)
	}
	return coinbaseHash
}
