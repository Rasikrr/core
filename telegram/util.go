package telegram

type KeyBoard [][]string
type KeyBoardRow []string

func CreateKeyboard(buttons []string, cols int) KeyBoard {
	keyBoard := make(KeyBoard, 0, len(buttons)/cols+1)
	for i := 0; i < len(buttons); i += cols {
		end := i + cols
		if end > len(buttons) {
			end = len(buttons)
		}
		row := make(KeyBoardRow, 0, cols)
		for _, button := range buttons[i:end] {
			row = append(row, button)
		}
		if len(row) > 0 {
			keyBoard = append(keyBoard, row)
		}
	}
	return keyBoard
}
