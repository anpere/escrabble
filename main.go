package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
  "sort"
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

/**
  * Given a dictionary body, generate the frequencies of each character
  * @param dictionary_corpus: a strong representing phrases in a dictionary
                              each word is on its own line
                              phrase declarations end with slashes or a newline
*/
func get_dictionary_frequencies(dictionary_corpus string) map[string]int {
	dict_frequencies := make(map[string]int)

	lines := strings.Split(dictionary_corpus, "\n")
	for _, line := range lines {
		word := strings.ToLower(strings.Split(line, "/")[0])
		for _, w_c := range word {
      char := string(w_c)
      if (!strings.ContainsAny(char, "0123456789\t\n\r'") && char!="") {
        dict_frequencies[char] += 1
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

	for char, freq := range frequency_map {
		fmt.Println("character:", char, "frequency", freq)
	}
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
		new_tree := tree_union(smaller, small)                   // this is a new tree
		rest := forest_difference(trees, []Tree{smaller, small}) // remove 2 leaves
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
  fmt.Println(smallest_tree, small_tree)
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
		codes[l_k] = "0 " + l_v
	}
	for r_k, r_v := range r_codes {
		codes[r_k] = "1 " + r_v
	}
	return codes
}

func main() {
	corpus_file := os.Args[1]
	corpus := get_corpus(corpus_file)
	hamming_code := generate_hamming_code(corpus)
	codes := hamming_code.GetCodes()

  var keys []string
  for k := range codes {
    keys = append(keys, k)
  }
  sort.Strings(keys)
	for _, word := range keys {
		fmt.Println("Word: ", word, "Code: ", codes[word])
	}

	fmt.Println("Hello, 世界")

}
