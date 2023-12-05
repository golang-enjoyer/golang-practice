package main

import (
	"fmt"
	"strings"
)

var globalWorld World

var globalPlayer Player

type Location struct {
	name           string
	ways           []string
	items          map[string][]string
	objectsToApply map[string]string
	description    string
	blockers       map[string]string
}

type Player struct {
	inventoryMap map[string][]string
	applyMap     map[string]string
	inventory    []string
	tasks        []string
}

type World struct {
	Locations       []Location
	currentLocation Location
	Player          Player
}

func initGame() {
	player := Player{
		inventoryMap: map[string][]string{
			"ключи":     {"рюкзак"},
			"конспекты": {"рюкзак"},
			"рюкзак":    {},
		},
		applyMap: map[string]string{
			"ключи": "дверь",
		},
		tasks: []string{
			"собрать рюкзак",
			"идти в универ",
		},
	}

	kitchen := Location{
		name: "кухня",
		ways: []string{"коридор"},
		items: map[string][]string{
			"стол": {"чай"},
		},
		description: "ты находишься на кухне, ",
	}

	corridor := Location{
		name:        "коридор",
		ways:        []string{"кухня", "комната", "улица"},
		items:       make(map[string][]string),
		description: "ничего интересного. ",
		objectsToApply: map[string]string{
			"дверь": "ключи",
		},
	}

	room := Location{
		name: "комната",
		ways: []string{"коридор"},
		items: map[string][]string{
			"стол": {"ключи", "конспекты"},
			"стул": {"рюкзак"},
		},
		description: "ты в своей комнате. ",
	}

	street := Location{
		name:        "улица",
		ways:        []string{"домой"},
		items:       make(map[string][]string),
		description: "на улице весна. ",
	}

	world := World{
		Locations: []Location{
			kitchen,
			corridor,
			room,
			street,
		},
		currentLocation: kitchen,
		Player:          player,
	}

	globalWorld = world
	globalPlayer = player
}

func addItems() string {
	itemsDescription := ""

	for k, v := range globalWorld.currentLocation.items {
		if len(v) != 0 {
			itemsDescription += "на " + k + "e: " + strings.Join(v, ", ") + ", "
		}
	}

	return itemsDescription
}

func addTasks() string {
	return "надо " + strings.Join(globalPlayer.tasks, " и ") + ". "
}

func possibleWays() string {
	return "можно пройти - " + strings.Join(globalWorld.currentLocation.ways, ", ")
}

func composeDescription() string {
	description := globalWorld.currentLocation.description
	description += addItems()
	description += addTasks()
	description += possibleWays()

	return description
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func findLocationInSlice(target string, slice []Location) Location {
	for _, element := range slice {
		if element.name == target {
			return element
		}
	}
	return Location{}
}

func findElementInSlice(target string, slice []string) (string, bool) {
	for _, element := range slice {
		if element == target {
			return element, true
		}
	}
	return "", false
}

func checkLocationToGo(location string) string {
	if stringInSlice(location, globalWorld.currentLocation.ways) {
		globalWorld.currentLocation = findLocationInSlice(location, globalWorld.Locations)
		return location + ", ничего интересного. " + possibleWays()
	} else {
		return "нет пути в " + location
	}
}

func canTakeItem(item string) bool {
	return len(globalPlayer.inventoryMap[item]) == 0
}

func addItemIfExists(itemToCheck string) bool {
	result := false
	for key, slice := range globalWorld.currentLocation.items {
		var newSlice []string
		for _, item := range slice {
			if item != itemToCheck {
				newSlice = append(newSlice, item)
			} else {
				globalPlayer.inventory = append(globalPlayer.inventory, itemToCheck)
				result = true
			}
		}
		globalWorld.currentLocation.items[key] = newSlice
	}

	return result
}

func takeItem(item string) string {
	if canTakeItem(item) {
		if addItemIfExists(item) {
			return "предмет добавлен в инвентарь: " + item
		} else {
			return "нет такого"
		}
	} else if globalPlayer.inventoryMap[item][0] == "рюкзак" {
		return "некуда класть"
	}

	return ""
}

func applyItem(item string, object string) string {
	if _, ok := findElementInSlice(item, globalPlayer.inventory); ok {
		// check that can apply
		if _, ok := globalWorld.currentLocation.objectsToApply[object]; ok {
			return "применено: " + item
		} else {
			return "не к чему применить"
		}
	} else {
		return "нет предмета в инвентаре - " + item
	}
}

func putOn(itemToCheck string) string {
	if _, ok := findElementInSlice(itemToCheck, globalWorld.currentLocation.items["стул"]); ok {
		globalPlayer.inventory = append(globalPlayer.inventory, itemToCheck)

		for k, v := range globalPlayer.inventoryMap {
			var newSlice []string
			for _, item := range v {
				if item != itemToCheck {
					newSlice = append(newSlice, item)
				}
			}

			for key, slice := range globalWorld.currentLocation.items {
				var newSlice []string
				for _, item := range slice {
					if item != itemToCheck {
						newSlice = append(newSlice, item)
					}
				}
				globalWorld.currentLocation.items[key] = newSlice
			}
			globalPlayer.inventoryMap[k] = newSlice
		}

		return "вы надели: " + itemToCheck
	} else {
		return "нет такого"
	}
}

func handleCommand(command string) string {
	splittedCommand := strings.Fields(command)

	switch splittedCommand[0] {
	case "осмотреться":
		return composeDescription()
	case "идти":
		return checkLocationToGo(splittedCommand[1])
	case "применить":
		return applyItem(splittedCommand[1], splittedCommand[2])
	case "взять":
		return takeItem(splittedCommand[1])
	case "надеть":
		return putOn(splittedCommand[1])
	default:
		return "неизвестная команда"
	}
}

func main() {
	initGame()
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("завтракать"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("надеть рюкзак"), "xxx")
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("взять телефон"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять конспекты"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти улица"))
	fmt.Println(handleCommand("применить ключи дверь"))
}
