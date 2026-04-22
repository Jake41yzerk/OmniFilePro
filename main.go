package main

import (
	"image/color"
	"os"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Warna dari icon lu
var (
	kuningOmni = color.NRGBA{R: 255, G: 215, B: 0, A: 255} // Kuning emas
	unguTua = color.NRGBA{R: 138, G: 43, B: 226, A: 255} // Ungu tua
	unguMuda = color.NRGBA{R: 186, G: 85, B: 211, A: 255} // Ungu muda
	putih = color.White
	hitam = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
)

func main() {
	a := app.New()
	w := a.NewWindow("OmniFilePro")
	w.Resize(fyne.NewSize(400, 700))

	// 1. SPLASH SCREEN
	splashBg := canvas.NewRectangle(kuningOmni)
	splashText := canvas.NewText("OmniFilePro", hitam)
	splashText.TextSize = 24
	splashText.TextStyle = fyne.TextStyle{Bold: true}
	splashText.Alignment = fyne.TextAlignCenter
	splash := container.NewStack(splashBg, container.NewCenter(splashText))
	w.SetContent(splash)

	// 2 detik kemudian minta izin + masuk app
	go func() {
		time.Sleep(2 * time.Second)
		dialog.ShowConfirm("Izin Diperlukan", "OmniFilePro butuh akses penyimpanan untuk membaca file. Izinkan?",
			func(ok bool) {
				if!ok {
					a.Quit()
				} else {
					w.SetContent(makeMainUI(w))
				}
			}, w)
	}()

	w.ShowAndRun()
}

func makeMainUI(w fyne.Window) fyne.CanvasObject {
	// ===== LATAR PUTIH GLOBAL =====
	bgPutih := canvas.NewRectangle(putih)

	// ===== HEADER KUNING =====
	headerKuning := canvas.NewRectangle(kuningOmni)
	headerKuning.Resize(fyne.NewSize(400, 50))
	headerText := canvas.NewText("OmniFilePro", hitam)
	headerText.TextStyle = fyne.TextStyle{Bold: true}
	header := container.NewStack(headerKuning, container.NewCenter(headerText))

	// ===== TAB EXPLORER =====
	var fileData []string
	files, err := os.ReadDir("/storage/emulated/0")
	if err!= nil {
		fileData = []string{"Gagal baca storage. Cek izin di Pengaturan HP."}
	} else {
		for _, f := range files {
			if f.IsDir() {
				fileData = append(fileData, "📁 "+f.Name())
			} else {
				fileData = append(fileData, "📄 "+f.Name())
			}
		}
	}

	fileList := widget.NewList(
		func() int { return len(fileData) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fileData[i])
		})

	explorerTab := container.NewBorder(
		widget.NewLabel("Internal Storage"), nil, nil, nil, fileList,
	)

	// ===== TAB KOMPRES + TOGGLE =====
	kompresAktif := widget.NewCheck("Aktifkan Fitur Kompres", func(on bool) {})
	kompresAktif.SetChecked(true)
	hapusAsli := widget.NewCheck("Hapus file asli setelah kompres", func(on bool) {})

	kompresTab := container.NewVBox(
		widget.NewLabelWithStyle("Pengaturan Kompres", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		kompresAktif,
		hapusAsli,
		widget.NewButton("Pilih File untuk Dikompres", func() {
			dialog.ShowInformation("Info", "Fitur pilih file belum jadi", w)
		}),
	)

	// ===== TAB DUPLIKAT =====
	duplikatTab := container.NewVBox(
		widget.NewLabelWithStyle("Cari File Duplikat", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewButtonWithIcon("Mulai Scan Duplikat", theme.SearchIcon(), func() {
			dialog.ShowInformation("Info", "Fitur scan belum jadi", w)
		}),
	)

	// ===== SIDEBAR GRADASI UNGU =====
	gradasiUngu := canvas.NewLinearGradient(unguTua, unguMuda, 0) // 0 = vertikal
	sidebarContent := container.NewVBox(
		widget.NewLabelWithStyle("MENU", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewAccordion(
			widget.NewAccordionItem("Penyimpanan",
				container.NewVBox(
					widget.NewButtonWithIcon("Laman Rumah", theme.HomeIcon(), func() {}),
					widget.NewButtonWithIcon("Internal", theme.FolderIcon(), func() {}),
				),
			),
			widget.NewAccordionItem("Tools",
				container.NewVBox(
					widget.NewButtonWithIcon("Analisis", theme.ComputerIcon(), func() {}),
					widget.NewButtonWithIcon("Bersihkan", theme.DeleteIcon(), func() {}),
				),
			),
		),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Pengaturan", theme.SettingsIcon(), func() {}),
	)
	sidebar := container.NewStack(gradasiUngu, container.NewPadded(sidebarContent))

	// ===== GABUNGIN SEMUA =====
	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Explorer", theme.FolderOpenIcon(), explorerTab),
		container.NewTabItemWithIcon("Kompres", theme.FileImageIcon(), kompresTab),
		container.NewTabItemWithIcon("Duplikat", theme.ViewRefreshIcon(), duplikatTab),
	)

	split := container.NewHSplit(sidebar, mainTabs)
	split.SetOffset(0.35) // Sidebar 35%

	// Stack: Layer 1 putih, Layer 2 header + konten
	ui := container.NewStack(
		bgPutih,
		container.NewBorder(header, nil, nil, nil, split),
	)

	return ui
}
