package partitions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"unsafe"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

type FdiskCmd struct {
	Size int
	Unit string
	Path string
	Type string
	Fit  string
	Name string
}

func (cmd *FdiskCmd) AssignParameters(command utils.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "size" {
			cmd.Size = parameter.IntValue
		} else if parameter.Name == "unit" {
			cmd.Unit = parameter.StringValue
		} else if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		} else if parameter.Name == "type" {
			cmd.Type = parameter.StringValue
		} else if parameter.Name == "fit" {
			cmd.Fit = parameter.StringValue
		} else if parameter.Name == "name" {
			cmd.Name = parameter.StringValue
		}
	}
}

func (cmd *FdiskCmd) Fdisk() {
	if cmd.Size != -1 {
		if cmd.Name != "" {
			if cmd.Path != "" {
				// ARRAY PARA GUARDAR LOS ESPACIOS LIBRES
				var avalaibleSpaces = make([]utils.VoidSpace, 0)

				// ABRO EL ARCHIVO
				file, err := os.OpenFile(cmd.Path, os.O_RDWR, 0777)

				if err != nil {
					log.Fatal("Error ", err)
					return
				}
				// LEO EL MBR
				MBR := utils.MBR{}
				size := int(unsafe.Sizeof(MBR))
				file.Seek(0, 0)
				data := utils.ReadBytes(file, size)
				buffer := bytes.NewBuffer(data)
				err1 := binary.Read(buffer, binary.BigEndian, &MBR)
				if err1 != nil {
					log.Fatal("Error ", err1)
					return
				}

				for i := 0; i < 4; i++ {
					if utils.ByteToString(MBR.Partitions[i].Part_type[:]) == "e" {
						// SI ES UNA PARTICION EXTENDIDA BUSCO ENTRE TODAS SUS LOGICAS
						actualEbr := utils.EBR{}
						// LEO EL PRIMER EBR
						sizeEbr := int(unsafe.Sizeof(actualEbr))
						file.Seek(int64(utils.ByteToInt(MBR.Partitions[i].Part_start[:])), 0)
						dataEbr := utils.ReadBytes(file, sizeEbr)
						bufferebr := bytes.NewBuffer(dataEbr)
						errEbr := binary.Read(bufferebr, binary.BigEndian, &actualEbr)
						if errEbr != nil {
							log.Fatal("Error ", errEbr)
							return
						}
						if utils.ByteToString(actualEbr.Part_name[:]) == cmd.Name || utils.ByteToString(MBR.Partitions[i].Part_name[:]) == cmd.Name {
							fmt.Println("Error: la particion con el nombre " + cmd.Name + " ya existe en el disco " + cmd.Path)
							return
						}
						// RECORRO TODOS LOS EBR
						for utils.ByteToInt(actualEbr.Part_next[:]) != -1 {
							size := int(unsafe.Sizeof(actualEbr))
							file.Seek(int64(utils.ByteToInt(actualEbr.Part_next[:])), 0)
							data := utils.ReadBytes(file, size)
							buffer := bytes.NewBuffer(data)
							err1 := binary.Read(buffer, binary.BigEndian, &actualEbr)
							if err1 != nil {
								log.Fatal("Error ", err1)
								return
							}
							if utils.ByteToString(actualEbr.Part_name[:]) == cmd.Name {
								fmt.Println("Error: la particion con el nombre " + cmd.Name + " ya existe en el disco " + cmd.Path)
								return
							}
						}
					} else {
						if utils.ByteToString(MBR.Partitions[i].Part_name[:]) == cmd.Name {
							fmt.Println("Error: la particion con el nombre " + cmd.Name + " ya existe en el disco " + cmd.Path)
							return
						}
					}
				}

				// CALCULO DE MULTIPLICADOR PARA ASIGNAR ESPACIO A LA PARTICION
				multiplicator := 1024
				if cmd.Unit == "b" {
					multiplicator = 1
				} else if cmd.Unit == "m" {
					multiplicator = 1024 * 1024
				}
				//fmt.Println(multiplicator)

				// ORDENA TODAS LAS PARTICIONES DE MENOR A MAYOR
				utils.BubbleSort(MBR.Partitions[:])
				// CALCULO DE TIPO DE PARTICIONES
				totalPartitions := 0
				totalExtended := 0
				extendedPartition := utils.Partition{}

				for i := 0; i < 4; i++ {
					//fmt.Println(utils.ByteToString(MBR.Partitions[i].Part_size[:]), utils.ByteToString(MBR.Partitions[i].Part_start[:]))
					// CUENTA CUANTAS PARTICIONES HAY DENTRO DEL MBR
					if utils.ByteToString(MBR.Partitions[i].Part_size[:]) != "-1" {
						totalPartitions++
					}

					if utils.ByteToString(MBR.Partitions[i].Part_type[:]) == "e" {
						extendedPartition = MBR.Partitions[i]
						totalExtended++
					}
				}

				if cmd.Type == "" || cmd.Type == "p" || cmd.Type == "e" {

					/*for i := 0; i < 4; i++ {
						fmt.Println(utils.ByteToString(MBR.Partitions[i].Part_start[:]))
					}*/
					// CALCULO LOS ESPACIOS VACIOS ENTRE PARTICIONES PRIMARIAS Y EXTENDIDAS
					if totalPartitions == 4 {
						fmt.Println("Error: la particion " + cmd.Name + " no se pude crear porque la suma de particiones extendidas y primarias llego a su límite")
						return
					} else if totalExtended > 0 && cmd.Type == "e" {
						fmt.Println("Error: la particion " + cmd.Name + " no se pude crear porque solamente puede existir una particion extendida en el disco")
						return
					}

					if totalPartitions != 0 {
						for i := 0; i < 4; i++ {
							if i == 0 {
								tmpSpace := utils.VoidSpace{}
								tmpSpace.Size = utils.ByteToInt(MBR.Partitions[i].Part_start[:]) - int(unsafe.Sizeof(MBR)) - 2
								tmpSpace.Start = int(unsafe.Sizeof(MBR)) + 1
								avalaibleSpaces = append(avalaibleSpaces, tmpSpace)
							} else if i == 3 {
								tmpSpace := utils.VoidSpace{}
								tmpSpace.Size = utils.ByteToInt(MBR.Mbr_size[:]) - (utils.ByteToInt(MBR.Partitions[i].Part_size[:]) + utils.ByteToInt(MBR.Partitions[i].Part_start[:])) - 1
								tmpSpace.Start = utils.ByteToInt(MBR.Partitions[i].Part_size[:]) + utils.ByteToInt(MBR.Partitions[i].Part_start[:]) + 1
								avalaibleSpaces = append(avalaibleSpaces, tmpSpace)
							} else {
								tmpSpace := utils.VoidSpace{}
								tmpSpace.Size = utils.ByteToInt(MBR.Partitions[i].Part_start[:]) - (utils.ByteToInt(MBR.Partitions[i-1].Part_size[:]) + utils.ByteToInt(MBR.Partitions[i-1].Part_start[:])) - 2
								tmpSpace.Start = utils.ByteToInt(MBR.Partitions[i-1].Part_size[:]) + utils.ByteToInt(MBR.Partitions[i-1].Part_start[:]) + 1
								avalaibleSpaces = append(avalaibleSpaces, tmpSpace)
							}
						}
					} else {
						tmpSpace := utils.VoidSpace{}
						tmpSpace.Size = utils.ByteToInt(MBR.Mbr_size[:]) - (int(unsafe.Sizeof(MBR)) + 1)
						tmpSpace.Start = int(unsafe.Sizeof(MBR)) + 1
						avalaibleSpaces = append(avalaibleSpaces, tmpSpace)
					}
					/*fmt.Println("-------------------------------")
					for i := 0; i < len(avalaibleSpaces); i++ {
						fmt.Println(avalaibleSpaces[i])
					}
					fmt.Println("-------------------------------")*/
					// ORDENO LOS ESPACIOS VACIOS DE MENOR A MAYOR TAMANIO
					utils.SortFreeSpaces(avalaibleSpaces[:])

					// VARIABLE PARA GUARDAR DONDE INICIA LA PARTICION CREADA
					selectVoidSpace := -1
					for i := 0; i < len(avalaibleSpaces); i++ {
						//fmt.Println("Estoo", avalaibleSpaces[i].Size, cmd.Size*multiplicator, avalaibleSpaces[i].Start)
						if avalaibleSpaces[i].Size >= cmd.Size*multiplicator {
							selectVoidSpace = avalaibleSpaces[i].Start
							break
						}
					}
					if selectVoidSpace != -1 {
						for i := 0; i < 4; i++ {
							if utils.ByteToInt(MBR.Partitions[i].Part_size[:]) == -1 {
								fit := ""
								ptype := ""
								if cmd.Fit == "" {
									fit = "wf"
								}
								if cmd.Type == "" {
									ptype = "p"
								} else {
									ptype = cmd.Type
								}

								// SI ES UNA PARTICION EXTENDIDA CREO LA CABECERA
								if ptype == "e" {
									// ESCRIBO EL EBR INICIAL
									EBR := utils.EBR{}
									copy(EBR.Part_status[:], "0")
									copy(EBR.Part_fit[:], "wf")
									copy(EBR.Part_start[:], "-1")
									copy(EBR.Part_size[:], "-1")
									copy(EBR.Part_next[:], "-1")
									copy(EBR.Part_name[:], "")
									// ME POSICIONO AL INICIO DE LA PARTICION EXTENDIDA
									file.Seek(int64(selectVoidSpace), 0)
									// ESCRIBO EL PRIMER EBR EN EL DISCO
									var bufferControl bytes.Buffer
									binary.Write(&bufferControl, binary.BigEndian, &EBR)
									utils.WriteBytes(file, bufferControl.Bytes())
								}

								copy(MBR.Partitions[i].Part_name[:], []byte(cmd.Name))
								copy(MBR.Partitions[i].Part_fit[:], []byte(fit))
								copy(MBR.Partitions[i].Part_size[:], []byte(strconv.Itoa(cmd.Size*multiplicator)))
								copy(MBR.Partitions[i].Part_start[:], []byte(strconv.Itoa(selectVoidSpace)))
								copy(MBR.Partitions[i].Part_type[:], []byte(ptype))
								// IMPRIMO EN CONSOLA LOS DATOS DE LA PARTICION CREADA
								utils.PrintPartition(MBR.Partitions[i])
								break
							}
						}
					} else {
						fmt.Println("Error: la particion " + cmd.Name + " no cabe en el disco " + cmd.Path)
					}
				} else {
					if totalExtended == 1 {
						fit := ""
						if cmd.Fit == "" {
							fit = "wf"
						} else {
							fit = cmd.Fit
						}
						// EBR A ESCRIBIR
						EBR := utils.EBR{}
						// LEO EL PRIMER EBR
						tempEbr := utils.EBR{} // VARIABLE PARA RECORRER LA LISTA DE EBR
						size := int(unsafe.Sizeof(tempEbr))
						file.Seek(int64(utils.ByteToInt(extendedPartition.Part_start[:])), 0)
						data := utils.ReadBytes(file, size)
						bufferebr := bytes.NewBuffer(data)
						err1 := binary.Read(bufferebr, binary.BigEndian, &tempEbr)
						if err1 != nil {
							log.Fatal("Error ", err1)
							return
						}
						//VARIABLE PARA GUARDAR DONDE SE VA A ESCRIBIR EL EBR
						startToWrite := 0
						if utils.ByteToString(tempEbr.Part_start[:]) == "-1" {
							// CALCULA QUE LA PARTICION QUEPA
							if (utils.ByteToInt(extendedPartition.Part_size[:]) - int(unsafe.Sizeof(EBR))) >= (cmd.Size * multiplicator) {
								// CALCULA DONDE INICIA LA PRIMERA PARTICION LOGICA
								start := utils.ByteToInt(extendedPartition.Part_start[:]) + int(unsafe.Sizeof(EBR)) + 1
								startToWrite = utils.ByteToInt(extendedPartition.Part_start[:])
								// ASIGNA LOS VALORES AL PRIMER EBR
								copy(EBR.Part_status[:], []byte("0"))
								copy(EBR.Part_fit[:], []byte(fit))
								copy(EBR.Part_start[:], []byte(strconv.Itoa(start)))
								copy(EBR.Part_size[:], []byte(strconv.Itoa(cmd.Size*multiplicator)))
								copy(EBR.Part_next[:], []byte("-1"))
								copy(EBR.Part_name[:], []byte(cmd.Name))
							} else {
								fmt.Println("Error: la particion lógica " + cmd.Name + " no cabe en el disco " + cmd.Path)
								return
							}
						} else {
							for utils.ByteToInt(tempEbr.Part_next[:]) != -1 {
								size := int(unsafe.Sizeof(tempEbr))
								file.Seek(int64(utils.ByteToInt(tempEbr.Part_next[:])), 0)
								data := utils.ReadBytes(file, size)
								buffer := bytes.NewBuffer(data)
								err1 := binary.Read(buffer, binary.BigEndian, &tempEbr)
								if err1 != nil {
									log.Fatal("Error ", err1)
									return
								}
								//fmt.Println(utils.ByteToString(tempEbr.Part_name[:]), utils.ByteToString(tempEbr.Part_status[:]))
							}
							if ((utils.ByteToInt(extendedPartition.Part_start[:]) + utils.ByteToInt(extendedPartition.Part_size[:])) - (utils.ByteToInt(tempEbr.Part_start[:]) + utils.ByteToInt(tempEbr.Part_size[:]))) >= (cmd.Size * multiplicator) {
								// CALCULA DONDE INICIA LA PARTICION LOGICA
								start := utils.ByteToInt(tempEbr.Part_start[:]) + utils.ByteToInt(tempEbr.Part_size[:]) + int(unsafe.Sizeof(EBR)) + 2
								// CALCULA DONDE SE VA A ESCRIBIR EL EBR
								startToWrite = utils.ByteToInt(tempEbr.Part_start[:]) + utils.ByteToInt(tempEbr.Part_size[:]) + 1
								// ASIGNA LOS VALORES AL PRIMER EBR
								copy(EBR.Part_status[:], []byte("0"))
								copy(EBR.Part_fit[:], []byte(fit))
								copy(EBR.Part_start[:], []byte(strconv.Itoa(start)))
								copy(EBR.Part_size[:], []byte(strconv.Itoa(cmd.Size*multiplicator)))
								copy(EBR.Part_next[:], []byte("-1"))
								copy(EBR.Part_name[:], []byte(cmd.Name))
								// APUNTO EL SIGUIENTE DE TEMPORAL AL INICIO DEL ACTUAL CREADO
								copy(tempEbr.Part_next[:], []byte(strconv.Itoa(startToWrite)))
								// REESCRIBO EL EBR TEMPORAL
								//fmt.Println(cmd.Name, startToWrite)
								//fmt.Println(utils.ByteToString(tempEbr.Part_name[:]), utils.ByteToInt(tempEbr.Part_start[:])-(int(unsafe.Sizeof(tempEbr))+1))
								file.Seek(int64(utils.ByteToInt(tempEbr.Part_start[:])-(int(unsafe.Sizeof(tempEbr))+1)), 0)
								var bufferControlTemp bytes.Buffer
								binary.Write(&bufferControlTemp, binary.BigEndian, &tempEbr)
								utils.WriteBytes(file, bufferControlTemp.Bytes())
							} else {
								fmt.Println("Error: la particion lógica " + cmd.Name + " no cabe en el disco " + cmd.Path)
								return
							}
						}
						// ESCRIBO EL EBR EN EL DISCO
						file.Seek(int64(startToWrite), 0)
						var bufferControl bytes.Buffer
						binary.Write(&bufferControl, binary.BigEndian, &EBR)
						utils.WriteBytes(file, bufferControl.Bytes())
						// IMPRIMO LOS DATOS DEL EBR CREADO
						utils.PrintEBR(EBR)
					} else {
						fmt.Println("Error: la particion " + cmd.Name + " no puede ser creada debido no existe particion extendida")
						return
					}
				}
				// REESCRIBO EL MBR EN EL DISCO
				file.Seek(0, 0)
				var bufferControl bytes.Buffer
				binary.Write(&bufferControl, binary.BigEndian, &MBR)
				utils.WriteBytes(file, bufferControl.Bytes())

				/*for i := 0; i < 4; i++ {
					fmt.Println(utils.ByteToString(MBR.Partitions[i].Part_start[:]))
				}*/
				// CIERRO EL ARCHIVO
				file.Close()
			} else {
				fmt.Println("Error: el parametro path es obligatorio en fdisk")
			}
		} else {
			fmt.Println("Error: el parametro name es obligatorio en fdisk")
		}
	} else {
		fmt.Println("Error: el parametro size es obligatorio en fdisk")
	}
}
