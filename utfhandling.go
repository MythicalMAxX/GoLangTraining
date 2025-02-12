package main

import (
    "bytes"
    "fmt"
    "unicode/utf8"
)

func main() {
    // Original string
    str := "Hello, 🤡✅"

    // Converting string to bytes
    byteSlice := []byte(str)
    fmt.Println("Byte Slice:", byteSlice)

    // Converting string to runes
    runeSlice := []rune(str)
    fmt.Println("Rune Slice:", runeSlice)

    // Iterating over bytes
    fmt.Print("Bytes: ")
    for _, b := range byteSlice {
        fmt.Printf("%d ", b)
    }
    fmt.Println()

    // Iterating over runes
    fmt.Print("Runes: ")
    for _, r := range runeSlice {
        fmt.Printf("%c ", r)
    }
    fmt.Println()

    // UTF-8 Handling: Counting runes
    runeCount := utf8.RuneCountInString(str)
    fmt.Println("Number of runes:", runeCount)

    // UTF-8 Handling: Iterating over runes in string
    fmt.Print("UTF-8 Runes: ")
    for i, r := range str {
        fmt.Printf("%c (at position %d) ", r, i)
    }
    fmt.Println()

    // Efficient concatenation using bytes.Buffer
    var buffer bytes.Buffer
    words := []string{"Hello", "World", "from", "Go"}
    for _, word := range words {
        buffer.WriteString(word)
        buffer.WriteString(" ")
    }
    result := buffer.String()
    fmt.Println("Concatenated string using bytes.Buffer:", result)

    // String manipulation using rune slices
    modifiedRunes := []rune(str)
    modifiedRunes[7] = '🤡'
    modifiedRunes[8] = '✅'
    modifiedStr := string(modifiedRunes)
    fmt.Println("Modified String:", modifiedStr)
}
