package ethereumservice

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumService struct {
	client *ethclient.Client
}

func New(client *ethclient.Client) EthereumService {
	return EthereumService{
		client: client,
	}
}
