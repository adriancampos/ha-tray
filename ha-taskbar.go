package main

// Todo:
// Allow for hot reloading config (maybe watch file?)
// In haws, make of HA's message id system
// I should be able to figure out domain without the user telling me. I /could/ strip off the first part of entity_id, but that sounds risky.
//     I'm pretty sure it's accessible in get_states, so I'll try that.
//
// Better Icon
// Allow for choice of standard or light/dark following system ui theme

import (
	"log"

	"github.com/adriancampos/ha-taskbar/config"
	"github.com/adriancampos/ha-taskbar/haws"
	"github.com/adriancampos/ha-taskbar/icon"
	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
)

var c *haws.Connection
var conf *config.Config

func main() {

	// Load config
	if _conf, err := config.LoadConfig("configuration.yaml"); err != nil {
		dlgs.Warning("Invalid Configuration", "configuration.yaml is either not found or invalid:\n"+err.Error())
		log.Fatalln(err)
	} else {
		conf = _conf
	}

var userEntityItems []entityItem
	systray.Run(onReady, nil)
	log.Println("Finished main(); exiting")
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("HomeAssistant")

	mURL := systray.AddMenuItem("Open Home Assistant", "Open web HA")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mURL.ClickedCh:
				browser.OpenURL(conf.Server.WebAddress)
			}
		}
	}()

	systray.AddSeparator()

	// Load entities

	type MItem struct {
		menuItem *systray.MenuItem
		entityID string
		onClick  func()
	}

	haEntities := make([]haws.ToggleableEntity, len(conf.ToggleableEntities))
	for i, te := range conf.ToggleableEntities {
		haEntities[i] = haws.ToggleableEntity{Domain: te.Domain, EntityID: te.EntityID}
	}
	log.Println("Loaded ToggleableEntities:", conf.ToggleableEntities)

	userMenuItems := make([]MItem, len(haEntities))
	for i, v := range haEntities {
		v := v
		i := i
		userMenuItems[i] = MItem{
			menuItem: systray.AddMenuItem(v.EntityID, "Toggle "+v.EntityID),
			entityID: v.EntityID,
			onClick: func() {
				haws.ToggleDevice(c, v)

				// To avoid UI flashes when the menu is next opened, assume that the toggle succeeded.
				// If the assumption is incorrect, it'll be corrected through the haws subscription.
				if userMenuItems[i].menuItem.Checked() {
					userMenuItems[i].menuItem.Uncheck()
				} else {
					userMenuItems[i].menuItem.Check()
				}
			},
		}
	}

	// Listen for menu close events
	go func() {
		for {
			<-systray.MenuClosedCh
			log.Println("systray menu hidden")
			// Close websocket connection when menu closes
			haws.Close(c)
		}
	}()

	// Listen for menu open events
	go func() {
		for {
			<-systray.MenuOpenedCh
			log.Println("systray menu shown")

			// Close any preexisting connection
			haws.Close(c)

			// Open connection to HA
			c = haws.OpenConnection(conf.Server.Address, conf.Server.AccessToken, func() { haws.RefreshAllEntities(c) }, func() {
				log.Println("Read error!")
			})

			// Handle clicks
			for _, mItem := range userMenuItems {
				mItem := mItem
				go func(mItem MItem) {
					for {
						// TODO: Should I add a switch case here and listen for menu closed events so that I can return so that there aren't a bunch of new threads?
						<-mItem.menuItem.ClickedCh
						mItem.onClick()
					}
				}(mItem)

				haws.SubscribeToggleableEntity(c, mItem.entityID, func(te haws.ToggleableEntity) {
					if te.State {
						mItem.menuItem.Check()
					} else {
						mItem.menuItem.Uncheck()
					}

					mItem.menuItem.SetTitle(te.FriendlyName)
				})

			}

		}
	}()
}
