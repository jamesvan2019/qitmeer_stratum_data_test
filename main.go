package main

import (
	"encoding/hex"
	"fmt"
	"github.com/Qitmeer/qitmeer/common"
	"github.com/Qitmeer/qitmeer/common/hash"
	"golang.org/x/crypto/blake2b"
	"math/big"
)

var limitBig = new(big.Int).Sub(new(big.Int).Lsh(common.Big1, 224), common.Big1)

func main() {
	// jobID := 1
	prevHashReversed := "a75ba110c3bdedb4b577e38d31bd171aff548cb6441e7dc25cf2fe9d537b9b11"
	cb1 := "0100010001"
	cb2 := "000000000000000000000000000000000000000000000000000000000000000003f1bd0008"
	cb3 := "162f7777772e6d656572706f6f6c2e636f6d2f32303230"
	cb4 := "ffffffffffffffff01007841cb020000001976a914499896c7814a6f49fa256bc5feaa5882a665339188ac0000000000000000"
	zeroHash := "0000000000000000000000000000000000000000000000000000000000000000"
	merkleBranches := []string{}
	version := "0000000c"
	version = ReverseByWidth(version, 1)
	nbits := "1b04f3e6"
	nbits = ReverseByWidth(nbits, 1)
	nTime := "5ee76f8e"
	nTime = ReverseByWidth(nTime, 1)
	// mainHeight := 1
	// needCleanJob := true
	ex1 := "0200682d"   // give by pool
	ex2 := "00000000"   // random value
	nonce := "c3913447" // nonce
	nonce = ReverseByWidth(nonce, 1)
	powType := "06" // 06 qitmeer_keccak256

	coinbaseHash := Blake2bd(cb1 + Blake2bd(cb2+ex1+ex2+cb3) + cb4)
	merkleHash := MakeMerkleRoot(coinbaseHash, merkleBranches)
	prevHash := prevHashReversed
	header := version + prevHash + merkleHash + zeroHash + nbits + nTime + nonce + powType
	fmt.Println(header)
	//output 0c000000a75ba110c3bdedb4b577e38d31bd171aff548cb6441e7dc25cf2fe9d537b9b1103028984f09f993c478c628d2b12334e7e1f65b0bba02e162ac92ed36340fda00000000000000000000000000000000000000000000000000000000000000000e6f3041b8e6fe75e473491c306
	b, _ := hex.DecodeString(header)
	// keccak := sha3.NewQitmeerKeccak256()
	// keccak.Write(b)
	// r := keccak.Sum(nil)
	// // output 87d41e1256914d78ffd10e1a8e3f0d6c2bf60cf0bb5f435de70f280100000000
	// h := hex.EncodeToString(r)
	h := hash.HashQitmeerKeccak256(b)
	target := limitBig.Div(limitBig, big.NewInt(50))
	fmt.Println(h)
	fmt.Println(fmt.Sprintf("%064x", target))
	//target 00000000051eb851eb851eb851eb851eb851eb851eb851eb851eb851eb851eb8
	if HashToBig(&h).Cmp(target) <= 0 {
		fmt.Println("match difficulty")
		return
	}
	fmt.Println("not match difficulty")
	if h.String() == "87d41e1256914d78ffd10e1a8e3f0d6c2bf60cf0bb5f435de70f280100000000" {
		fmt.Println("check success")
		return
	}
	fmt.Println("check failed")
	// submit data
	//[submit]{PoolUser, jobID, ExtraNonce2, timestampStr,nonceStr}
	// params="[work01 1 00000000 5ee76f8e c3913447 ]"
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
func HashToBig(hash *hash.Hash) *big.Int {
	// A Hash is in little-endian, but the big package wants the bytes in
	// big-endian, so reverse them.
	buf := *hash
	blen := len(buf)
	for i := 0; i < blen/2; i++ {
		buf[i], buf[blen-1-i] = buf[blen-1-i], buf[i]
	}

	return new(big.Int).SetBytes(buf[:])
}

// Make MerkleRoot
func MakeMerkleRoot(coinbaseHash string, merkleBranches []string) string {
	for i := 0; i < len(merkleBranches); i++ {
		coinbaseHash += merkleBranches[i]
		coinbaseHash = Blake2bd(coinbaseHash)
	}
	return coinbaseHash
}
