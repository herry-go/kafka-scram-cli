package scram

import (
	"crypto/sha256"
	"crypto/sha512"
	"github.com/xdg-go/scram"
	"hash"
)

// SHA256 hash function
func SHA256() hash.Hash { return sha256.New() }

// SHA512 hash function
func SHA512() hash.Hash { return sha512.New() }

// XDGSCRAMClient implements sarama.SCRAMClient
type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	return x.ClientConversation.Step(challenge)
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
} 