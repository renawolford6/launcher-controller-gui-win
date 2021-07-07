/**
 * Copyright (c) 2021 BlockDev AG
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
package main

import (
	"fmt"
	"net"

	"github.com/lxn/win"

	"github.com/lxn/walk"
)

type Model struct {
	state        int
	inTray       bool
	pipeListener net.Listener

	icon *walk.Icon
	mw   *walk.MainWindow
	lv   *LogView

	// docker
	lbDocker    *walk.Label
	lbContainer *walk.Label

	// inst
	lbInstallationState *walk.Label
	progressBar         *walk.ProgressBar

	// common
	btnCmd  *walk.PushButton
	btnCmd2 *walk.PushButton

	dlg chan int
}

const (
	STATUS_FRAME = 0
	INSTALL_NEED = -1
	INSTALL_     = -2
	INSTALL_FIN  = -3
)

var mod Model

func init() {
	mod.dlg = make(chan int)
}

func (m *Model) ShowMain() {
	m.mw.Show()
	win.BringWindowToTop(m.mw.Handle())

	win.ShowWindow(m.mw.Handle(), win.SW_SHOW)
	win.ShowWindow(m.mw.Handle(), win.SW_SHOWNORMAL)
	//win.ShowWindow(m.mw.Handle(), win.SW_RESTORE)

	//win.SetWindowPos(m.mw.Handle(), win.HWND_NOTOPMOST, 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE)
	//win.SetWindowPos(m.mw.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE)
	//win.SetWindowPos(m.mw.Handle(), win.HWND_NOTOPMOST, 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE)
}

func (m *Model) SetState(s int) {
	m.state = s
	m.Invalidate()
}

const frameI = 1
const frameS = 2

func (m *Model) Invalidate() {
	if m.state == 0 {
		m.mw.Children().At(frameI).SetVisible(false)
		m.mw.Children().At(frameS).SetVisible(true)

	}
	if m.state == INSTALL_NEED {
		m.mw.Children().At(frameI).SetVisible(true)
		m.mw.Children().At(frameS).SetVisible(false)
		m.HideProgress()

		m.btnCmd.SetEnabled(true)
		m.btnCmd.SetText("Install")
		m.btnCmd.SetFocus()

		m.lbInstallationState.SetText("Docker desktop is requred to run exit node.\r\n" +
			"Press button to begin installation.")

		m.lbDocker.SetText("OK")
	}
	if m.state == INSTALL_ {
		m.btnCmd.SetEnabled(false)
		m.lbInstallationState.SetText("Downloading installation packages.\r\n" + "_")
	}
	if m.state == INSTALL_FIN {
		m.lbInstallationState.SetText("Installation successfully finished!\r\n_")
		m.btnCmd.SetEnabled(true)
		m.btnCmd.SetText("Finish !")
	}
}

func (m *Model) BtnOnClick() {
	fmt.Println("BtnOnClick", m.state)

	if m.state == INSTALL_FIN {
		//m.SetState(0)
		m.dlg <- 0
	}
	if m.state == INSTALL_NEED {
		m.SetState(INSTALL_)
		m.dlg <- 0
	}
}

func (m *Model) WaitDialogueComplete() {
	//log.Println("WaitDialogueComplete>")
	<-m.dlg
}

func (m *Model) HideProgress() {
	m.progressBar.SetVisible(false)
}

func (m *Model) PrintProgress(progress int) {
	m.lv.AppendText(fmt.Sprintf("Download %d %%\r\n", progress))
	m.progressBar.SetVisible(true)
	m.progressBar.SetValue(progress)
}
