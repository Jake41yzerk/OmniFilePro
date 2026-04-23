package main

import (
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	currentPath = "/storage/emulated/0"
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
		cekIzinDanMasuk(w)
	}()

	w.ShowAndRun()
}

func cekIzinDanMasuk(w fyne.Window) {
	_, err := os.ReadDir(currentPath)
	if err!= nil && os.IsPermission(err) {
		dialog.ShowConfirm("Izin Belum Lengkap",
			"Android 11+ butuh 'Akses semua file'.\n\nPencet OK buat buka Settings, terus nyalain izin buat OmniFilePro.",
			func(ok bool) {
				if ok && runtime.GOOS == "android" {
					exec.Command("am", "start", "-a", "android.settings.MANAGE_ALL_FILES_ACCESS_PERMISSION").Run()
				}
				w.SetContent(makeMainUI(w))
			}, w)
	} else {
		w.SetContent(makeMainUI(w))
	}
}

func makeMainUI(w fyne.Window) fyne.CanvasObject {
	bgPutih := canvas.NewRectangle(putih)

	// ===== HEADER + TOMBOL REFRESH =====
	headerKuning := canvas.NewRectangle(kuningOmni)
	headerText := canvas.NewText("OmniFilePro", hitam)
	headerText.TextStyle = fyne.TextStyle{Bold: true}
	
	refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		w.SetContent(makeMainUI(w))
	})

	sidebarOpen := true
	gradasiUngu := canvas.NewVerticalGradient(unguTua, unguMuda)
	gradasiUngu.Resize(fyne.NewSize(200, 700))

	sidebarContent := container.NewVBox(
		widget.NewButtonWithIcon("Internal", theme.FolderIcon(), func() {
			currentPath = "/storage/emulated/0"
			w.SetContent(makeMainUI(w))
		}),
		widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() {
			currentPath = "/storage/emulated/0/Download"
			w.SetContent(makeMainUI(w))
		}),
		widget.NewButtonWithIcon("DCIM", theme.FileImageIcon(), func() {
			currentPath = "/storage/emulated/0/DCIM"
			w.SetContent(makeMainUI(w))
		}),
		widget.NewSeparator(),
		widget.NewButtonWithIcon("Buka Settings Izin", theme.SettingsIcon(), func() {
			if runtime.GOOS == "android" {
				exec.Command("am", "start", "-a", "android.settings.MANAGE_ALL_FILES_ACCESS_PERMISSION").Run()
			}
		}),
	)
	sidebar := container.NewStack(gradasiUngu, container.NewPadded(sidebarContent))
	sidebar.Resize(fyne.NewSize(200, 700))

	// ===== EXPLORER =====
	var fileData []string
	var fileIsDir []bool
	
	files, err := os.ReadDir(currentPath)
	if err!= nil {
		fileData = []string{
			"ERROR: " + err.Error(),
			"",
			"Cara benerin:",
			"1. Pencet tombol Settings di sidebar",
			"2. Nyalain 'Akses semua file'",
			"3. Pencet tombol refresh ↻",
		}
		fileIsDir = []bool{false, false, false}
	} else {
		if currentPath!= "/storage/emulated/0" {
			fileData = append(fileData, "📁..")
			fileIsDir = append(fileIsDir, true)
		}
		if len(files) == 0 {
			fileData = append(fileData, "(Folder ini kosong)")
			fileIsDir = append(fileIsDir, false)
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

	pathLabel := canvas.NewText("Path: "+currentPath, hitam)
	pathLabel.TextSize = 11

	fileList := widget.NewList(
		func() int { return len(fileData) },
		func() fyne.CanvasObject { return canvas.NewText("template", hitam) },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = fileData[i]
			o.Refresh()
		})
	
	fileList.OnSelected = func(id widget.ListItemID) {
		fileList.Unselect(id)
		if id < len(fileIsDir) && fileIsDir[id] {
			selectedName := fileData[id]
			if len(selectedName) > 3 { selectedName = selectedName[3:] }
			
			if selectedName == ".." {
				currentPath = filepath.Dir(currentPath)
			} else {
				currentPath = filepath.Join(currentPath, selectedName)
			}
			w.SetContent(makeMainUI(w))
		}
	}

	judulExplorer := canvas.NewText("Explorer", hitam)
	judulExplorer.TextStyle = fyne.TextStyle{Bold: true}
	explorerTab := container.NewBorder(
		container.NewVBox(judulExplorer, pathLabel), nil, nil, nil, fileList,
	)

	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Explorer", theme.FolderOpenIcon(), explorerTab),
	)

	kontenArea := container.NewStack(mainTabs)
	
	drawerBtn := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		if sidebarOpen { sidebar.Hide(); sidebarOpen = false } else { sidebar.Show(); sidebarOpen = true }
	})

	header := container.NewStack(
		headerKuning,
		container.NewBorder(nil, nil, container.NewHBox(container.NewPadded(drawerBtn)), container.NewPadded(refreshBtn), container.NewCenter(headerText)),
	)

	ui := container.NewStack(
		bgPutih,
		container.NewBorder(header, nil, sidebar, nil, kontenArea),
	)

	return ui
}
