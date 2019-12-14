package fly

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

// StreamParsing read the Json in chunks
func streamBParsing(input *os.File, c chan flypayBTransaction) {

	// Error Handling
	he := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	// Start reading the file
	dec := json.NewDecoder(bufio.NewReader(input))

	// Start of the expect object
	t, err := dec.Token()
	he(err)

	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		log.Fatal("Expected object")
	}

	// Read props
	for dec.More() {
		t, err = dec.Token()
		he(err)

		// transactions array
		if t == "transactions" {
			t, err := dec.Token()
			he(err)
			if delim, ok := t.(json.Delim); !ok || delim != '[' {
				log.Fatal("Expected array")
			}

			// Read transactions
			for dec.More() {
				// Read next transaction
				ft := flypayBTransaction{}
				he(dec.Decode(&ft))
				c <- ft
			}

			// Array closing delim
			t, err = dec.Token()
			he(err)
			if delim, ok := t.(json.Delim); !ok || delim != ']' {
				log.Fatal("Expected array closing")
			}
		}
	}

	// Object closing delim
	t, err = dec.Token()
	he(err)
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		log.Fatal("Expected object closing")
	}

	close(c)
}
