package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

type KeyPair struct {
	private *ecdsa.PrivateKey
	public  *ecdsa.PublicKey
	Address common.Address
}

type TrafficGenerator struct {
	KeyPairs []KeyPair
	client   *ethclient.Client
	Config   Config
}

func main() {
	app := cli.NewApp()
	app.Name = "da-traffic-generator"
	app.Usage = "EigenDA Traffic Generator"
	app.Description = "Service for generating traffic to EigenDA disperser"
	app.Flags = Flags
	app.Action = trafficGeneratorMain
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}

func trafficGeneratorMain(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	g, err := NewTrafficGenerator(config)
	if err != nil {
		fmt.Println("new traffic gen error", err)
	}

	ctx_, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	fmt.Println("g.Config.NumInstances", g.Config.NumInstances)

	// every instance will use its own wallet
	for i := 0; i < int(g.Config.NumInstances); i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			_ = g.StartTraffic(ctx_, j)
		}(i)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	cancel()
	wg.Wait()
	return nil
}

func NewTrafficGenerator(config *Config) (*TrafficGenerator, error) {
	chainClient, err := ethclient.Dial(config.Hostname)
	if err != nil {
		return nil, err
	}

	keyPairs := make([]KeyPair, len(config.SignerPrivateKeys))

	privateKeyList := config.SignerPrivateKeys
	Addresses := config.Addresses

	for i := 0; i < len(privateKeyList); i++ {
		var privateKey *ecdsa.PrivateKey
		privateKey, err = crypto.HexToECDSA(privateKeyList[i])
		if err != nil {
			return nil, fmt.Errorf("NewClient: cannot parse private key: %w", err)
		}
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot get public keys")
		}
		keyPairs[i].private = privateKey
		keyPairs[i].public = publicKeyECDSA
		keyPairs[i].Address = common.HexToAddress(Addresses[i])
	}

	return &TrafficGenerator{
		client:   chainClient,
		KeyPairs: keyPairs,
		Config:   *config,
	}, nil
}

func (g *TrafficGenerator) StartTraffic(ctx context.Context, fromIndex int) error {
	ticker := time.NewTicker(g.Config.RequestInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			toIndex := (fromIndex + 1) % int(g.Config.NumInstances)
			tx, err := g.CraftTx(context.Background(), g.KeyPairs[fromIndex], g.KeyPairs[toIndex])
			if err != nil {
				fmt.Println("failed to craft a tx", "err:", err)
			}

			err = g.SendRequest(ctx, tx)
			if err != nil {
				fmt.Println("failed to send blob request", "err:", err)
			}

			g.GetBalance(g.KeyPairs[fromIndex].Address)
		}
	}

	return nil
}

func (g *TrafficGenerator) GetBalance(address common.Address) {
	balance, err := g.client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address, "balance", balance)
}

func (g *TrafficGenerator) SendRequest(ctx context.Context, tx *types.Transaction) error {
	return g.client.SendTransaction(context.Background(), tx)

}

func (g *TrafficGenerator) CraftTx(ctx context.Context, from, to KeyPair) (*types.Transaction, error) {
	fromAddress := crypto.PubkeyToAddress(*from.public)
	nonce, err := g.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	value := big.NewInt(1000)           // in wei (1 eth)
	gasLimit := uint64(3000000)        // in units
	gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	toAddress := to.Address

	data := make([]byte, g.Config.PadSize)
	_, err = rand.Read(data)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := g.client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), from.private)
	if err != nil {
		return nil, err
	}
	binaryTx, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	fmt.Println("signedTx len after rlp encoding", len(binaryTx))

	return signedTx, nil
}
