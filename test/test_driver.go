package main

import (
	"../driver"
	"fmt"

)






func main() {


	tmp := driver.Elev_Init()

	switch tmp {
	case 1:
		fmt.Println("Elevator initiated succesfully")
	case 0:
		fmt.Println("Error during elevator initiation!")
	}


	driver.Elev_Set_Motor_Direction(-1)






}
