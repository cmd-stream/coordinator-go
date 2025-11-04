package main

import (
	"os"
	"reflect"

	man "github.com/cmd-stream/coordinator-go/manager"
	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/workflw-go/manager"),
	)
	assert.EqualError(err, nil)

	ct := reflect.TypeFor[man.CmdEnvelope]()
	err = g.AddStruct(ct)
	assert.EqualError(err, nil)

	err = g.AddDTS(ct)
	assert.EqualError(err, nil)

	ot := reflect.TypeFor[man.OutcomeEnvelope]()
	err = g.AddStruct(ot)
	assert.EqualError(err, nil)

	err = g.AddDTS(ot)
	assert.EqualError(err, nil)

	err = g.AddDefinedType(reflect.TypeFor[man.ServiceID]())
	assert.EqualError(err, nil)

	err = g.AddStruct(reflect.TypeFor[man.WorkflowID]())
	assert.EqualError(err, nil)

	err = g.AddDefinedType(reflect.TypeFor[man.PartID]())
	assert.EqualError(err, nil)

	err = g.AddDefinedType(reflect.TypeFor[man.PartSeq]())
	assert.EqualError(err, nil)

	bs, err := g.Generate()
	assert.EqualError(err, nil)

	err = os.WriteFile("./mus-format.gen.go", bs, 0644)
	assert.EqualError(err, nil)
}
