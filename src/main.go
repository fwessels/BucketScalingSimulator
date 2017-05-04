package main

import (
	sim "github.com/fwessels/BucketScalingSimulator"
)

func main() {
	sim.CreateBucket("CT")
	sim.CreateBucket("MR")
	sim.CreateBucket("CR")
	sim.CreateBucket("XA")
	sim.CreateBucket("US")

	for i := 0; i < 200; i++ {
		sim.CreateObject("CT", 2*512*512)
	}

	for i := 0; i < 10; i++ {
		sim.CreateObject("CR", 2*2096*2096)
	}

	for i := 0; i < 100; i++ {
		sim.CreateObject("CT", 2*512*512)
	}

	for i := 0; i < 200; i++ {
		sim.CreateObject("MR", 2*256*256)
	}

	for i := 0; i < 200; i++ {
		sim.CreateObject("CT", 2*512*512)
	}

	sim.PrintSlots()

	sim.ExpandSlots(4, 100*1024*1024)

	sim.PrintSlots()

	sim.CreateBucket("DA")

	for i := 0; i < 10; i++ {
		sim.CreateObject("DA", 2*1024*1024)
	}

	for i := 0; i < 10; i++ {
		sim.CreateObject("CT", 2*512*512)
	}

	for i := 0; i < 10; i++ {
		sim.CreateObject("MR", 2*256*256)
	}

	for i := 0; i < 2; i++ {
		sim.CreateObject("CR", 2*2096*2096)
	}

	sim.CreateBucket("PT")
	sim.CreateBucket("NM")
	sim.CreateBucket("CT2")

	for i := 0; i < 5128+7800; i++ {
		sim.CreateObject("PT", 2*128*128)
	}

	sim.PrintSlots()

	sim.ExpandSlots(4, 120*1024*1024)

	sim.PrintSlots()

	for i := 0; i < 700; i++ {
		sim.CreateObject("CT2", 2*512*512)
	}

	sim.PrintSlots()

}
