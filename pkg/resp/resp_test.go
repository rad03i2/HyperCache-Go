package resp
import("bytes";"testing")
func TestArrayRoundTripParsing(t *testing.T){r:=NewReader(bytes.NewBufferString("*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"));v,e:=r.ReadValue();if e!=nil||len(v.Array)!=2||string(v.Array[1].Bulk)!="key"{t.Fatalf("bad parse: %#v %v",v,e)}}
func TestWriter(t *testing.T){var b bytes.Buffer;w:=NewWriter(&b);_ = w.WriteBulk([]byte("مرحبا"));if b.Len()==0{t.Fatal("empty output")}}
