package main

import (
    "fmt"
    "github.com/eris-ltd/thelonious/monk"
    "github.com/eris-ltd/thelonious/monkchain"
    "github.com/eris-ltd/thelonious/monkutil"
    "github.com/eris-ltd/eris-std-lib/go-tests"
)

func main(){
    monkchain.NoGenDoug = true
    e := monk.NewEth(nil)
    e.ReadConfig("eth-config.json")
    e.Init() 
    e.Start() 

    addr := e.DeployContract("var-tests.lll", "lll")
    fmt.Println(addr)

    e.Commit()

    state := e.Pipe.World().State()

    // test single
    s := vars.GetSingle(addr, "mySingle", state)
    fmt.Println(monkutil.Bytes2Hex(s))

    // test array
    t := vars.GetArrayElement(addr, "myArray", 2, state)
    fmt.Println(monkutil.Bytes2Hex(t))
    t = vars.GetArrayElement(addr, "myArray", 5, state)
    fmt.Println(monkutil.Bytes2Hex(t))
    t = vars.GetArrayElement(addr, "myArray", 6, state)
    fmt.Println(monkutil.Bytes2Hex(t))

    // test linked list
    l := vars.GetLinkedListElement(addr, "myLL", "balls", state)
    fmt.Println(monkutil.Bytes2Hex(l))
    l = vars.GetLinkedListElement(addr, "myLL", "paws", state)
    fmt.Println(monkutil.Bytes2Hex(l))
    
}
