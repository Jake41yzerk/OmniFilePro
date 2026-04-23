package main

import (
    "fmt"
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
    kuningOmni  = color.NRGBA{R: 255, G: 215, B: 0, A: 255}
    kuningTua   = color.NRGBA{R: 255, G: 180, B: 0, A: 255}
    unguTua     = color.NRGBA{R: 60, G: 0, B: 90, A: 200} // lebih gelap & transparan
    unguMuda    = color.NRGBA{R: 138, G: 43, B: 226, A: 200}
    putih       = color.White
    hitam       = color.NRGBA{R: 20, G: 20, B: 20, A: 255}
    currentPath = "/storage/emulated/0"
)

func main() {
    a := app.New()
    w := a.NewWindow("OmniFilePro")
    w.Resize(fyne.NewSize(400, 700))

    // Splash Screen
    splashBg := canvas.NewRectangle(kuningOmni)
    splashText := canvas.NewText("OmniFilePro", hitam)
    splashText.TextSize = 24
    splashText.TextStyle = fyne.TextStyle{Bold: true}
    
    w.SetContent(container.NewStack(
        splashBg, 
        container.NewCenter(splashText),
    ))

    go func() {
        time.Sleep(1 * time.Second)
        cekIzinDanMasuk(w)
    }()

    w.ShowAndRun()
}

// Cek izin storage, kalau deny -> lempar ke Settings
func cekIzinDanMasuk(w fyne.Window) {
    _, err := os.ReadDir(currentPath)
    if err != nil && os.IsPermission(err) {
        dialog.ShowConfirm("Izin Belum Lengkap",
            "Android 11+ butuh 'Akses semua file'.\n\nPencet OK untuk buka Settings, lalu nyalakan izin untuk OmniFilePro.",
            func(ok bool) {
                if ok && runtime.GOOS == "android" {
                    _ = exec.Command("am", "start",
                        "-a", "android.settings.MANAGE_ALL_FILES_ACCESS_PERMISSION").Run()
                }
                w.SetContent(makeMainUI(w))
            }, w)
        return
    }
    w.SetContent(makeMainUI(w))
}

func makeMainUI(w fyne.Window) fyne.CanvasObject {
    bgPutih := canvas.NewRectangle(putih)

    // ===== HEADER GRADIENT KUNING =====
    headerGradient := canvas.NewLinearGradient(kuningOmni, kuningTua, 90)
    headerGradient.SetMinSize(fyne.NewSize(400, 80))

    title := canvas.NewText("OmniFilePro", hitam)
    title.TextStyle = fyne.TextStyle{Bold: true}
    title.TextSize = 20

    // Tombol refresh izin
    refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
        cekIzinDanMasuk(w)
    })

    // Variabel kontrol sidebar
    var sidebar *fyne.Container
    var overlay *fyne.Container
    sidebarOpen := false

    // ===== EXPLORER DATA =====
    fileData, fileIsDir := loadDirectoryData()

    pathLabel := canvas.NewText("Path: "+currentPath, hitam)
    pathLabel.TextSize = 11

    fileList := widget.NewList(
        func() int { return len(fileData) },
        func() fyne.CanvasObject {
            return canvas.NewText("item", hitam)
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            t := o.(*canvas.Text)
            t.Text = fileData[i]
            t.Refresh()
        },
    )

    fileList.OnSelected = func(id widget.ListItemID) {
        fileList.Unselect(id)
        if id >= len(fileIsDir) {
            return
        }
        if !fileIsDir[id] {
            return
        }
        
        selectedName := fileData[id]
        if len(selectedName) > 3 && (selectedName[:3] == "📁 " || selectedName[:3] == "📄 ") {
            selectedName = selectedName[3:]
        }
        
        if selectedName == ".." {
            currentPath = filepath.Dir(currentPath)
        } else {
            currentPath = filepath.Join(currentPath, selectedName)
        }
        w.SetContent(makeMainUI(w))
    }

    judulExplorer := canvas.NewText("Explorer", hitam)
    judulExplorer.TextStyle = fyne.TextStyle{Bold: true}

    explorerTab := container.NewBorder(
        container.NewVBox(judulExplorer, pathLabel),
        nil, nil, nil,
        fileList,
    )

    // ===== MENU HOME =====
    menuGrid := makeHomeIconGrid(w)

    mainTabs := container.NewAppTabs(
        container.NewTabItemWithIcon("Home", theme.HomeIcon(), menuGrid),
        container.NewTabItemWithIcon("Explorer", theme.FolderOpenIcon(), explorerTab),
    )

    contentArea := container.NewStack(mainTabs)

    // ===== SIDEBAR =====
    sideGradient := canvas.NewLinearGradient(unguTua, unguMuda, 0)
    sideGradient.SetMinSize(fyne.NewSize(220, 700))

    sideButtons := container.NewVBox(
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
        widget.NewButtonWithIcon("Music", theme.VolumeUpIcon(), func() {
            currentPath = "/storage/emulated/0/Music"
            w.SetContent(makeMainUI(w))
        }),
        widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
            dialog.ShowInformation("Info", "Settings akan segera hadir.", w)
        }),
    )

    sidebarInner := container.NewBorder(
        layout.NewSpacer(), nil, nil, nil,
        container.NewPadded(sideButtons),
    )

    sidebar = container.NewStack(sideGradient, sidebarInner)
    sidebar.Resize(fyne.NewSize(220, 700))
    sidebar.Hide()

    // ===== OVERLAY =====
    overlayBg := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0}) // Start transparent
    overlayBg.Hide()

    overlay = container.NewStack(
        overlayBg,
        container.NewHBox(sidebar, layout.NewSpacer()),
    )
    overlay.Hide()

    // ===== TOMBOL MENU =====
    drawerBtn := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
        sidebarOpen = !sidebarOpen
        animateSidebar(sidebarOpen, sidebar, overlayBg, overlay)
    })

    headerContent := container.NewBorder(
        nil, nil,
        container.NewHBox(container.NewPadded(drawerBtn)),
        container.NewHBox(container.NewPadded(refreshBtn)),
        container.NewCenter(title),
    )

    header := container.NewStack(headerGradient, headerContent)

    // ROOT LAYOUT
    root := container.NewStack(
        bgPutih,
        container.NewBorder(header, nil, nil, nil, contentArea),
        overlay,
    )

    return root
}

