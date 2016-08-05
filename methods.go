package gopoet

import (
	"bytes"
	"fmt"
)

type MethodSpec struct {
	// CodeBlock

	FuncSpec
	ReceiverName    string
	Receiver        ImportSpec
	IsValueReceiver bool
}

var _ CodeBlock = (*MethodSpec)(nil)

func NewMethodSpec(name, receiverName string, isValueReceiver bool, receiver ImportSpec) *MethodSpec {
	if receiver.NeedsQualifier() {
		panic(fmt.Sprintf("Cannot create method on a struct that needs to be qualified, got %#v", receiver))
	}

	return &MethodSpec{
		FuncSpec: FuncSpec{
			Name: name,
		},
		ReceiverName:    receiverName,
		IsValueReceiver: isValueReceiver,
		Receiver:        receiver,
	}
}

func (m *MethodSpec) String() string {
	writer := NewCodeWriter()

	writer.WriteStatement(m.createSignature())

	for _, st := range m.Statements {
		writer.WriteStatement(st)
	}

	writer.WriteStatement(Statement{
		BeforeIndent: -1,
		Format:       "}",
	})

	return writer.String()
}

func (m *MethodSpec) createSignature() Statement {
	formatStr := bytes.Buffer{}
	signature, args := m.Signature()

	formatStr.WriteString("func ")
	formatStr.WriteString("(")
	formatStr.WriteString(m.ReceiverName)
	formatStr.WriteString(" ")
	if !m.IsValueReceiver {
		formatStr.WriteString("*")
	}
	formatStr.WriteString(m.Receiver.GetName())
	formatStr.WriteString(") ")
	formatStr.WriteString(signature)
	formatStr.WriteString(" {")

	return Statement{
		AfterIndent: 1,
		Format:      formatStr.String(),
		Arguments:   args,
	}
}
