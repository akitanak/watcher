package app

import (
	"testing"
)

func TestChannelWriter_Write(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		in      []byte
		want    []byte
		wantLen int
	}{
		"if input is empty, return empty": {
			in:      []byte{},
			want:    []byte{},
			wantLen: 0,
		},
		"if input is not empty, return input": {
			in:      []byte("hello"),
			want:    []byte("hello"),
			wantLen: 5,
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sendTo := make(chan []byte)
			w := NewChannelWriter(sendTo)
			writerFn := func() {
				gotLen, err := w.Write(tc.in)
				if err != nil {
					t.Fatalf("expect no error, but got %+v", err)
				}
				if gotLen != tc.wantLen {
					t.Errorf("want %d, but got %d", tc.wantLen, gotLen)
				}
			}

			go writerFn()

			got := <-sendTo
			if string(got) != string(tc.want) {
				t.Errorf("want %s, but got %s", tc.want, got)
			}
		})
	}
}
