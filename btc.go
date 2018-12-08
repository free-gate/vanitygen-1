package main

import (
	"strings"
	"time"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	log "github.com/sirupsen/logrus"
)

type BtcCmd struct {
	Network string `long:"network" default:"mainnet" description:"network"`
	Pattern string `long:"pattern" default:"1kkk" description:"match pattern"`
}

func (bc *BtcCmd) Execute(args []string) ([]string ,error) {
	if len(args) >= 1 {
		bc.Pattern = args[0]
	}
	if len(bc.Pattern) > 6 {
		//log.Fatal("Quitting, this pattern would take too much time.")
		log.Info("Quitting, this pattern would take too much time.")
	} else if len(bc.Pattern) > 4 {
		log.Info("This pattern could take awhile, please wait.")
	}

	beginTime := time.Now()
	//prefix := strings.ToLower(bc.Pattern)
	prefix := bc.Pattern
	chainParams := &chaincfg.MainNetParams
	
	fmt.Println("*****", bc, bc.Pattern, args[0])

	switch bc.Network {
	case "mainnet":
		chainParams = &chaincfg.MainNetParams
	case "testnet":
		chainParams = &chaincfg.TestNet3Params
	default:
		chainParams = &chaincfg.MainNetParams
	}

	var numAttempts int64 = 0
	foundAddr := ""
	foundWif := ""
	for {
		numAttempts++

		privKey, err := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			log.Fatalf("Failed to create private key, err: %v", err)
		}

		addrPubKey, err := btcutil.NewAddressPubKey(
			privKey.PubKey().SerializeUncompressed(), chainParams)
		if err != nil {
			log.Fatalf("Failed to calculate public key, err: %v", err)
		}

		rcvAddr := addrPubKey.AddressPubKeyHash().EncodeAddress()
		if matchPrefix(rcvAddr, prefix) {
			foundAddr = rcvAddr
			wif, err := btcutil.NewWIF(privKey, chainParams, false)
			if err != nil {
				log.Fatalf("failed to get wif: %s", err)
			}
			foundWif = wif.String()
			fmt.Println("privKey:",privKey, "\naddrPubKey:", addrPubKey)
			break
		}
	}
	
	var result []string
	//result = append(result, time.Since(beginTime))
	result = append(result, foundAddr)
	result = append(result, foundWif)

	log.Infof("\nElapsed: %s\naddr: %s\nwif: %s\nattempts: %d.\n",
		time.Since(beginTime), foundAddr, foundWif, numAttempts)

	return result, nil
}

// Case-insensitive otherwise search performance suffers
func matchPrefix(address string, prefix string) bool {
	// compare search pattern to the left-most substr
	//lower := strings.ToLower(address)
	//fmt.Println("lower:", lower, "address:", address, "prefix:", prefix)
	//return strings.HasPrefix(lower, prefix)
	return strings.HasPrefix(address, prefix)
}

var btcCmd BtcCmd

func init() {
	parser.AddCommand("btc", "get a BTC vanity address", "The ticker command get a BTC vanity address", &btcCmd)
//	fmt.Println("btc.go init:",btcCmd)
}
