package main

import (
    "bytes"
    "fmt"
    "strings"
)

func main() {
    // Efficient concatenation using strings.Builder
    var builder strings.Builder
    words := []string{"Hello", "World", "from", "Go"}
    for _, word := range words {
        builder.WriteString(word)
        builder.WriteString(" ")
    }
    result := builder.String()
    fmt.Println("Concatenated string using strings.Builder:", result)

    // Splitting and joining strings using strings package
    sentence := "Go is a great programming language"
    words = strings.Split(sentence, " ")
    fmt.Println("Split string:", words)

    rejoined := strings.Join(words, "-")
    fmt.Println("Joined string:", rejoined)

    // Efficient concatenation using bytes.Buffer
    var buffer bytes.Buffer
    for _, word := range words {
        buffer.WriteString(word)
        buffer.WriteString(" ")
    }
    result = buffer.String()
    fmt.Println("Concatenated string using bytes.Buffer:", result)

    // String formatting using fmt.Sprintf
    name := "Alice"
    age := 30
    formatted := fmt.Sprintf("Name: %s, Age: %d", name, age)
    fmt.Println("Formatted string using fmt.Sprintf:", formatted)

    // Handling Unicode strings with runes
    unicodeStr := "Hello, ✅"
    fmt.Print("Unicode string characters: ")
    for _, r := range unicodeStr {
        fmt.Printf("%c ", r)
    }
    fmt.Println()
}
