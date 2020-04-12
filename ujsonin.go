package main

import (
    "fmt"
    "io/ioutil"
    "strings"
)

func main() {
    content, _ := ioutil.ReadFile("test.json")
    root := parse( content )
    DumpNodeHash( root, 1 )
}

type NodeHash struct {
    hash map [ string ] *JNode
}

type NodeStr struct {
    str string
}

type JNode struct {
    parent   *JNode
    nodeType int
    hash     *NodeHash
    str      *NodeStr
}

func NewNodeHash() (* NodeHash) {
    node := &NodeHash{ hash: make( map [ string ] *JNode ) }
    return node
}

func NewNodeStr( str string ) ( *NodeStr ) {
    return &NodeStr{ str: str }
}

func DumpNodeHash( self *NodeHash, depth int ) {
    fmt.Printf("{\n")
    
    for key, val := range self.hash {
        fmt.Printf("%s\"%s\":",strings.Repeat("  ",depth), key);
        DumpJNode( val, depth );
    }
    depth--;
    fmt.Printf("%s}\n",strings.Repeat("  ",depth))
}

func DumpJNode( self *JNode, depth int ) {
    if self.nodeType == 1 {
        DumpNodeHash( self.hash, depth+1 )
    } else if self.nodeType == 2 {
        fmt.Printf("\"%s\"\n", self.str.str )
    }
}

func parse( data [] byte ) (*NodeHash) {
    size := len( data )
    
    pos := 1
    keyStart := 0
    strStart := 0
    key := ""
    var let byte
    
    root := NewNodeHash()
    cur := &JNode{ nodeType: 1, hash: root }
Hash:
    let = data[pos]
    pos++
    if pos >= size { goto Done } 
    if let == '"' {
        goto KeyName1
    }
    if let == '}' && cur.parent != nil {
        cur = cur.parent
    }
    if let == '/' && pos < (size-1) {
        if data[pos] == '/' {
            pos++
            goto HashComment
        }
        if data[pos] == '*' {
            pos++
            goto HashComment2
        }
    }
    goto Hash
HashComment:
    let = data[pos]
    pos++
    if let == 0x0d || let == 0x0a {
        goto Hash
    }
    goto HashComment
HashComment2:
    let = data[pos]
    pos++
    if( let == '*' && pos < (size-1) && data[pos] == '/' ) { pos++; goto Hash; }
    goto HashComment2;
KeyName1:
    keyStart = pos
    pos++
KeyNameX:
    let = data[pos]
    pos++
    if let == '"' {
        key = string( data[ keyStart : pos - 1 ] )
        goto Colon
    }
    goto KeyNameX
Colon:
    let = data[pos]
    pos++
    if let == ':' {
        goto AfterColon
    }
    goto Colon
AfterColon:
    let = data[pos]
    pos++
    if let == '"' {
        goto String1
    }
    if let == '{' {
        newHash := NewNodeHash()
        newJNode := &JNode{ nodeType: 1, hash: newHash }
        newJNode.parent = cur;
        if cur.nodeType == 1 {
            cur.hash.hash[ key ] = newJNode
        }
        cur = newJNode
        goto Hash
    }
    if let == '/' && pos < (size-1) {
        if data[pos] == '/' {
            pos++
            goto AC_Comment
        }
        if data[pos] == '*' {
            pos++
            goto AC_Comment2
        }
    }
    // if( let == 't' || let == 'f' ) ... for true/false
    // if( let >= '0' && let <= '9' ) ... for numbers
    // if( let == '[' ) ... for array
    goto AfterColon;
AC_Comment:
    let = data[pos]
    pos++
    if let == 0x0d || let == 0x0a {
        goto AfterColon
    }
    goto AC_Comment
AC_Comment2:
    let = data[pos]
    pos++
    if let == '*' && pos < (size-1) && data[pos] == '/' {
        pos++
        goto Hash
    }
    goto AC_Comment2
String1:
    let = data[pos]
    pos++
    if( let == '"' ) {
       nodeStr := NewNodeStr( "" )
       newJNode := &JNode{ nodeType: 2, str: nodeStr } 
       cur.hash.hash[ key ] = newJNode
       goto AfterVal;
    }
    strStart = pos - 1
StringX:
    let = data[pos]
    pos++
    if let == '"' {
       if cur.nodeType == 1 {
           newStr := NewNodeStr( string( data[ strStart : pos - 1 ] ) )
           newJNode := &JNode{ nodeType: 2, str: newStr }
           cur.hash.hash[ key ] = newJNode
       }
       goto AfterVal
    }
    goto StringX;   
AfterVal:
    // who cares about commas in between things; we can just ignore them :D
    goto Hash
Done:
    return root
}