package support

import (
	"bytes"
	"encoding/json"

	man "github.com/cmd-stream/coordinator-go/manager"
	muss "github.com/mus-format/mus-stream-go"
)

func MakeOutcomeSerializerMUS(ser muss.Serializer[man.Outcome]) OutcomeSerializerMUS {
	return OutcomeSerializerMUS{ser}
}

type OutcomeSerializerMUS struct {
	ser muss.Serializer[man.Outcome]
}

func (s OutcomeSerializerMUS) Marshal(o man.Outcome) (bs []byte, err error) {
	bs = make([]byte, 0, s.ser.Size(o))
	buf := bytes.NewBuffer(bs)
	s.ser.Marshal(o, buf)
	bs = buf.Bytes()
	return
}

func (s OutcomeSerializerMUS) Unmarshal(bs []byte) (o man.Outcome, err error) {
	buf := bytes.NewBuffer(bs)
	o, _, err = s.ser.Unmarshal(buf)
	return
}

type OutcomeSerializerJSON struct{}

func (s OutcomeSerializerJSON) Marshal(o man.Outcome) (bs []byte, err error) {
	return json.Marshal(o)
}

func (s OutcomeSerializerJSON) Unmarshal(bs []byte) (o man.Outcome, err error) {
	err = json.Unmarshal(bs, &o)
	return
}
