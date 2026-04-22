package main

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "os/exec"
    "log"
)

func main() {
    a := app.NewWithID("com.omnifile.pro")
    w := a.NewWindow("OmniFile Pro")

    // Tab 1: Explorer
    explorer := widget.NewLabel(" Internal Storage")

    // Tab 2: Kompres Video
    kompres := widget.NewButton("Kompres Video 100MB  8MB", func() {
        cmd := exec.Command("ffmpeg", "-i", "/sdcard/video.mp4", "-crf", "28", "/sdcard/output.mp4")
        if err := cmd.Run(); err != nil {
            log.Println("Error:", err)
        }
    })

    // Tab 3: Duplikat
    duplikat := widget.NewButton("Hapus File Duplikat", func() {
        log.Println("Scan duplikat...")
    })

    // Tab 4: Tersembunyi
    tersembunyi := widget.NewButton("Lihat File Tersembunyi", func() {
        log.Println("List .hidden files...")
    })

    // Tab 5: Pembersih
    pembersih := widget.NewButton("Bersihkan Cache", func() {
        log.Println("Hapus cache...")
    })

    // Tab 6: FTP
    ftp := widget.NewButton("Jalanin FTP Server", func() {
        log.Println("FTP on port 2121...")
    })

    tabs := container.NewAppTabs(
        container.NewTabItem("Explorer", explorer),
        container.NewTabItem("Kompres", kompres),
        container.NewTabItem("Duplikat", duplikat),
        container.NewTabItem("Tersembunyi", tersembunyi),
        container.NewTabItem("Pembersih", pembersih),
        container.NewTabItem("FTP", ftp),
    )

    w.SetContent(tabs)
    w.ShowAndRun()
}
