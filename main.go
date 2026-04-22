package main

import (
	"image/color"
	"os"
	"path/filepath"
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

var (
	kuningOmni = color.NRGBA{R: 255, G: 215, B: 0, A: 255}
	unguTua = color.NRGBA{R: 138, G: 43, B: 226, A: 255}
	unguMuda = color.NRGBA{R: 186, G: 85, B: 211, A: 255}
	putih = color.White
	hitam = color.NRGBA{R: 20, G: 20, B: 20, A: 255}
	currentPath = "/storage/emulated/0" // Path aktif sekarang
)

func main() {
	a := app.New()
	w := a.NewWindow("OmniFilePro")
	w.Resize(fyne.NewSize(400, 700))

	splashBg := canvas.NewRectangle(kuningOmni)
	splashText := canvas.NewText("OmniFilePro", hitam)
	splashText.TextSize = 24
	splashText.TextStyle = fyne.TextStyle{Bold: true}
	w.SetContent(container.NewStack(splashBg, container.NewCenter(splashText)))

	go func() {
		time.Sleep(2 * time.Second)
		dialog.ShowConfirm("Izin Diperlukan", "OmniFilePro butuh akses penyimpanan. Izinkan?",
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
	bgPutih := canvas.NewRectangle(putih)

	// ===== HEADER KUNING + TOMBOL DRAWER =====
	headerKuning := canvas.NewRectangle(kuningOmni)
	headerText := canvas.NewText("OmniFilePro", hitam)
	headerText.TextStyle = fyne.TextStyle{Bold: true}
	
	sidebarOpen := true
	
	// ===== SIDEBAR DRAWER GRADASI UNGU =====
	gradasiUngu := canvas.NewVerticalGradient(unguTua, unguMuda)
	gradasiUngu.Resize(fyne.NewSize(200, 700))

	buatLabelItem := func(teks string) *fyne.Container {
		txt := canvas.NewText(teks, putih)
		txt.TextSize = 14
		return container.NewHBox(layout.NewSpacer(), txt, layout.NewSpacer())
	}

	sidebarContent := container.NewVBox(
		buatLabelItem("MENU"),
		widget.NewSeparator(),
		widget.NewButtonWithIcon("Internal", theme.FolderIcon(), func() {
			currentPath = "/storage/emulated/0"
			w.SetContent(makeMainUI(w)) // Refresh ke root
		}),
		widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() {
			currentPath = "/storage/emulated/0/Download"
			w.SetContent(makeMainUI(w))
		}),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Pengaturan", theme.SettingsIcon(), func() {}),
	)
	
	sidebar := container.NewStack(gradasiUngu, container.NewPadded(sidebarContent))
	sidebar.Resize(fyne.NewSize(200, 700))

	// ===== EXPLORER YANG BISA BUKA FOLDER =====
	var fileData []string
	var fileIsDir []bool // Nandain mana folder mana file
	
	files, err := os.ReadDir(currentPath)
	if err!= nil {
		fileData = []string{"Error: " + err.Error()}
		fileIsDir = []bool{false}
	} else {
		// Tambahin tombol ".." buat back kecuali di root
		if currentPath!= "/storage/emulated/0" {
			fileData = append(fileData, "📁..")
			fileIsDir = append(fileIsDir, true)
		}
		for _, f := range files {
			if f.IsDir() {
				fileData = append(fileData, "📁 "+f.Name())
				fileIsDir = append(fileIsDir, true)
			} else {
				fileData = append(fileData, "📄 "+f.Name())
				fileIsDir = append(fileIsDir, false)
			}
		}
	}

	pathLabel := canvas.NewText(currentPath, hitam)
	pathLabel.TextSize = 12

	fileList := widget.NewList(
		func() int { return len(fileData) },
		func() fyne.CanvasObject { return canvas.NewText("template", hitam) },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = fileData[i]
			o.Refresh()
		})
	
	// INI LOGIKA BUKA FOLDERNYA BOS
	fileList.OnSelected = func(id widget.ListItemID) {
		fileList.Unselect(id) // Biar gak nge-select biru
		
		if fileIsDir[id] {
			selectedName := fileData[id][3:] // Buang emoji "📁 "
			if selectedName == ".." {
				// Kalo pencet ".." = balik ke folder parent
				currentPath = filepath.Dir(currentPath)
			} else {
				// Kalo folder biasa = masuk ke dalem
				currentPath = filepath.Join(currentPath, selectedName)
			}
			w.SetContent(makeMainUI(w)) // Refresh UI pake path baru
		} else {
			dialog.ShowInformation("File", "Buka file "+fileData[id]+" belum jadi", w)
		}
	}

	judulExplorer := canvas.NewText("Explorer", hitam)
	judulExplorer.TextStyle = fyne.TextStyle{Bold: true}
	explorerTab := container.NewBorder(
		container.NewVBox(judulExplorer, pathLabel), nil, nil, nil, fileList,
	)

	// ===== TAB LAINNYA =====
	judulKompres := canvas.NewText("Pengaturan Kompres", hitam)
	judulKompres.TextStyle = fyne.TextStyle{Bold: true}
	kompresTab := container.NewVBox(judulKompres, widget.NewCheck("Aktifkan Fitur Kompres", func(on bool) {}))

	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Explorer", theme.FolderOpenIcon(), explorerTab),
		container.NewTabItemWithIcon("Kompres", theme.FileImageIcon(), kompresTab),
	)

	// ===== DRAWER BUTTON =====
	kontenArea := container.NewStack(mainTabs)
	
	drawerBtn := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		if sidebarOpen {
			sidebar.Hide()
			sidebarOpen = false
		} else {
			sidebar.Show()
			sidebarOpen = true
		}
	})

	header := container.NewStack(
		headerKuning,
		container.NewBorder(nil, nil, container.NewPadded(drawerBtn), nil, container.NewCenter(headerText)),
	)

	ui := container.NewStack(
		bgPutih,
		container.NewBorder(header, nil, sidebar, nil, kontenArea),
	)

	return ui
}
