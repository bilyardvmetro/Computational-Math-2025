package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Print("Выберете, что хотите решить (введите цифру)\n1) Решить нелинейное уравнение\n2) Решить систему нелинейных уравнений\n Enter: ")
	var option int
	ReadInt(in, &option, true)

	if option == 1 {
		LinearEquation(in, out)
	} else if option == 2 {
		LinearSystem(in, out)
	} else {
		GetOut(OptionError{})
	}
}
