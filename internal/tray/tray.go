package tray

import (
	"fmt"
	"time"
	_ "embed"

	"fyne.io/systray"
)

type Options struct {
	StatusFn       func() (active bool, since time.Time, paused bool)
	OnPauseToggle  func()
	OnSleepNow     func()
	OnHibernateNow func()
	OnOpenLog      func()
	OnOpenConfig   func()
	OnExit         func()
}

//go:embed sleep_guardian.ico
var trayIcon []byte

func Run(o Options) {
	systray.Run(func() {
		systray.SetTitle("Sleep guardian")
		systray.SetIcon(trayIcon)

		mStatus := systray.AddMenuItem("Статус: …", "")
		mPause := systray.AddMenuItemCheckbox("Пауза", "Приостановить страж", false)
		mSleep := systray.AddMenuItem("Уснуть сейчас", "")
		mHiber := systray.AddMenuItem("Гибернация", "")
		systray.AddSeparator()
		mOpenLog := systray.AddMenuItem("Открыть лог", "")
		mOpenCfg := systray.AddMenuItem("Открыть конфиг", "")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Выход", "")

		go func() {
			for {
				a, since, paused := o.StatusFn()
				state := "Ожидание"
				if a {
					state = "Активно"
				}
				if paused {
					state += " (пауза)"
				}
				mStatus.SetTitle(fmt.Sprintf("Статус: %s • с %s", state, since.Format("15:04:05")))
				time.Sleep(1 * time.Second)
			}
		}()

		go func() {
			for {
				select {
				case <-mPause.ClickedCh:
					o.OnPauseToggle()
					mPause.Check()
				case <-mSleep.ClickedCh:
					o.OnSleepNow()
				case <-mHiber.ClickedCh:
					o.OnHibernateNow()
				case <-mOpenLog.ClickedCh:
					o.OnOpenLog()
				case <-mOpenCfg.ClickedCh:
					o.OnOpenConfig()
				case <-mQuit.ClickedCh:
					o.OnExit()
					systray.Quit()
					return
				}
			}
		}()
	}, func() {})
}
