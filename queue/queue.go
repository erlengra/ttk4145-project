// Package queue defines and stores queues for the lift. Two queues are used:
// One local queue containing all orders assigned to this particular lift,
// and one remote containing all external orders assigned to any lift on the
// network. The remote order also stores information about the IP of the lift
// assigned to each order, and it has a timer attached to each order. The
// timer makes sure that if an order is assigned to a lift, and we never
// receive an 'order complete' message, the order will still be handled.
package queue

import (
	def "../config"
	"../network"
	"fmt"
	"log"
	"time"
	"strconv"
)

// queue defines a queue, a 2D array of orderStatuses representing the
// buttons on the lift panel.
type queue struct {
	matrix [def.NUM_FLOORS][def.NUM_BUTTONS]orderStatus
}

// orderStatus defines the status of an order: Whether it is active, which
// lift is assigned to take it, and how long it has been active. (The latter
// two are only used in the remote queue.)
type orderStatus struct {
	active bool
	addr   string      `json:"-"`
	timer  *time.Timer `json:"-"`
}

var inactive = orderStatus{active: false, addr: "", timer: nil}

var local queue
var remote queue

var updateLocal = make(chan bool)
var takeBackup = make(chan bool, 10)
var OrderTimeoutChan = make(chan def.OrderButton)
var newOrder chan bool
//var newOrder = make(chan bool)

func Init(newOrderTemp chan bool, outgoingMsg chan network.Packet) {
	newOrder = newOrderTemp
	go updateLocalQueue()
	//runBackup(outgoingMsg)
	log.Println(def.ColG, "Queue initialised.", def.ColN)
}

// AddLocalOrder adds an order to the local queue.
func AddLocalOrder(floor int, button def.ButtonType) {
	//fmt.Println("The order was added locally?")
	local.setOrder(floor, button, orderStatus{true, "", nil})

	newOrder <- true
	//fmt.Println("HEIR")

}

// AddRemoteOrder adds an order to the remote queue, and spawns a timer
// for the order. (If the order times out, it will be taken care of.)
func AddRemoteOrder(floor int, button def.ButtonType, addr string) {
	alreadyExist := IsRemoteOrder(floor, button)
	remote.setOrder(floor, button, orderStatus{true, addr, nil})
	if !alreadyExist {
		go remote.startTimer(floor, button)
	}
	updateLocal <- true
}

// RemoveRemoteOrdersAt removes all orders at the given floor from the remote
// queue.
func RemoveRemoteOrdersAt(floor int) {
	var b def.ButtonType
	for b = def.BUTTON_CALL_UP; b < def.NUM_BUTTONS; b++ {
		remote.stopTimer(floor, b)
		remote.setOrder(floor, b, inactive)
	}
	updateLocal <- true
}

// RemoveOrdersAt removes all orders at the given floor in local and remote queue.
func RemoveOrdersAt(floor int, outgoingMsg chan<- network.Packet) {
	var b def.ButtonType
	for b = def.BUTTON_CALL_UP; b < def.NUM_BUTTONS; b++ {
		remote.stopTimer(floor, b)
		local.setOrder(floor, b, inactive)
		remote.setOrder(floor, b, inactive)
	}
	outgoingMsg <- network.Packet{Receiver_address: "129.241.187.155:"+strconv.Itoa(network.LocalListenPort), Sender_address: string(network.GetOwnID()), Data: []byte("Testmsg"), Length:7}
}

// ShouldStop returns whether the lift should stop when it reaches the given
// floor, going in the given direction.
func ShouldStop(floor int, dir def.MotorDirection) bool {
	return local.shouldStop(floor, dir)
}

// ChooseDirection returns the direction the lift should continue after the
// current floor, going in the given direction.
func ChooseDirection(floor int, dir def.MotorDirection) def.MotorDirection {
	return local.chooseDirection(floor, dir)
}

// IsLocalOrder returns whether there in an order with the given floor and
// button in the local queue.
func IsLocalOrder(floor int, button def.ButtonType) bool {
	return local.isOrder(floor, button)
}

// IsRemoteOrder returns true if there is a order with the given floor and
// button in the remote queue.
func IsRemoteOrder(floor int, button def.ButtonType) bool {
	return remote.isOrder(floor, button)
}

