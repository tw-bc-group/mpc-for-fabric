package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	"github.com/LabZion/HEaaS/common"
	pb "github.com/LabZion/HEaaS/fhe"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/ldsec/lattigo/bfv"
	"google.golang.org/grpc"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "172.26.32.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

var account = "fan@torchz.net"

// KeyPair is a pair of bfv public and private keys
type KeyPair struct {
	PublicKey []byte
	SecretKey []byte
}

// Bid is a bid
type Bid struct {
	LimitPriceDistance int
	CreditDistance     int
}

// generateKeysRemote gets a new pair of fhe keys
func generateKeysRemote(client pb.FHEClient) KeyPair {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keyPair, err := client.GenerateKey(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("%v.GenerateKey(_) = _, %v: ", client, err)
	}
	return KeyPair{
		PublicKey: keyPair.PublicKey,
		SecretKey: keyPair.SecretKey,
	}
}

// storeKey store a pair of fhe keys
func storeKey(client pb.FHEClient, account string, keyPair KeyPair) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.StoreKey(ctx, &pb.StoreKeyRequest{
		Account: account,
		KeyPair: &pb.KeyPair{
			PublicKey: keyPair.PublicKey,
			SecretKey: keyPair.SecretKey,
		},
	})
	if err != nil {
		log.Fatalf("%v.StoreKey(_) = _, %v: ", client, err)
	}
	return
}

// setAsk set an ask for account
func setAsk(client pb.FHEClient, keyPair KeyPair, account string, limit int) {
	params := common.GetParams()

	sk := bfv.SecretKey{}
	_ = sk.UnmarshalBinary(keyPair.SecretKey)
	encryptorSk := bfv.NewEncryptorFromSk(params, &sk)

	limitCiphertextBytes := common.EncryptInt(encryptorSk, limit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.SetAsk(ctx, &pb.AskRequest{
		Account:              account,
		LimitPriceCipherText: limitCiphertextBytes,
	})
	if err != nil {
		log.Fatalf("%v.SetAsk(_) = _, %v: ", client, err)
	}
	return
}

// setBid set an bid for account
func setBid(client pb.FHEClient, keyPair KeyPair, targetAccount string, account string, limit int, credit int) {
	params := common.GetParams()

	pk := bfv.PublicKey{}
	pk.UnmarshalBinary(keyPair.PublicKey)
	encryptorPk := bfv.NewEncryptorFromPk(params, &pk)

	limitCiphertextBytes := common.EncryptInt(encryptorPk, limit)
	creditCiphertextBytes := common.EncryptInt(encryptorPk, credit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.SetBid(ctx, &pb.BidRequest{
		TargetAccount:        targetAccount,
		Account:              account,
		LimitPriceCipherText: limitCiphertextBytes,
		CreditCipherText:     creditCiphertextBytes,
	})
	if err != nil {
		log.Fatalf("%v.SetBid(_) = _, %v: ", client, err)
	}
	return
}

// getEligibleBids fetch all eligible bids
func getEligibleBids(client pb.FHEClient, keyPair KeyPair, account string) (uint64, []Bid) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eligibleBidResponse, err := client.EligibleBid(ctx, &pb.EligibleBidRequest{
		Account: account,
	})
	if err != nil {
		log.Fatalf("%v.EligibleBid(_) = _, %v: ", client, err)
	}
	// Decrypting Bids
	bids := []Bid{}
	params := common.GetParams()

	sk := bfv.SecretKey{}
	_ = sk.UnmarshalBinary(keyPair.SecretKey)
	decryptor := bfv.NewDecryptor(params, &sk)

	for _, bid := range eligibleBidResponse.Bids {
		limitPriceDistance := common.DecryptInt(decryptor, bid.LimitPriceDistanceCiphertext)
		creditDistance := common.DecryptInt(decryptor, bid.CreditDistanceCiphertext)

		bids = append(bids, Bid{
			LimitPriceDistance: limitPriceDistance,
			CreditDistance:     creditDistance,
		})
	}
	return eligibleBidResponse.TotalBidNumber, bids
}

func fetchPublicKey(client pb.FHEClient, account string) KeyPair {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keyPair, err := client.FetchPublicKey(ctx, &pb.FetchPublicKeyRequest{
		Account: account,
	})
	if err != nil {
		log.Fatalf("%v.FetchPublicKey(_) = _, %v: ", client, err)
	}
	return KeyPair{
		PublicKey: keyPair.PublicKey,
		SecretKey: keyPair.SecretKey,
	}
}

type SmartContract struct {
	contractapi.Contract
}

func newClient() pb.FHEClient {
	var opts []grpc.DialOption

	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := pb.NewFHEClient(conn)

	return client
}

func (s *SmartContract) SetAsk(_ contractapi.TransactionContextInterface) (string, error) {
	client := newClient()
	kp := generateKeysRemote(client)
	storeKey(client, account, kp)
	limit := 100
	fmt.Println("Saving Ask.")
	setAsk(client, kp, account, limit)
	return "SetAsk succeeded", nil
}

func (s *SmartContract) GetBids(_ contractapi.TransactionContextInterface) (string, error) {
	client := newClient()
	kp := fetchPublicKey(client, account)
	number, bids := getEligibleBids(client, kp, account)
	return fmt.Sprintf("total bid number: %d, bids: %#v\n", number, bids), nil
}

func (s *SmartContract) SetBids(_ contractapi.TransactionContextInterface) (string, error) {
	client := newClient()
	kp := fetchPublicKey(client, account)

	limit := 100
	credit := 630
	setBid(client, kp, account, "alice@gmail.com", limit+10, credit+100)
	setBid(client, kp, account, "bob@gmail.com", limit, credit-100)
	setBid(client, kp, account, "evan@gmail.com", limit-10, credit)

	return "SetBids succeeded", nil
}

func main() {
	chainCode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create demo chainCode: %s", err.Error())
		return
	}

	if err := chainCode.Start(); err != nil {
		fmt.Printf("Error starting demo chainCode: %s", err.Error())
	}
}
