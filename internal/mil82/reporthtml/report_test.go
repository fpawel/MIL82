package reporthtml

import (
	"encoding/json"
	"fmt"
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
