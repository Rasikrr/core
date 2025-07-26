package telegram

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Util(t *testing.T) {
	buttons := []string{"Button1", "Button2", "Button3", "Button4", "Button5"}
	keyboard := CreateKeyboard(buttons, 2)
	require.Equal(t, 3, len(keyboard))
	for _, row := range keyboard {
		require.True(t, len(row) <= 2)
		t.Logf("%v", row)
	}
}
