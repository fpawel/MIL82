package mil82

import (
	"encoding/json"
	"fmt"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/mil82/internal/data"
	"testing"
)

func BenchmarkReportParty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = reportParty(1)
	}
}

func TestReportParty(t *testing.T) {
	reportParty(1)
	b, _ := json.MarshalIndent(reportParty(1), "", "  ")
	fmt.Printf("%s", string(b))
}

func TestReportPartyHtml(t *testing.T) {
	b, _ := json.MarshalIndent(reportParty(1), "", "  ")
	fmt.Printf("%s", string(b))
}

func init() {
	data.DB = gohelp.OpenSqliteDBx(`C:\GOPATH\src\github.com\fpawel\mil82\build\mil82.sqlite`)
}
