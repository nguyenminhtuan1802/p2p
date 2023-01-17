package network

import (
	"fmt"
	"strconv"
	"tnguyen/blockchainexample/transaction"
	"time"
)

// A node represents a peer in P2P network
// Node contains a server that listens to peer node request 
// Clients in a node are used to send request to this node's peers
// Note: Channel provides data to respond to peers' requests. In this
//       example the data is transaction id, but channel can support other data types
type Node struct {
	Protocol 				string
	Id       				int
	Server   				*Server
	Clients  				map[int]*Client
	PublicTransactionList	*transaction.TransactionList
	PrivateTransactionList	*transaction.TransactionList
	Channel					chan int
	Validator				ConsensusValidator
	FinishRunning			bool
}

func NodeConstructor(protocol string, id int,
	 publicTrans *transaction.TransactionList,
	 privateTrans *transaction.TransactionList,) *Node {
	node := new(Node)
	node.Protocol = protocol
	node.Id = id
	node.Clients = make(map[int]*Client)
	node.Channel = make(chan int, 1)

	node.PrivateTransactionList = privateTrans
	node.PublicTransactionList = publicTrans
	node.FinishRunning = false

	return node
}

func (n *Node) Start() {	
	n.Server = ServerConstructor(n.Protocol, addresses[n.Id], n)
	go n.Server.Start()
}

func (n *Node) Run() {
	time.Sleep(5 * time.Second)
	for {
		if n.Validator.Update() == false {
			fmt.Println("[NODE ", n.Id, "] DONE")
			for _,v := range n.PublicTransactionList.Transactions {
				fmt.Println("[NODE ", n.Id, "] :", v.Data)
			}
			n.FinishRunning = true
			break
		}
	}
}

// Connect to another peer
func (n *Node) Connect(another *Node) {
	_, ok := n.Clients[another.Id]

	if ok {
		return
	}

	client := ClientConstructor(another.Protocol, addresses[another.Id])
	client.Connect()
	anotherClient := ClientConstructor(n.Protocol, addresses[n.Id])
	anotherClient.Connect()

	n.Clients[another.Id] = client
	another.Clients[n.Id] = anotherClient
}

func (n *Node) PingAllConnections() {
	var i = 0
	for _, v := range n.Clients {
		v.Send("Ping from Node " + strconv.Itoa(n.Id) + " to Node " + v.Address)
		i++
		go v.Receive()
	}
	fmt.Println("[NODE ", n.Id, "] has", i, " connections")
}

func (n *Node) ShutDown() {
	n.Server.ShutDown()
}

func (n *Node) GetNextUndecidedTransaction() (int, *transaction.Transaction) {
	index := n.PrivateTransactionList.CurrentIdx
	tran := n.PrivateTransactionList.Transactions[index]
	return index, tran
}

func (n *Node) GetTransactionData(index int) int {
	publicTransaction := n.PublicTransactionList.GetTransactionData(index)
	if (publicTransaction != nil) {
		return publicTransaction.Data
	} else {
		privateTransaction := n.PrivateTransactionList.GetTransactionData(index)
		if (privateTransaction != nil) {
			return privateTransaction.Data
		} else {
			panic("Invalid index")
		}
	}
}

type ConsensusValidator interface {
	Update() bool
}