func animateSidebar(open bool, sidebar *fyne.Container, overlayBg *canvas.Rectangle, overlay *fyne.Container) {
    if open {
        overlay.Show()
        overlayBg.Show()
        sidebar.Show()
    }

    anim := fyne.NewAnimation(time.Millisecond*200, func(f float32) {
        if !open {
            f = 1 - f
        }
        overlayBg.FillColor = color.NRGBA{0, 0, 0, uint8(100 * f)}
        overlayBg.Refresh()
    })
    
    anim.Curve = fyne.AnimationEaseInOut
    anim.SetOnCompleted(func() {
        if !open {
            overlay.Hide()
            overlayBg.Hide()
            sidebar.Hide()
        }
    })
    anim.Start()
}

func loadDirectoryData() ([]string, []bool) {
    var fileData []string
    var fileIsDir []bool

    files, err := os.ReadDir(currentPath)
    if err != nil {
        fileData = []string{
            "ERROR: " + err.Error(),
            "",
            "Cara benerin:",
            "1. Pencet tombol refresh ↻ di kanan atas",
            "2. Nyalain 'Akses semua file'",
            "3. Balik ke app",
        }
        fileIsDir = make([]bool, len(fileData))
        return fileData, fileIsDir
    }

    if currentPath != "/storage/emulated/0" {
        fileData = append(fileData, "📁 ..")
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
    return fileData, fileIsDir
}

func makeHomeIconGrid(w fyne.Window) fyne.CanvasObject {
    type menuItem struct {
        icon fyne.Resource
        text string
        path string
    }
    items := []menuItem{
        {theme.HomeIcon(), "Internal", "/storage/emulated/0"},
        {theme.FolderIcon(), "Download", "/storage/emulated/0/Download"},
        {theme.FileImageIcon(), "DCIM", "/storage/emulated/0/DCIM"},
        {theme.VolumeUpIcon(), "Music", "/storage/emulated/0/Music"},
        {theme.SettingsIcon(), "Settings", ""},
    }

    var buttons []fyne.CanvasObject
    for _, m := range items {
        m := m
        btn := widget.NewButtonWithIcon(m.text, m.icon, func() {
            if m.path != "" {
                currentPath = m.path
                w.SetContent(makeMainUI(w))
            } else {
                dialog.ShowInformation("Settings", "Belum diimplementasi.", w)
            }
        })
        btn.Importance = widget.HighImportance
        buttons = append(buttons, btn)
    }
    
    grid := container.NewGridWithColumns(3, buttons...)
    return container.NewCenter(grid)
}
