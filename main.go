package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
  "sort"
  "strconv"
)

type Tree struct {
	left   *Tree
	right  *Tree
	key    string
	weight int
}

type Forest struct {
	trees []Tree
}
var IGNORE string = "ïkwàçèò0123456789-­­­­',.() "
var ACCENTS string = "áéïóú"
var TOO_RARE string = "wüïkyzx" + ACCENTS

var HANDLE_ACCENTS bool = false
func handle_accents(char string) string {
  if !HANDLE_ACCENTS {
    return char
  }
  if (char == "á") {
    return "a"
  } else if (char == "é") {
    return "e"
  } else if (char == "í" || char == "ï") {
    return "i"
  } else if (char == "ó") {
    return "o"
  } else if (char == "ú" || char == "ü") {
    return "u"
  } else {
    return char
  }
}
/**
  * Given a Spanish dictionary body, generate the frequencies of each character
  * @param dictionary_corpus: a string representing phrases in a dictionary
                              each phrase is on its own line
                              phrase declarations end with slashes or a newline
*/
func get_dictionary_frequencies(dictionary_corpus string) map[string]int {
	dict_frequencies := make(map[string]int)
  // Split the corpus into lines, and get the phrase in the beginning.
	lines := strings.Split(dictionary_corpus, "\n")
  var prev_char string = "";
	for _, line := range lines {
		word := strings.ToLower(strings.Split(line, "/")[0])
		for _, w_c := range word {
      char := string(w_c)
      if (!strings.ContainsAny(char, IGNORE) && len(char) > 0) {
        if (prev_char=="l" || prev_char=="r" ) && prev_char == char {
          dict_frequencies[char+char] += 1
          fmt.Println(char+char)
          dict_frequencies[char] -= 1 // let's not count the previous instance
        } else if (prev_char == "c") && (char == "h") {
          dict_frequencies["ch"] += 1
          dict_frequencies["c"] -= 1
        }
        prev_char = char
        dict_frequencies[handle_accents(char)] += 1
      }
		}
	}
	return dict_frequencies
}

/**
 * Given the name of the file, gets the contents of the file
 */
func get_corpus(file_name string) string {
	corpus, err := ioutil.ReadFile(file_name)
	if err != nil {
		log.Fatal(err)
	}
	return string(corpus)
}

/**
 * Given a corpus string, generates a coding from letters in the corpus
 * to the codeword it should be
 *
 */
func generate_hamming_code(corpus string) Tree {

	frequency_map := get_dictionary_frequencies(corpus)
	// Go through the frequencies and generate the hamming
	tree_freqs := generate_forest(frequency_map)
	ham_tree := generate_ham_tree(tree_freqs)
	return ham_tree
}

func generate_ham_tree(trees []Tree) Tree {
	if len(trees) == 1 {
		// This is the base case where we have already generateed the ham tree
		return trees[0]
	} else {
		// Find the smallest trees
		smaller, small := find_two_mins(trees)
		// Coalesce the tiny trees
		new_tree := tree_union(smaller, small)  // this is a new tree
		rest := forest_difference(trees, []Tree{smaller, small})  // remove 2 leaves
		return generate_ham_tree(append(rest, new_tree))
	}
}
func find_two_mins(trees []Tree) (Tree, Tree) {
	var smallest_value int = math.MaxInt64
	var smallest_tree, small_tree Tree
	for _, tree := range trees {
		if tree.weight <= smallest_value {
			small_tree = smallest_tree
			smallest_value = tree.weight
			smallest_tree = tree
		}
	}
	return smallest_tree, small_tree
}
func tree_union(left Tree, right Tree) Tree {
	return Tree{
		&left,
		&right,
		"",
		left.weight + right.weight,
	}
}
func forest_difference(left []Tree, right []Tree) []Tree {
	var diff []Tree
	right_temp := map[Tree]bool{}
	for _, t_r := range right {
		right_temp[t_r] = true
	}
	for _, t_l := range left {
		if _, ok := right_temp[t_l]; !ok {
			diff = append(diff, t_l)
		}
	}
	return diff
}

func generate_forest(weight_maps map[string]int) []Tree {
	var forest []Tree
	for key, value := range weight_maps {
		leaf := Tree{
			nil,
			nil,
			key,
			value,
		}
		forest = append(forest, leaf)
	}
	return forest
}

func (t *Tree) GetCodes() map[string]string {
	if t.left == nil || t.right == nil {
		return map[string]string{t.key: ""}
	}
	l_codes := t.left.GetCodes()
	r_codes := t.right.GetCodes()
	codes := make(map[string]string)
	for l_k, l_v := range l_codes {
		codes[l_k] = "0" + l_v
	}
	for r_k, r_v := range r_codes {
		codes[r_k] = "1" + r_v
	}
  delete(codes, "")
	return codes
}


func check(e error) {
  if e != nil {
      panic(e)
  }
}

func gen_piece_count(mode string) {
  var piece_count map[string]int
  var msg string
  if mode=="ES" {
    piece_count = map[string]int{
      "A":12, "E" :12, "O" :9, "I" :6, "S" :12, "N" :10, "R" :10,"U" :10, "L" :8,
      "T" :8,"D" :10,"G" :4, "C" :8,"B" :4,"M" :4,"P" :4,"H" :4,"F" :2,"V" :2,
      "Y" :2,"CH" :2,"Q" :2,"J" :2,"LL" :2,"Ñ" :2,"RR" :2,"X" :2, "Z":2,
      "á":12 , "é":12 , "ó":9, "ü":2, "ú":5, "í":6,
    }
  } else if mode=="EN" {
    piece_count   = map[string]int{
      "a":9,"b":2,"c":2,"d":4,"e":12,"f":2,"g":3,"h":2,"i":9,"j":1,"k":1,"l":4,
      "m":2,"n":6,"o":8,"p":2,"q":1,"r":6,"s":4,"t":6,"u":4,"v":2,"w":2,"x":1,
      "y":2,"z":1,
    }
  }
  for letter, count := range piece_count {
    msg += (strings.ToLower(letter) + ":" + strconv.Itoa(count)+"\n")
  }
  err := ioutil.WriteFile("piece_count.txt", []byte(msg), 0644)
  check(err)
}
func main() {
	corpus_file := os.Args[1]
  var mode string
  if strings.Contains(corpus_file, "Spanish"){
      mode = "ES"
  } else {
    mode = "EN"
  }

	corpus := get_corpus(corpus_file)
	hamming_code := generate_hamming_code(corpus)
  frequencies := get_dictionary_frequencies(corpus)
  var min int = math.MaxInt64
  for char, freq := range frequencies {
    if strings.ContainsAny(char, TOO_RARE) {
      continue
    }
    if freq < min {
      min = freq
    }
  }
  character_count := make(map[string]int)
  for char, freq := range frequencies {
    character_count[char] = 12 * freq / 85080
  }
  //fmt.Println(character_count)
	codes := hamming_code.GetCodes()

  var keys []string
  for k := range codes {
    keys = append(keys, k)
  }
  sort.Strings(keys)
  var msg = ""
	for _, w := range keys {
    msg += (string(w) + ":" + strconv.Itoa(len(codes[w])) + "\n")
		fmt.Println(w,":",len(codes[w]))
	}

  gen_piece_count(mode)

  err := ioutil.WriteFile("piece_value.txt", []byte(msg), 0644)

  check(err)
}
