# Snow Consensus Simulator
* Simulate 200 nodes in a peer-to-peer network using localhost <br />
* Each node starts with a list of empty public transaction list and a non-empty private transaction list <br />
* Using the snow consensus algorithm, the network gradually comes to a consensus on the public transaction list <br />
* Each transaction is modelled with an ID <br />
* To run <br />
    1. cd ./apps
    2. go build -o  ./../bin/nodeApp to build the node application executable
    3. run test Test_200Nodes_200Processes under test/snowball_test.go (in debug mode)