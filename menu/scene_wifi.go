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
		icon:  "reload",
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
					icon:        "menu_network",
					stringValue: func() string { return ludos.NetworkStatus(network) },
					callbackOK: func() {
						list.segueNext()
						menu.Push(buildKeyboard(
							"Passpharse for "+network.SSID,
							func(pass string) {
								go func() {
									err := ludos.ConnectNetwork(network, pass)
									if err != nil {
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
	HintBar(&Props{},
		Hint(&Props{}, "key-up-down", "NAVIGATE"),
		Hint(&Props{}, "key-x", "BACK"),
		Hint(&Props{Hidden: s.children[0].callbackOK == nil}, "key-x", "CONNECT"),
	)()
}
