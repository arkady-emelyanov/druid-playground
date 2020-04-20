package main

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/pkg/namesgenerator"
	"os"
)

type Merchant struct {
	MerchantId uint
	Name       string
}

func main() {
	merchantChan := make(chan *Merchant, 1000)
	writeDoneChan := make(chan struct{}, 0)

	// receive transactions and put all of them to file
	go func(merchantChan <-chan *Merchant) {
		fp, _ := os.Create("generated/merchants.json")
		defer fp.Close()

		i := 0
		for tr := range merchantChan {
			if i%10000 == 0 {
				fmt.Println("Syncing file to disk, written:", i)
				fp.Sync()
			}

			b, err := json.Marshal(tr)
			if err != nil {
				panic(err)
			}

			if _, err := fp.Write(b); err != nil {
				panic(err)
			}
			if _, err := fp.Write([]byte("\n")); err != nil {
				panic(err)
			}
			i++
		}

		writeDoneChan <- struct{}{}
	}(merchantChan)

	for i := 0; i < 1001; i++ {
		name := namesgenerator.GetRandomName(0)
		merchantChan <- &Merchant{
			MerchantId: uint(i),
			Name:       name,
		}
	}

	close(merchantChan)
	<-writeDoneChan
}