// ReassignOrders finds all orders assigned to a dead lift, removes them from
// the remote queue, and sends them on the network as new, unassigned orders.
func ReassignOrders(deadAddr string, outgoingMsg chan<- network.Packet) {
	//var f def.ButtonType
	for f := 0; f < def.NUM_FLOORS; f++ {
		var b def.ButtonType
		for b = 0; b < def.NUM_BUTTONS; b++ {
			if remote.matrix[f][b].addr == deadAddr {
				remote.setOrder(f, b, inactive)
				outgoingMsg <- network.Packet{Receiver_address: "129.241.187.155:"+strconv.Itoa(network.LocalListenPort), Sender_address: string(network.GetOwnID()), Data: []byte("Testmsg"), Length:7}
			}
		}
	}
}

// printQueues prints local and remote queue to screen in a somewhat legible
// manner.
func printQueues() {
	fmt.Printf(def.ColC)
	fmt.Println("Local   Remote")
	for f := def.NUM_FLOORS - 1; f >= 0; f-- {

		s1 := ""
		if local.isOrder(f, def.BUTTON_CALL_UP) {
			s1 += "↑"
		} else {
			s1 += " "
		}
		if local.isOrder(f, def.BUTTON_CALL_COMMAND) {
			s1 += "×"
		} else {
			s1 += " "
		}
		fmt.Printf(s1)
		if local.isOrder(f, def.BUTTON_CALL_DOWN) {
			fmt.Printf("↓   %d  ", f+1)
		} else {
			fmt.Printf("    %d  ", f+1)
		}

		s2 := "   "
		if remote.isOrder(f, def.BUTTON_CALL_UP) {
			fmt.Printf("↑")
			s2 += "(↑ " + remote.matrix[f][def.BUTTON_CALL_UP].addr[12:15] + ")"
		} else {
			fmt.Printf(" ")
		}
		if remote.isOrder(f, def.BUTTON_CALL_DOWN) {
			fmt.Printf("↓")
			s2 += "(↓ " + remote.matrix[f][def.BUTTON_CALL_DOWN].addr[12:15] + ")"
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%s", s2)
		fmt.Println()
	}
	fmt.Printf(def.ColN)
}

// updateLocalQueue checks remote queue for new orders assigned to this lift
// and copies them to the local queue.
func updateLocalQueue() {
	for {
		<-updateLocal
		for f := 0; f < def.NUM_FLOORS; f++ {
			var b def.ButtonType
			for b = 0; b < def.NUM_BUTTONS; b++ {
				if remote.isOrder(f, b) {
					if b != def.BUTTON_CALL_COMMAND && remote.matrix[f][b].addr == def.Laddr {
						if !local.isOrder(f, b) {
							local.setOrder(f, b, orderStatus{true, "", nil})
							newOrder <- true
						}
					}
				}
			}
		}
	}
}



func (q *queue) startTimer(floor int, button def.ButtonType) {
	const orderTimeout = 30 * time.Second

	q.matrix[floor][button].timer = time.NewTimer(orderTimeout)
	<-q.matrix[floor][button].timer.C
	OrderTimeoutChan <- def.OrderButton{Type: button, Floor: floor}
}

func (q *queue) stopTimer(floor int, button def.ButtonType) {
	if q.matrix[floor][button].timer != nil {
		q.matrix[floor][button].timer.Stop()
	}
}

func (q *queue) isEmpty() bool {
	for f := 0; f < def.NUM_FLOORS; f++ {
		for b := 0; b < def.NUM_BUTTONS; b++ {
			if q.matrix[f][b].active {
				return false
			}
		}
	}
	return true
}

func (q *queue) setOrder(floor int, button def.ButtonType, status orderStatus) {
	fmt.Println("Here!")
	// Ignore if order to be set is equal to order already in queue.
	if q.isOrder(floor, button) == status.active {
		return
	}


	q.matrix[floor][button] = status

	//def.SyncLightsChan <- true

	takeBackup <- true
	printQueues()

	fmt.Println("Here!")

}

func (q *queue) isOrder(floor int, button def.ButtonType) bool {
	return q.matrix[floor][button].active
}

func (q *queue) isOrdersAbove(floor int) bool {
	for f := floor + 1; f < def.NUM_FLOORS; f++ {
		var b def.ButtonType
		for b = 0; b < def.NUM_BUTTONS; b++ {
			if q.isOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		var b def.ButtonType
		for b = 0; b < def.NUM_BUTTONS; b++ {
			if q.isOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) chooseDirection(floor int, dir def.MotorDirection) def.MotorDirection {
	if q.isEmpty() {
		return def.STOP_Direction
	}
	switch dir {
	case def.DOWN_Direction:
		if q.isOrdersBelow(floor) && floor > 0 {
			return def.DOWN_Direction
		} else {
			return def.UP_Direction
		}
	case def.UP_Direction:
		if q.isOrdersAbove(floor) && floor < def.NUM_FLOORS-1 {
			return def.UP_Direction
		} else {
			return def.DOWN_Direction
		}
	case def.STOP_Direction:
		if q.isOrdersAbove(floor) {
			return def.UP_Direction
		} else if q.isOrdersBelow(floor) {
			return def.DOWN_Direction
		} else {
			return def.STOP_Direction
		}
	default:
	// 	def.CloseConnectionChan <- true
	// 	def.Restart.Run()
	// 	log.Printf("%sChooseDirection(): called with invalid direction %d, returning stop%s\n", def.ColR, dir, def.ColN)
	 	return 0
	}
}

func (q *queue) shouldStop(floor int, dir def.MotorDirection) bool {
	switch dir {
	case def.DOWN_Direction:
		return q.isOrder(floor, def.BUTTON_CALL_COMMAND) ||
			q.isOrder(floor, def.BUTTON_CALL_COMMAND) ||
			floor == 0 ||
			!q.isOrdersBelow(floor)
	case def.UP_Direction:
		return q.isOrder(floor, def.BUTTON_CALL_UP) ||
			q.isOrder(floor, def.BUTTON_CALL_COMMAND) ||
			floor == def.NUM_FLOORS-1 ||
			!q.isOrdersAbove(floor)
	case def.STOP_Direction:
		return q.isOrder(floor, def.BUTTON_CALL_COMMAND) ||
			q.isOrder(floor, def.BUTTON_CALL_UP) ||
			q.isOrder(floor, def.BUTTON_CALL_COMMAND)
	default:
		def.CloseConnectionChan <- true
		//def.Restart.Run()
		log.Fatalln(def.ColR, "This direction doesn't exist", def.ColN)
	}
	return false
}

func (q *queue) deepCopy() *queue {
	queueCopy := new(queue)
	for f := 0; f < def.NUM_FLOORS; f++ {
		for b := 0; b < def.NUM_BUTTONS; b++ {
			queueCopy.matrix[f][b] = q.matrix[f][b]
		}
	}
	return queueCopy
}



// CalculateCost returns how much effort it is for this lift to carry out
// the given order. Each sheduled stop and each travel between adjacent
// floors on the way towards target will add cost 2. Cost 1 is added if the
// lift starts between floors.
func CalculateCost(targetFloor, targetButton, prevFloor, currFloor int, currDir def.MotorDirection) int {
	q := local.deepCopy()
	q.setOrder(targetFloor, def.BUTTON_CALL_COMMAND, orderStatus{true, "", nil})

	cost := 0
	floor := prevFloor
	dir := currDir

	if currFloor == -1 {
		// Between floors, add 1 cost.
		cost++
	} else if dir != def.STOP_Direction {
		// At floor, but moving, add 2 cost.
		cost += 2
	}
	floor, dir = incrementFloor(floor, dir)

	// Simulate how the lift will move, and accumulate cost until it 'reaches' target.
	// Break after 10 iterations to assure against a stuck loop.
	for n := 0; !(floor == targetFloor && q.shouldStop(floor, dir)) && n < 10; n++ {
		if q.shouldStop(floor, dir) {
			cost += 2
			q.setOrder(floor, def.BUTTON_CALL_UP, inactive)
			q.setOrder(floor, def.BUTTON_CALL_DOWN, inactive)
			q.setOrder(floor, def.BUTTON_CALL_COMMAND, inactive)
		}
		dir = q.chooseDirection(floor, dir)
		floor, dir = incrementFloor(floor, dir)
		cost += 2
	}
	return cost
}

func incrementFloor(floor int, dir def.MotorDirection) (int, def.MotorDirection) {
	switch dir {
	case def.DOWN_Direction:
		floor--
	case def.UP_Direction:
		floor++
	case def.STOP_Direction:
		// Don't increment.
	default:
		def.CloseConnectionChan <- true
		//def.Restart.Run()
		log.Fatalln(def.ColR, "incrementFloor(): invalid direction, not incremented", def.ColN)
	}

	if floor <= 0 && dir == def.DOWN_Direction {
		dir = def.UP_Direction
		floor = 0
	}
	if floor >= def.NUM_FLOORS-1 && dir == def.UP_Direction {
		dir = def.DOWN_Direction
		floor = def.NUM_FLOORS - 1
	}
	return floor, dir
}