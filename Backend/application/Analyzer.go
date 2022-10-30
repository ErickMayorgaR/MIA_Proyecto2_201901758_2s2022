package application

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

type Comment struct {
	Value string
}

func AnalyzerF(script_path string, isFile bool) []utils.Command {
	var tree = make([]utils.Command, 0)
	input := ""
	// script o comando
	if isFile {
		// leer si es archivo
		input = strings.ToLower(readFile(script_path)) + "\n"
	} else {

		input = script_path
	}

	tempCommand := newCommand("")
	tempPar := newParameter("", "", -1)
	tempWord := ""

	findValue := false

	tempIntValue := -1
	tempStringValue := ""
	isIntValue := false

	isComment := false

	valueFound := false
	catchSpaces := false

	for i, character := range input {
		letter := string(character)
		// si ha iniciado un comentario
		if letter == "#" {
			tree = append(tree, tempCommand)
			tempCommand = newCommand("comment")
			isComment = true
		}

		if !isComment {
			if !findValue {
				if letter != " " && letter != "\n" {
					tempWord += letter
					if tempWord == "mkdisk" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "rmdisk" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "fdisk" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mount" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mkfs" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mkdir" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mkfile" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "rep" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "login" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "logout" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mkgrp" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "rmgrp" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mkuser" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "rmusr" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "exec" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "pause" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
						// PARAMETROS
					} else if tempWord == "-size=" {
						tempPar = newParameter("size", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = true
					} else if tempWord == "-fit=" {
						tempPar = newParameter("fit", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-unit=" {
						tempPar = newParameter("unit", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-path=" {
						tempPar = newParameter("path", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-ruta=" {
						tempPar = newParameter("ruta", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-type=" {
						tempPar = newParameter("type", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-name=" {
						tempPar = newParameter("name", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-id=" {
						tempPar = newParameter("id", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-cont=" {
						tempPar = newParameter("cont", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-usuario=" {
						tempPar = newParameter("usuario", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-password=" {
						tempPar = newParameter("password", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-pwd=" {
						tempPar = newParameter("pwd", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-grp=" {
						tempPar = newParameter("grp", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
						// parametros de una sola letra
					} else if tempWord == "-p" {
						// si llego al final de la cadena entonces guarda el parametro y el comando
						if i == (len(input) - 1) {
							tempPar = newParameter("-p", "-p", -1)
							tempWord = ""
							tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
							// guardo comando en el arbol de comandos
							tree = append(tree, tempCommand)
						} else {
							// si el siguiente es un espacio o un "-" guarda el parametro
							if string(input[i+1]) == " " || string(input[i+1]) == "-" || string(input[i+1]) == "\n" || string(input[i+1]) == "#" {
								tempPar = newParameter("-p", "-p", -1)
								tempWord = ""
								tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
								// si no es un caracter de separacion continua analizando
							} else {
								continue
							}
						}
					} else if tempWord == "-r" {
						// si llego al final de la cadena entonces guarda el parametro y el comando
						if i == (len(input) - 1) {
							tempPar = newParameter("-r", "-r", -1)
							tempWord = ""
							tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
							// guardo comando en el arbol de comandos
							tree = append(tree, tempCommand)
						} else {
							// si el siguiente es un espacio o un "-" guarda el parametro
							if string(input[i+1]) == " " || string(input[i+1]) == "-" || string(input[i+1]) == "\n" || string(input[i+1]) == "#" {
								tempPar = newParameter("-r", "-r", -1)
								tempWord = ""
								tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
								// si no es un caracter de separacion continua analizando
							} else {
								continue
							}
						}
					}
				}
			} else {
				tempWord += letter
				// bandera si viene cadena con comillas
				if letter == "\"" {
					if catchSpaces {
						catchSpaces = false
					} else {
						catchSpaces = true
					}
				}

				if letter == " " && catchSpaces && valueFound {
					tempWord += letter
				}
				if letter == " " && !valueFound {
					continue
				} else {
					valueFound = true
				}
				// si hay un espacio se siguen buscando comandos
				if (letter == " " || i == (len(input)-1) || letter == "\n") && valueFound && !catchSpaces {
					if isIntValue {
						convertedValue, _ := strconv.Atoi(strings.TrimSpace(tempWord))
						tempIntValue = convertedValue
					} else {
						// String sin saltos
						tempStringValue = strings.TrimSuffix(strings.TrimSpace(tempWord), "\n")
					}
					tempPar.IntValue = tempIntValue
					tempPar.StringValue = tempStringValue
					tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
					tempWord = ""
					tempStringValue = ""
					tempIntValue = -1
					findValue = false
					valueFound = false // != " "

					if i == (len(input) - 1) {
						tree = append(tree, tempCommand)
					}
				}
			}
		} else {
			tempWord += letter
			// mientras no sea salto de linea o no se haya llegado al final de la cadena, se concatena
			if letter == "\n" || i == (len(input)-1) {
				tempPar = newParameter("value", strings.TrimSpace(tempWord), -1)
				tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
				tempWord = ""
				isComment = false
				// Si llega al final se almacena el comando
				if i == (len(input) - 1) {
					tree = append(tree, tempCommand)
				}
			}
		}

		// si llega al final, se almacena el comando
		if i == (len(input) - 1) {
			tree = append(tree, tempCommand)
		}
	}
	return tree
}

func newParameter(name string, stringValue string, intValue int) utils.Parameter {
	return utils.Parameter{Name: name, StringValue: stringValue, IntValue: intValue}
}

func newCommand(name string) utils.Command {
	temp := make([]utils.Parameter, 1)
	return utils.Command{Parameters: temp, Name: name}
}

func readFile(script_path string) string {
	datosComoBytes, err := ioutil.ReadFile(script_path /*"./test.txt"*/)
	if err != nil {
		log.Fatal(err)
	}
	// convirtiendo bytes a string
	datosComoString := string(datosComoBytes)

	return datosComoString
}

func (cmd *Comment) AssignParameters(command utils.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "value" {
			cmd.Value = parameter.StringValue
		}
	}
}

func (cmd *Comment) ShowComment() {
	//fmt.Println(cmd.Value)
}
