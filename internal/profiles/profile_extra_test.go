package profiles

import "testing"

// ensure nil probe input is handled gracefully
func TestEvaluateResolveFree_NilProbe(t *testing.T) {
	ev := EvaluateResolveFree(nil)
	if ev.OK {
		t.Fatalf("expected not OK for nil probe")
	}
}
