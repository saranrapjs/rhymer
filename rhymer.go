// Rhyming dictionary in go

package rhymer

import (
    "os"
    "bufio"
    "strings"
    "runtime"
    "path"
)

type phonTrie struct {
    leaves map[string]*phonTrie
    words map[string]bool
}

type rhymer struct {
    dictionary map[string][][]string
    trie phonTrie
}

// Create a new rhymer by reading the pronounciation dictionary
func NewRhymer() *rhymer {
    r := new(rhymer)

    // Read the file
    _, filename, _, _ := runtime.Caller(0)
    file, err := os.Open(path.Join(path.Dir(filename), "data", "reduxdict"))
    check(err)
    scanner := bufio.NewScanner(file)
    defer file.Close()
    check(err)

    r.dictionary = make(map[string][][]string)
    r.trie.leaves = make(map[string]*phonTrie)

    // Scan the file line by line
    for scanner.Scan() {
        // Split the line by whitespace
        f := strings.Fields(scanner.Text())
        // f[0] is the string, and f[1:] is the pronounciation
        r.dictionary[f[0]] = append(r.dictionary[f[0]], f[1:])

        // cur := r.trie
        // for i := len(f)-1; i >= 0; i-- {
            
        // }
    }

    return r
}

// Get the pronounciation phonemes for string s
func (r *rhymer) Pronounce(s string) [][]string {
    return r.dictionary[strings.ToUpper(s)]
}

// Find all words that rhyme with string s
func (r *rhymer) FindRhymesByWord(s string) []string {
    s = strings.ToUpper(s)
    if _, ok := r.dictionary[s]; !ok {
        return []string{}
    }
    return r.FindRhymes(r.dictionary[s][0])
}

// Find all words that rhyme with slice of pronounciation phonemes
func (r *rhymer) FindRhymes(s []string) []string {
    // Strip away the leading constanant sounds
    var words []string

    // Calculate the offset, if it's the error case, we return empty because there's no vowel sounds
    offset := vowelOffset(s)
    if offset == -1 {
        return []string{}
    }
    toRhyme := s[offset:]
    minLen := len(toRhyme)

    for k, v := range r.dictionary {
        for _, pronounce := range v {
            // Make sure that the word we compare against is at least as long as the word we rhyme
            if len(pronounce) < minLen {
                continue
            }
            // Check if they rhyme
            if rhymeTo(pronounce, toRhyme) {
                words = append(words, k)
            }
        }
    }
    return words
}

// Checks whether or not s1 and s2 rhyme, returns 1 if they do, 0 if they don't, and -1 when one word is unknown
func (r *rhymer) Rhymes(s1, s2 string) int {
    s1 = strings.ToUpper(s1)
    s2 = strings.ToUpper(s2)

    p1 := r.Pronounce(s1)
    p2 := r.Pronounce(s2)

    // Return -1 if one of the words is unknown
    if len(p1) == 0 || len(p2) == 0 {
        return -1
    }

    // Return 1 if any of the prounounciations rhyme
    for _, v := range p1 {
        for _, w := range p2 {
            if rhymeToUnordered(v, w) {
                return 1
            }
        }
    }
    return 0
}

// Find how many phonemes the vowel is offset so we know where to start rhyming
func vowelOffset(s []string) int {
    for i, v := range s {
        switch v[0] {
            case 'A', 'E', 'I', 'O', 'U':
                return i
        }
    }
    return -1
}

func rhymeTo(l, s []string) bool {
    diff := len(l) - len(s)
    ret := true

    offset := vowelOffset(s)

    // Check if the words sound the same, ignoring the first constanant sounds of the shorter word
    for i, v := range s[offset:] {
        if l[diff+i+offset] != v {
            ret = false
        }
    }

    return ret
}

func rhymeToUnordered(a1, a2 []string) bool {
    // Find the word with less rhymable phonemes
    var longer []string
    var shorter []string
    if len(a1)-vowelOffset(a1) > len(a2)-vowelOffset(a2) {
        longer = a1
        shorter = a2
    } else {
        longer = a2
        shorter = a1
    }

    return rhymeTo(longer, shorter)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}
