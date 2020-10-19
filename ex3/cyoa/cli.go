package cyoa

import (
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"
	"strconv"
	"strings"
)

func (cli *StoryCLI) runCLI() error {
	if err := keyboard.Open(); err != nil {
		return err
	}
	defer func() {
		_ = keyboard.Close()
	}()

	chapter := cli.IntroChapter

	for {
		err := cli.renderParagraph(chapter)
		if err != nil {
			return err
		}

		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		if key == keyboard.KeyEsc {
			break
		}

		if char == '0' {
			chapter = cli.IntroChapter
			continue
		}

		optionNum, err := strconv.Atoi(string(char))
		if err == nil {
			optionNum -= 1
			if options := cli.Story[chapter].Options; len(options) > optionNum {
				chapter = options[optionNum].Chapter
			}
		}
	}

	return nil
}

// renderParagraph prints paragraph's content to console
func (cli *StoryCLI) renderParagraph(chapter string) error {
	var b strings.Builder

	tm.Clear()
	tm.MoveCursor(0, 0)

	fmt.Fprint(&b, tm.Bold("Choose Your Own Adventure"))

	fmt.Fprintf(&b, "\n\nChapter: %s\n\n", cli.Story[chapter].Title)

	for _, p := range cli.Story[chapter].Paragraphs {
		fmt.Fprintln(&b, "  ", p)
	}

	fmt.Fprintln(&b, "\n", tm.Bold("Press key to select option:"))

	for i, p := range cli.Story[chapter].Options {
		opt := fmt.Sprintf("[%d] %s\n", i+1, p.Text)
		fmt.Fprintln(&b, tm.Color(opt, tm.GREEN))
	}

	fmt.Fprint(&b, "\n", tm.Color("[0] Go to intro", tm.YELLOW))
	fmt.Fprint(&b, "\n\n", tm.Color("[ESC] Exit", tm.RED))

	_, err := tm.Print(b.String())
	if err != nil {
		return err
	}

	tm.Flush()

	return nil
}
