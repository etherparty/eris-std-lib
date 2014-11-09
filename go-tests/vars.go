package vars

import (
    "math/big"
    "bytes"
    "github.com/eris-ltd/thelonious/monkstate"
    "github.com/eris-ltd/thelonious/monkutil"
    "github.com/eris-ltd/thelonious/monkcrypto"
)

// each var comes with 3 permissions: add, rm, mod
var StdVarSize = 4

// location of a var is 1 followed by the first 8 bytes of its sha3
func VariName(name string) []byte{
    h := monkcrypto.Sha3Bin(monkutil.PackTxDataArgs2(name))
    base := append([]byte{1}, h[:8]...)
    base = append(base, bytes.Repeat([]byte{0}, 32-len(base))...)
    return base
}

/*
    All interfaces are in strings, but array indices are ints.
*/


// Single Type Variable.
// (+ @@variname 1)
func GetSingle(addr []byte , name string, state *monkstate.State) []byte {
    obj := state.GetStateObject(addr)
    base := VariName(name)
    base[31] = byte(StdVarSize+1)
    return (obj.GetStorage(monkutil.BigD(base))).Bytes()
}

func GetArrayElement(addr []byte, name string, index int, state *monkstate.State) []byte{
    return GetKeyedArrayElement(addr, name, "0x0", index, state)
}

func GetArray(addr []byte, name string, state *monkstate.State) [][]byte{
    return GetKeyedArray(addr, name, "0x0", state)
}

// key must come as hex!!!!!!
func GetKeyedArrayElement(addr []byte, name, key string, index int, state *monkstate.State) []byte {
    bigBase := big.NewInt(0)
    bigBase2 := big.NewInt(0)

    obj := state.GetStateObject(addr)
    base := VariName(name)

    // how big are the elements stored in this array:
    sizeLocator := make([]byte, len(base))
    copy(sizeLocator, base)
    sizeLocator = append(sizeLocator[:31], byte(StdVarSize+1))
    elementSizeBytes := (obj.GetStorage(monkutil.BigD(sizeLocator))).Bytes()
    elementSize := monkutil.BigD(elementSizeBytes).Uint64()

    // key should be trailing 20 bytes
    if len(key) >= 2 && key[:2] == "0x"{
        key = key[2:]
    }
    if l := len(key); l > 40{
        key = key[l-40:]
    }

    // what slot does the array start at:
    keyBytes := monkutil.PackTxDataArgs2("0x"+key)
    keyBytesShift := append(keyBytes[3:], []byte{1,0,0}...)
    slotBig := bigBase.Add(monkutil.BigD(base), monkutil.BigD(keyBytesShift))

    //numElements := obj.GetStorage(slotBig)

    // which slot (row), and where in that slot (col) is the element we want:
    entriesPerRow :=  int64(256 / elementSize)
    rowN := int64(index) / entriesPerRow
    colN := int64(index) % entriesPerRow

    row := bigBase.Add(big.NewInt(1), bigBase.Add(slotBig, big.NewInt(rowN))).Bytes()
    rowStorage := (obj.GetStorage(monkutil.BigD(row))).Bytes()
    rowStorageBig := monkutil.BigD(rowStorage)

    elSizeBig := monkutil.BigD(elementSizeBytes)
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

func GetKeyedArray(addr []byte, name, key string, state *monkstate.State) [][]byte{
    return nil
}

func GetLinkedListElement(addr []byte, name, key string, state *monkstate.State) []byte{
    bigBase := big.NewInt(0)

    obj := state.GetStateObject(addr)
    base := VariName(name)

    // key should be trailing 20 bytes
    if l := len(key); l > 20{
        key = key[l-20:]
    }

    // get slot for this keyed element of linked list
    keyBytes := monkutil.PackTxDataArgs2(key)
    keyBytesShift := append(keyBytes, []byte{1,0,0}...)[3:]
    slotBig := bigBase.Add(monkutil.BigD(base), monkutil.BigD(keyBytesShift))

    // value is right at slot    
    v := obj.GetStorage(slotBig)          
    return v.Bytes()
}

