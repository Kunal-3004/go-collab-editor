package domain

import (
	"errors"
	"sort"
)

type OpType string

const (
	OpInsert OpType = "INSERT"
	OpDelete OpType = "DELETE"
)

type Operation struct {
	Type     OpType  `json:"type"`
	Position float64 `json:"position"`
	Value    string  `json:"value"`
	ClientID string  `json:"clientId"`
}

type Document struct {
	ID         string
	Operations []Operation
}

func NewDocument(id string) *Document {
	return &Document{
		ID:         id,
		Operations: []Operation{},
	}
}

func (d *Document) ApplyOperation(op Operation) error {
	if op.Value == "" && op.Type == OpInsert {
		return errors.New("cannot insert empty value")
	}

	d.Operations = append(d.Operations, op)

	sort.Slice(d.Operations, func(i, j int) bool {
		return d.Operations[i].Position < d.Operations[j].Position
	})

	return nil
}
