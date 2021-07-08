/**
 * Copyright (c) 2021 BlockDev AG
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

const (
	docker = "docker"
	group  = "docker-users"
)

func checkSystemsAndTry() {
	mod.Invalidate()
	dckr := os.Getenv("ProgramFiles") + "\\Docker\\Docker\\resources\\bin\\" + docker

	for {
		ex := cmdRun(mod.lv, dckr, []string{"ps"})
		switch ex {
		case 0:
			mod.lbDocker.SetText("Running [OK]")

			ex := cmdRun(mod.lv, dckr, strings.Split("container start myst", " "))
			switch ex {
			case 0:
				mod.lbContainer.SetText("Running [OK]")
				mod.btnCmd2.SetEnabled(true)

			default:
				log.Printf("Failed to start cmd: %v", ex)
				mod.lbContainer.SetText("Installing")

				ex := cmdRun(mod.lv, dckr, strings.Split("run --cap-add NET_ADMIN -d -p 4449:4449 --name myst -v myst-data:/var/lib/mysterium-node mysteriumnetwork/myst:latest service --agreed-terms-and-conditions", " "))
				if ex == 0 {
					mod.lbDocker.SetText("Running [OK]")
					continue
				}
			}

		case 1:
			mod.lbDocker.SetText("Starting..")
			mod.lbContainer.SetText("-")

			if isProcessRunning("Docker Desktop.exe") {
				break
			}
			dd := os.Getenv("ProgramFiles") + "\\Docker\\Docker\\Docker Desktop.exe"
			cmd := exec.Command(dd, "-Autostart")
			if err := cmd.Start(); err != nil {
				log.Printf("Failed to start cmd: %v", err)
			}
			break

		default:
			mod.SetState(ST_INSTALL_NEED)
			mod.WaitDialogueComplete()
			mod.SetState(ST_INSTALL_INPROGRESS)

			if !CheckWindowsVersion() {
				mod.lbInstallationState2.SetText("Reason:\r\nYou must be running Windows 10 version 1607 (the Anniversary update) or above.")
				mod.SetState(ST_INSTALL_ERR)
				mod.WaitDialogueComplete()

				// exit
				win.SendMessage(mod.mw.Handle(), win.WM_CLOSE, 0, 0)
				return
			}

			list := []struct{ url, name string }{
				{"https://desktop.docker.com/win/stable/amd64/Docker%20Desktop%20Installer.exe", "DockerDesktopInstaller.exe"},
				{"https://wslstorestorage.blob.core.windows.net/wslblob/wsl_update_x64.msi", "wsl_update_x64.msi"},
			}
			for fi, v := range list {
				if _, err := os.Stat(os.Getenv("TMP") + v.name); err != nil {

					mod.lbInstallationState2.SetText(fmt.Sprintf("%d of %d: %s", fi+1, len(list), v.name))
					mod.PrintProgress(0)

					err := DownloadFile(os.Getenv("TMP")+"\\"+v.name, v.url, mod.PrintProgress)
					if err != nil {
						mod.lbInstallationState2.SetText("Reason:\r\nDownload failed")
						mod.SetState(ST_INSTALL_ERR)
						mod.WaitDialogueComplete()
						return
					}
				}
			}

			err := runMeElevated("msiexec.exe", "/I wsl_update_x64.msi /quiet", os.Getenv("TMP"))
			if err != nil {
				mod.lbInstallationState2.SetText("Reason:\r\nCommand failed: msiexec.exe /I wsl_update_x64.msi")
				mod.SetState(ST_INSTALL_ERR)
				mod.WaitDialogueComplete()
				return
			}
			ex := cmdRun(mod.lv, os.Getenv("TMP")+"\\DockerDesktopInstaller.exe", []string{"install", "--quiet"})
			if ex != 0 {
				mod.lbInstallationState2.SetText("Reason:\r\nDockerDesktopInstaller failed")
				mod.SetState(ST_INSTALL_ERR)
				mod.WaitDialogueComplete()
				return
			}

			if !checkExe() {
				installExe()
			}
			if !CurrentGroupMembership(group) {
				// request to logout //

				ret := walk.MsgBox(mod.mw, "Installation", "Log of from the current session to finish the installation.", walk.MsgBoxTopMost|walk.MsgBoxYesNo|walk.MsgBoxIconExclamation)
				if ret == win.IDYES {
					windows.ExitWindowsEx(windows.EWX_LOGOFF, 0)
					return
				}
				mod.SetState(ST_INSTALL_ERR)
				mod.lbInstallationState2.SetText("Log of from the current session to finish the installation.")
				mod.WaitDialogueComplete()
				return
			}

			mod.SetState(ST_INSTALL_FIN)
			mod.WaitDialogueComplete()
			mod.SetState(ST_STATUS_FRAME)
			continue
		}

		time.Sleep(10 * time.Second)
	}
}
