package menu

import (
	"github.com/libretro/ludo/ludos"
	ntf "github.com/libretro/ludo/notifications"
)

type sceneWiFi struct {
	entry
}

func buildWiFi() Scene {
	var list sceneWiFi
	list.label = "WiFi Menu"

	list.children = append(list.children, entry{
		label: "Looking for networks",
		icon:  "reset",
	})

	list.segueMount()

	go func() {
		networks, err := ludos.ScanNetworks()
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
		}

		if len(networks) > 0 {
			list.children = []entry{}
			for _, network := range networks {
				network := network
				list.children = append(list.children, entry{
					label:       network.SSID,
					icon:        "wifi",
					stringValue: func() string { return ludos.NetworkStatus(network) },
					callbackOK: func() {
						list.segueNext()
						menu.Push(buildKeyboard(
							"Passphrase for "+network.SSID,
							func(pass string) {
								go func() {
									if err := ludos.ConnectNetwork(network, pass); err != nil {
										ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
									}
								}()
							},
						))
					},
				})
				list.segueMount()
				menu.tweens.FastForward()
			}
		} else {
			list.children[0].label = "No network found"
			list.children[0].icon = "close"
		}
	}()

	return &list
}

func (s *sceneWiFi) Entry() *entry {
	return &s.entry
}

func (s *sceneWiFi) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneWiFi) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneWiFi) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneWiFi) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneWiFi) render() {
	genericRender(&s.entry)
}

func (s *sceneWiFi) drawHintBar() {
	w, h := menu.Window.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 88*menu.ratio, 0, hintBgColor)
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 2*menu.ratio, 0, sepColor)

	_, upDown, _, a, b, _, _, _, _, _ := hintIcons()

	lstack := float32(75) * menu.ratio
	rstack := float32(w) - 96*menu.ratio
	stackHintLeft(&lstack, upDown, "Navigate", h)
	if s.children[0].callbackOK != nil {
		stackHintRight(&rstack, a, "Connect", h)
	}
	stackHintRight(&rstack, b, "Back", h)
}
