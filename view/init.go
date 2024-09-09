package view

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/mnbjhu/plog/input"
)

var conf input.Config

func Init() {
	another := true

	configForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Input").
				Options(
					huh.NewOption("stdout", "stdout"),
					huh.NewOption("stderr", "stderr"),
				).Value(&conf.Input),
			huh.NewConfirm().Title("Add Columns Now").Value(&another),
		).Title("Configuration"),
	)

	err := configForm.Run()
	if err != nil {
		panic(err)
	}

	for another {
		title := ""
		expr := ""
		width := "8"
		newColForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Title").Value(&title),
				huh.NewSelect[string]().
					Title("Pattern").
					Options(huh.NewOption("Word", "(?<%s>\\w+)"),
						huh.NewOption("Integer", "(?<%s>\\d+)"),
						huh.NewOption("Float", "(?<%s>\\d+\\.\\d+)"),
						huh.NewOption("Any", "(?<%s>\\w+)"),
						huh.NewOption("Bracketed", "\\[(?<%s>[^\\[]+)\\]"),
						huh.NewOption("Rest", "rest"),
					).Value(&expr),
				huh.NewInput().Title("Width").
					Validate(func(b string) error {
						_, err := strconv.Atoi(b)
						return err
					}).
					Value(&width),
			).Title(fmt.Sprintf("Column %d", len(conf.Columns)+1)),
		)
		err := newColForm.Run()
		if err != nil {
			panic(err)
		}
		if expr == "rest" {
			expr = "(?<%s>.*)"
			another = false
		}
		w, _ := strconv.Atoi(width)
		conf.Columns = append(conf.Columns, input.ColumnDef{
			Title: title,
			Width: w,
		})
		if conf.Regex != "" {
			conf.Regex = conf.Regex + "\\s+" + fmt.Sprintf(expr, title)
		} else {
			conf.Regex = fmt.Sprintf(expr, title)
		}
	}
	conf.Regex = fmt.Sprintf("^%s$", conf.Regex)
	conf.Save()
}
