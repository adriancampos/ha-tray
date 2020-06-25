# HA Taskbar

A taskbar app to toggle Home Assistant entities.  
Currently supports turning available switches, lights and input_boolean's on and off. 


## Installation
Grab the latest executable from [Releases](https://github.com/adriancampos/ha-taskbar/releases). I've only built and tested a Windows binary. More to come in the future. For the time being, please build the project if you're on another platform. 


## Setup
For now, all configuration is done through `configuration.yaml`. To start, you'll need to include your **server address**, **long lived access token** and some **entities**.

### Token
To create a token within HA, login to HA and click on your profile.
Under Long Lived Access Tokens, create a new token, give it a name place this token into the config.


### Entities
Entity configuration should look similar to that of Home Assistant. Each entity should have a(n):
* entity_id. This needs to match your HA config exactly.
* device_type. This is required. It will be pulled from HA in the future.

### Example config
```yaml
server:
    server_address: demo.home-assistant.io:80
    web_address: https://demo.home-assistant.io  # optional. Default is "https://" + server_address
    access_token: long-lived-page-access-token-abcdef-123456

entities:
  - entity_id: light.led_strip
    domain: light
  - entity_id: light.desk
    domain: light
  - entity_id: light.ceiling
    domain: light
  - entity_id: switch.plug01
    domain: switch
```

## Building
TODO

## TODO
- [ ] Support scenes, input_booleans, automations, scripts, etc.
- [ ] Allow the user to customize menu order, dividers, submenus
- [ ] Pull entity icon and place on menu item 
- [x] Allow hot reloading of config
- [ ] Infer domain from websocket messages
- [ ] Refactor haws to make use of HA's message id system
- [ ] Optional light/dark icon to match taskbar aesthetic
- [ ] Better error checking for config
