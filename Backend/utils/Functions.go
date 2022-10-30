package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func WriteBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func GetDate() string {

	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return fecha
}

func GetRandom() int {
	res := rand.Intn(319)
	return res
}

func ByteToString(byteArray []byte) string {
	temp := bytes.Trim(byteArray, "\x00")
	str := string(temp[:])
	return str
}

func ByteToInt(byteArray []byte) int {
	return ToInt(ByteToString(byteArray))
}

func ToInt(stringValue string) int {
	res, err := strconv.Atoi(stringValue)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func BubbleSort(array []Partition) {
	for i := 0; i < len(array)-1; i++ {
		for j := 0; j < len(array)-i-1; j++ {
			if ToInt(ByteToString(array[j].Part_start[:])) > ToInt(ByteToString(array[j+1].Part_start[:])) {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
}

func SortFreeSpaces(array []VoidSpace) {
	for i := 0; i < len(array)-1; i++ {
		for j := 0; j < len(array)-i-1; j++ {
			if array[j].Size > array[j+1].Size {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
}

// FUNCION PARA BORRAR ELEMENTO POR INIDICE DE UN ARRAY
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// FUNCION PARA OBTENER EL CONTENIDO DE UN ARCHIVO
func ReadFile(file_name string) string {
	// LEER EL ARREGLO DE BYTES DEL ARCHIVO
	datosComoBytes, err := ioutil.ReadFile(file_name)
	if err != nil {
		log.Fatal(err)
	}
	// CONVERTIR EL ARREGLO A STRING
	datosComoString := string(datosComoBytes)

	return datosComoString
}

func GraphDot(dot_path string, report_path string) {
	// OBTENGO LA EXTENSION
	extension := report_path[strings.LastIndex(report_path, ".")+1:]
	// CREA GRAFICA A TRAVES DE COMANDOS EN CONSOLA
	ver, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(ver, "-T"+extension, dot_path).Output()
	mode := int(0777)
	ioutil.WriteFile(report_path, cmd, os.FileMode(mode))
}

// FUNCION PARA DETERMINAR LA CANTIDAD DE BLOQUES A USAR DEPENDIENDO
// DEL TAMANIO DE UNA CADENA
func GetBlocksNumber(size int) int {

	var number_blocks int
	if size > 64 {
		number_blocks = (size / 64)
	} else {
		number_blocks = 1
	}

	if size == 0 {
		number_blocks = 0
	}

	if ((size % 64) != 0) && (size > 64) {
		number_blocks++
	}
	return number_blocks
}

func PrintEBR(ebr EBR) {
	fmt.Println("---------EBR CREADO-----------")
	fmt.Println("Name: ", ByteToString(ebr.Part_name[:]))
	fmt.Println("Fit: ", ByteToString(ebr.Part_fit[:]))
	fmt.Println("Next: ", ByteToString(ebr.Part_next[:]))
	fmt.Println("Size: ", ByteToString(ebr.Part_size[:]))
	fmt.Println("Start: ", ByteToString(ebr.Part_start[:]))
	fmt.Println("Status: ", ByteToString(ebr.Part_status[:]))
	fmt.Println("------------------------------")
}

func PrintPartition(partition Partition) {
	fmt.Println("------ PARTICION CREADA ------")
	fmt.Println("Name: ", ByteToString(partition.Part_name[:]))
	fmt.Println("Fit: ", ByteToString(partition.Part_fit[:]))
	fmt.Println("Size: ", ByteToString(partition.Part_size[:]))
	fmt.Println("Start: ", ByteToString(partition.Part_start[:]))
	fmt.Println("Status: ", ByteToString(partition.Part_status[:]))
	fmt.Println("Type: ", ByteToString(partition.Part_type[:]))
	fmt.Println("------------------------------")
}

func NewByte() [64]byte {
	var ver [64]byte
	return ver
}
