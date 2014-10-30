package vars

import (
    "math/big"
    "bytes"
    "github.com/eris-ltd/thelonious/ethstate"
    "github.com/eris-ltd/thelonious/ethutil"
    "github.com/eris-ltd/thelonious/ethcrypto"
)

var StdVarSize = 3

// location of a var is 1 followed by the first 8 bytes of its sha3
func VariName(name string) []byte{
    h := ethcrypto.Sha3Bin(ethutil.PackTxDataArgs(name))
    base := append([]byte{1}, h[:8]...)
    base = append(base, bytes.Repeat([]byte{0}, 32-len(base))...)
    return base
}

/*
    All interfaces are in strings, but array indices are ints.


*/


// Single Type Variable.
// (+ @@variname 1)
func GetSingle(addr, name string, state *ethstate.State) []byte {
    byteAddr := ethutil.Hex2Bytes(addr)
    obj := state.GetStateObject(byteAddr)
    base := VariName(name)
    base[31] = byte(StdVarSize+1)
    return (obj.GetStorage(ethutil.BigD(base))).Bytes()
}

func GetArrayIndex(addr, name string, index int, state *ethstate.State) []byte{
    return GetKeyedArrayIndex(addr, name, "0x0", index, state)
}

func GetKeyedArrayIndex(addr, name, key string, index int, state *ethstate.State) []byte {
    bigBase := big.NewInt(0)
    bigBase2 := big.NewInt(0)

    byteAddr := ethutil.Hex2Bytes(addr)
    obj := state.GetStateObject(byteAddr)
    base := VariName(name)

    // how big are the elements stored in this array:
    sizeLocator := make([]byte, len(base))
    copy(sizeLocator, base)
    sizeLocator = append(sizeLocator[:31], byte(StdVarSize+1))
    elementSizeBytes := (obj.GetStorage(ethutil.BigD(sizeLocator))).Bytes()
    elementSize := ethutil.BigD(elementSizeBytes).Uint64()

    // key should be trailing 20 bytes

    // what slot does the array start at:
    keyBytes := ethutil.UserHex2Bytes(key)
    keyBytesShift := append(keyBytes, []byte{1,0,0}...)
    slotBig := bigBase.Add(ethutil.BigD(base), ethutil.BigD(keyBytesShift))

    //numElements := obj.GetStorage(slotBig)

    // which slot (row), and where in that slot (col) is the element we want:
    entriesPerRow :=  int64(256 / elementSize)
    rowN := int64(index) / entriesPerRow
    colN := int64(index) % entriesPerRow

    row := bigBase.Add(big.NewInt(1), bigBase.Add(slotBig, big.NewInt(rowN))).Bytes()
    rowStorage := (obj.GetStorage(ethutil.BigD(row))).Bytes()
    rowStorageBig := ethutil.BigD(rowStorage)

    elSizeBig := ethutil.BigD(elementSizeBytes)
    // row storage gives us a big number, from which we need to pull
    // an element of size elementsize.
    // so divide it by 2^(colN*elSize) and take modulo 2^elsize 
    // divide row storage by 2^(colN*elSize)
    colBig := bigBase.Exp(big.NewInt(2), bigBase.Mul(elSizeBig, big.NewInt(colN)), nil)
    r := bigBase.Div(rowStorageBig, colBig)
    w := bigBase2.Exp(big.NewInt(2), elSizeBig, nil)
    v := bigBase.Mod(r, w)
    return v.Bytes()
}

func getLinkedList(addr, name, key string, state *ethstate.State){

}




