// Question: Implement a remote control system for a home automation system where each button on the remote can be programmed to perform different actions (e.g., turn on lights, play music).
// Scenario: Create Command objects for each action and a RemoteControl class that invokes these commands.

package main

import "fmt"

// Command interface
type command interface{
	execute()
}

// Light struct
type Light struct{}

func (L *Light) turnOn(){
	fmt.Println("Light is turned ON")
}

type LightOnCommand struct{
	light *Light
}

func (Loc *LightOnCommand) execute(){
	Loc.light.turnOn()
}

type RemoteControl struct{
	command command
}

func (Rc *RemoteControl) pressButton(){
	Rc.command.execute()
}

func main(){
	Light := &Light{}
	LightOnCommand := &LightOnCommand{light: Light}
	RemoteControl := &RemoteControl{command: LightOnCommand}
	RemoteControl.pressButton()

}