var addresses = map[int]string{
	1:  "localhost:60000",
	2:  "localhost:60001",
	3:  "localhost:60002",
	4:  "localhost:60003",
	5:  "localhost:60004",
	6:  "localhost:60005",
	7:  "localhost:60006",
	8:  "localhost:60007",
	9:  "localhost:60008",
	10: "localhost:60090",
	11: "localhost:60010",
	12: "localhost:60011",
	13: "localhost:60012",
	14: "localhost:60013",
	15: "localhost:60014",
	16: "localhost:60015",
	17: "localhost:60016",
	18: "localhost:60017",
	19: "localhost:60018",
	20: "localhost:60019",
	21:  "localhost:60020",
	22:  "localhost:60021",
	23:  "localhost:60022",
	24:  "localhost:60023",
	25:  "localhost:60024",
	26:  "localhost:60025",
	27:  "localhost:60026",
	28:  "localhost:60027",
	29:  "localhost:60028",
	30: "localhost:60029",
	31: "localhost:60030",
	32: "localhost:60031",
	33: "localhost:60032",
	34: "localhost:60033",
	35: "localhost:60034",
	36: "localhost:60035",
	37: "localhost:60036",
	38: "localhost:60037",
	39: "localhost:60038",
	40: "localhost:60039",
	41: "localhost:60040",
	42: "localhost:60041",
	43: "localhost:60042",
	44: "localhost:60043",
	45: "localhost:60044",
	46: "localhost:60045",
	47: "localhost:60046",
	48: "localhost:60047",
	49: "localhost:60048",
	50: "localhost:60049",
	51: "localhost:60050",
	52: "localhost:60051",
	53: "localhost:60052",
	54: "localhost:60053",
	55: "localhost:60054",
	56: "localhost:60055",
	57: "localhost:60056",
	58: "localhost:60057",
	59: "localhost:60058",
	60: "localhost:60059",
	61: "localhost:60060",
	62: "localhost:60061",
	63: "localhost:60062",
	64: "localhost:60063",
	65: "localhost:60064",
	66: "localhost:60065",
	67: "localhost:60066",
	68: "localhost:60067",
	69: "localhost:60068",
	70: "localhost:60069",
	71: "localhost:60070",
	72: "localhost:60071",
	73: "localhost:60072",
	74: "localhost:60073",
	75: "localhost:60074",
	76: "localhost:60075",
	77: "localhost:60076",
	78: "localhost:60077",
	79: "localhost:60078",
	80: "localhost:60079",
	81: "localhost:60080",
	82: "localhost:60081",
	83: "localhost:60082",
	84: "localhost:60083",
	85: "localhost:60084",
	86: "localhost:60085",
	87: "localhost:60086",
	88: "localhost:60087",
	89: "localhost:60088",
	90: "localhost:60089",
	91: "localhost:61000",
	92: "localhost:60091",
	93: "localhost:60092",
	94: "localhost:60093",
	95: "localhost:60094",
	96: "localhost:60095",
	97: "localhost:60096",
	98: "localhost:60097",
	99: "localhost:60098",
	100: "localhost:60099",
	101: "localhost:60100",
	102: "localhost:60101",
	103: "localhost:60102",
	104: "localhost:60103",
	105: "localhost:60104",
	106: "localhost:60105",
	107: "localhost:60106",
	108: "localhost:60107",
	109: "localhost:60108",
	110: "localhost:60109",
	111: "localhost:60110",
	112: "localhost:60111",
	113: "localhost:60112",
	114: "localhost:60113",
	115: "localhost:60114",
	116: "localhost:60115",
	117: "localhost:60116",
	118: "localhost:60117",
	119: "localhost:60118",
	120: "localhost:60119",
	121: "localhost:60120",
	122: "localhost:60121",
	123: "localhost:60122",
	124: "localhost:60123",
	125: "localhost:60124",
	126: "localhost:60125",
	127: "localhost:60126",
	128: "localhost:60127",
	129: "localhost:60128",
	130: "localhost:60129",
	131: "localhost:60130",
	132: "localhost:60131",
	133: "localhost:60132",
	134: "localhost:60133",
	135: "localhost:60134",
	136: "localhost:60135",
	137: "localhost:60136",
	138: "localhost:60137",
	139: "localhost:60138",
	140: "localhost:60139",
	141: "localhost:60140",
	142: "localhost:60141",
	143: "localhost:60142",
	144: "localhost:60143",
	145: "localhost:60144",
	146: "localhost:60145",
	147: "localhost:60146",
	148: "localhost:60147",
	149: "localhost:60148",
	150: "localhost:60149",
	151: "localhost:60150",
	152: "localhost:60151",
	153: "localhost:60152",
	154: "localhost:60153",
	155: "localhost:60154",
	156: "localhost:60155",
	157: "localhost:60156",
	158: "localhost:60157",
	159: "localhost:60158",
	160: "localhost:60159",
	161: "localhost:60160",
	162: "localhost:60161",
	163: "localhost:60162",
	164: "localhost:60163",
	165: "localhost:60164",
	166: "localhost:60165",
	167: "localhost:60166",
	168: "localhost:60167",
	169: "localhost:60168",
	170: "localhost:60169",
	171: "localhost:60170",
	172: "localhost:60171",
	173: "localhost:60172",
	174: "localhost:60173",
	175: "localhost:60174",
	176: "localhost:60175",
	177: "localhost:60176",
	178: "localhost:60177",
	179: "localhost:60178",
	180: "localhost:60179",
	181: "localhost:60180",
	182: "localhost:60181",
	183: "localhost:60182",
	184: "localhost:60183",
	185: "localhost:60184",
	186: "localhost:60185",
	187: "localhost:60186",
	188: "localhost:60187",
	189: "localhost:60188",
	190: "localhost:60189",
	191: "localhost:60190",
	192: "localhost:60191",
	193: "localhost:60192",
	194: "localhost:60193",
	195: "localhost:60194",
	196: "localhost:60195",
	197: "localhost:60196",
	198: "localhost:60197",
	199: "localhost:60198",
	200: "localhost:60199",
}