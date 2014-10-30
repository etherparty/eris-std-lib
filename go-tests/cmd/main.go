package main

import (
    "fmt"
    "github.com/eris-ltd/thelonious/monk"
    "github.com/eris-ltd/thelonious/ethutil"
    "github.com/eris-ltd/eris-std-lib/go-tests"
)

func main(){
   e := monk.NewEth(nil)
   e.ReadConfig("eth-config.json")
   e.Init() 
   e.Start() 

   addr := e.DeployContract("var-tests.lll", "lll")
   fmt.Println(addr)
   e.Commit()
   
   state := e.Pipe.World().State()

   s := vars.GetSingle(addr, "mySingle", state)
   fmt.Println(ethutil.Bytes2Hex(s))

   t := vars.GetArrayIndex(addr, "myArray", 2, state)
   fmt.Println(ethutil.Bytes2Hex(t))

   t = vars.GetArrayIndex(addr, "myArray", 5, state)
   fmt.Println(ethutil.Bytes2Hex(t))

   t = vars.GetArrayIndex(addr, "myArray", 6, state)
   fmt.Println(ethutil.Bytes2Hex(t))

}
