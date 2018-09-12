package greeting

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"
)

type errorWriter struct {
	Err error
}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, w.Err
}

func TestGreeting_Do(t *testing.T) {
	cases := map[string]struct {
		writer    io.Writer
		clock     Clock
		msg       string
		expectErr bool
	}{
		"04時": {
			writer: new(bytes.Buffer),
			clock:  mockClock(t, "2018/08/31 04:00:00"),
			msg:    "おはよう",
		},
		"09時": {
			writer: new(bytes.Buffer),
			clock:  mockClock(t, "2018/08/31 09:00:00"),
			msg:    "おはよう",
		},
		"10時": {
			writer: new(bytes.Buffer),
			clock:  mockClock(t, "2018/08/31 10:00:00"),
			msg:    "こんにちは",
		},
		"16時": {
			writer: new(bytes.Buffer),
			clock:  mockClock(t, "2018/08/31 16:00:00"),
			msg:    "こんにちは",
		},
		"17時": {
			writer: new(bytes.Buffer),
			clock:  mockClock(t, "2018/08/31 17:00:00"),
			msg:    "こんばんは",
		},
		"03時": {
			writer: new(bytes.Buffer),
			clock:  mockClock(t, "2018/08/31 03:00:00"),
			msg:    "こんばんは",
		},
		"エラー": {
			writer:    &errorWriter{Err: errors.New("error")},
			expectErr: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			g := Greeting{
				Clock: tc.clock,
			}

			switch err := g.Do(tc.writer); true {
			case err == nil && tc.expectErr:
				t.Error("expected error did not occur")
			case err != nil && !tc.expectErr:
				t.Error("unexpected error", err)
			}

			if buff, ok := tc.writer.(*bytes.Buffer); ok {
				msg := buff.String()
				if msg != tc.msg {
					t.Errorf("greeting msg wnot %s but got %s", tc.msg, msg)
				}
			}
		})
	}
}

func mockClock(t *testing.T, v string) Clock {
	t.Helper()
	now, err := time.Parse("2006/01/02 15:04:05", v)
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	return ClockFunc(func() time.Time {
		return now
	})
}
