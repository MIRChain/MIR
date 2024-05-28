package qbfttypes

import (
	"fmt"
	"io"
	"math/big"

	"github.com/MIRChain/MIR/consensus/istanbul"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/rlp"
)

type Preprepare[P crypto.PublicKey] struct {
	CommonPayload
	Proposal                  istanbul.Proposal
	JustificationRoundChanges []*SignedRoundChangePayload
	JustificationPrepares     []*Prepare
}

func NewPreprepare[P crypto.PublicKey](sequence *big.Int, round *big.Int, proposal istanbul.Proposal) *Preprepare[P] {
	return &Preprepare[P]{
		CommonPayload: CommonPayload{
			code:     PreprepareCode,
			Sequence: sequence,
			Round:    round,
		},
		Proposal: proposal,
	}
}

func (m *Preprepare[P]) EncodePayloadForSigning() ([]byte, error) {
	return rlp.EncodeToBytes(
		[]interface{}{
			m.Code(),
			[]interface{}{m.Sequence, m.Round, m.Proposal},
		})
}

func (m *Preprepare[P]) EncodeRLP(w io.Writer) error {
	return rlp.Encode(
		w,
		[]interface{}{
			[]interface{}{
				[]interface{}{m.Sequence, m.Round, m.Proposal},
				m.signature,
			},
			[]interface{}{
				m.JustificationRoundChanges,
				m.JustificationPrepares,
			},
		})
}

func (m *Preprepare[P]) DecodeRLP(stream *rlp.Stream) error {
	var message struct {
		SignedPayload struct {
			Payload struct {
				Sequence *big.Int
				Round    *big.Int
				Proposal *types.Block[P]
			}
			Signature []byte
		}
		Justification struct {
			RoundChanges []*SignedRoundChangePayload
			Prepares     []*Prepare
		}
	}
	if err := stream.Decode(&message); err != nil {
		return err
	}
	m.code = PreprepareCode
	m.Sequence = message.SignedPayload.Payload.Sequence
	m.Round = message.SignedPayload.Payload.Round
	m.Proposal = message.SignedPayload.Payload.Proposal
	m.signature = message.SignedPayload.Signature
	m.JustificationPrepares = message.Justification.Prepares
	m.JustificationRoundChanges = message.Justification.RoundChanges
	return nil
}

func (m *Preprepare[P]) String() string {
	return fmt.Sprintf("code: %d, sequence: %d, round: %d, proposal: %v", m.code, m.Sequence, m.Round, m.Proposal.Hash().Hex())
}
