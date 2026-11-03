// Copyright (C) 2024-2026 Murilo Gomes <profmugomes.com.br>
// SPDX-License-Identifier: GPL-2.0-only

// Site: https://www.profmugomes.com.br

package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	"github.com/profmugomes/mgcolumnview/v3"
	"github.com/profmugomes/mgdialogbox/v3"
	"github.com/profmugomes/mgrun/v3"
	"github.com/profmugomes/mgsmartflow/v3"

	c "profmugomes/micheckexif/controls"
)

var a fyne.App

const VERSION_APP = "2.0.0"

func main() {
	c.LoadTranslations()

	a = app.NewWithID("br.com.profmugomes.micheckexif")
	a.SetIcon(resourceMicheckexifPng)
	w := a.NewWindow("MiCheckExif")
	w.CenterOnScreen()
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(680, 697))

	sPathHome, _ := os.UserHomeDir()

	sMenu := fyne.NewMainMenu(
		fyne.NewMenu(c.T("About"),
			fyne.NewMenuItem(c.T("Support MiCheckExif"), func() {
				sURL, _ := url.Parse("https://github.com/sponsors/profmugomes")
				a.OpenURL(sURL)
			}),
			fyne.NewMenuItem(c.T("About MiCheckExif"), showAbout),
		),
	)

	flow := mgsmartflow.New()

	lstArquivos := mgcolumnview.NewColumnView([]string{"", c.T("Files")}, []float32{20, 680}, true)
	lstInfo := mgcolumnview.NewColumnView([]string{c.T("File"), c.T("Information")}, []float32{480, 170}, false)

	btnAddArquivos := widget.NewButton(c.T("Add File"), func() {
		mgdialogbox.NewOpenFile(a, c.T("Select the file(s)"), []string{".png", ".jpeg", ".jpg", ".gif", ".webp", ".pdf"}, true, func(filenames []string) {
			for _, filename := range filenames {
				lstArquivos.AddRow([]string{filename})
			}
		})
	})

	btnRemoverArquivo := widget.NewButton(c.T("Remove"), func() {
		lstArquivos.RemoveSelected()
	})
	btnRemoverTodos := widget.NewButton(c.T("Remove All"), func() {
		lstArquivos.RemoveAll()
	})

	flow.AddColumn(btnAddArquivos, btnRemoverArquivo, btnRemoverTodos)

	flow.AddRow(lstArquivos)
	flow.Resize(lstArquivos, lstArquivos.MinSize().Width, 200)

	label1 := widget.NewLabel(c.T("Remove original file"))
	label1.TextStyle = fyne.TextStyle{Bold: true}
	flow.AddRow(label1)

	cboArquivoOriginal := widget.NewSelect([]string{c.T("Yes"), c.T("No")}, nil)
	cboArquivoOriginal.SetSelectedIndex(1)
	flow.AddRow(cboArquivoOriginal)

	lblInfo := widget.NewLabel("")

	btnGetInfo := widget.NewButton(c.T("Get Data"), func() {
		lstInfo.RemoveAll()

		go func() {
			for _, filename := range lstArquivos.ListAll() {
				fyne.Do(func() {
					lblInfo.SetText(c.T("Getting Data: %s", filename.Data[0]))
				})

				c := mgrun.New(fmt.Sprintf("exiftool \"%s\"", filename.Data[0]))
				c.SetDir(sPathHome)

				c.OnStdout(func(s string) {
					for row := range strings.SplitSeq(s, "\n") {
						fyne.Do(func() {
							lstInfo.AddRow([]string{row, ""})
						})
					}
				})

				c.Run()
			}

			fyne.Do(func() {
				lblInfo.SetText("")
			})
		}()
	})

	btnRemoveData := widget.NewButton(c.T("Remove Data"), func() {
		lstInfo.RemoveAll()

		go func() {
			for _, filename := range lstArquivos.ListAll() {
				fyne.Do(func() {
					lblInfo.SetText(c.T("Removing Data: %s", filename.Data[0]))
				})

				sRun := mgrun.New(fmt.Sprintf("exiftool -all= \"%s\"", filename.Data[0]))
				sRun.SetDir(sPathHome)
				var sResult strings.Builder
				sRun.OnStdout(func(s string) {
					sResult.WriteString(s)
				})

				sRun.Run()

				if cboArquivoOriginal.SelectedIndex() == 0 {
					if _, err := os.Stat(filename.Data[0] + "_original"); err == nil {
						os.Remove(filename.Data[0] + "_original")
					}
				}

				if strings.Contains(sResult.String(), "0 image files updated") {
					fyne.Do(func() {
						lstInfo.AddRow([]string{filename.Data[0], c.T("File not changed")})
					})
				} else {
					fyne.Do(func() {
						lstInfo.AddRow([]string{filename.Data[0], c.T("Data removed")})
					})
				}
			}

			fyne.Do(func() {
				lblInfo.SetText("")
			})
		}()
	})

	flow.AddColumn(btnGetInfo, btnRemoveData)

	flow.AddRow(lblInfo)

	flow.AddRow(lstInfo)
	flow.Resize(lstInfo, lstArquivos.MinSize().Width, 250)

	w.SetMainMenu(sMenu)
	w.SetContent(flow.Container)
	w.ShowAndRun()
}
