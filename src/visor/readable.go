package visor

import (
	"github.com/skycoin/skycoin/src/coin"
)

// Transaction wraps around coin.Transaction, tagged with its status.  This allows us
// to include unconfirmed txns
type Transaction struct {
	Transaction coin.Transaction
	Status      TransactionStatus
	Time        uint64
}

// TransactionStatus represents the transaction status
type TransactionStatus struct {
	Confirmed bool
	// If confirmed, how many blocks deep in the chain it is. Will be at least 1 if confirmed.
	Height uint64
	// If confirmed, the sequence of the block in which the transaction was executed
	BlockSeq uint64
}

// NewUnconfirmedTransactionStatus creates unconfirmed transaction status
func NewUnconfirmedTransactionStatus() TransactionStatus {
	return TransactionStatus{
		Confirmed: false,
		Height:    0,
		BlockSeq:  0,
	}
}

// NewConfirmedTransactionStatus creates confirmed transaction status
func NewConfirmedTransactionStatus(height uint64, blockSeq uint64) TransactionStatus {
	if height == 0 {
		logger.Panic("Invalid confirmed transaction height")
	}
	return TransactionStatus{
		Confirmed: true,
		Height:    height,
		BlockSeq:  blockSeq,
	}
}

// TransactionInput includes the UxOut spent in a transaction and the calculated hours of the output at spending time
type TransactionInput struct {
	UxOut           coin.UxOut
	CalculatedHours uint64
}

// NewTransactionInput creates a TransactionInput.
// calculateHoursTime is the time against which the CalculatedHours should be computed
func NewTransactionInput(ux coin.UxOut, calculateHoursTime uint64) (TransactionInput, error) {
	// The overflow bug causes this to fail for some transactions, allow it to pass
	calculatedHours, err := ux.CoinHours(calculateHoursTime)
	if err != nil {
		logger.Critical().Warningf("Ignoring NewTransactionInput ux.CoinHours failed: %v", err)
		calculatedHours = 0
	}

	return TransactionInput{
		UxOut:           ux,
		CalculatedHours: calculatedHours,
	}, nil
}

// BlockchainMetadata encapsulates useful information from the coin.Blockchain
type BlockchainMetadata struct {
	// Most recent block
	HeadBlock *coin.SignedBlock
	// Number of unspent outputs in the coin.Blockchain
	Unspents uint64
	// Number of known unconfirmed txns
	Unconfirmed uint64
}

// NewBlockchainMetadata creates blockchain meta data
func NewBlockchainMetadata(head *coin.SignedBlock, unconfirmedLen, unspentsLen uint64) (*BlockchainMetadata, error) {
	return &BlockchainMetadata{
		HeadBlock:   head,
		Unspents:    unspentsLen,
		Unconfirmed: unconfirmedLen,
	}, nil
}
