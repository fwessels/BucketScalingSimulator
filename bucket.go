package BucketScalingSimulator

import (
	"fmt"
	"math/rand"
)

type BucketInfo struct {
	slots []*Slot
}

var Buckets map[string]*BucketInfo

func init() {
	Buckets = make(map[string]*BucketInfo)
}

const bucketGrowWatermark = 0.75

func (bi BucketInfo) needToGrowBucket() bool {

	if len(bi.slots) == getTotalSlots() {
		return false // Cannot grow any further
	}

	bucketUsed := 0
	bucketAvail := 0

	for _, slot := range bi.slots {
		bucketUsed += slot.used
		bucketAvail += slot.avail
	}

	filledPct := float64(bucketUsed) / float64(bucketUsed + bucketAvail)
	//fmt.Println("filled", filledPct)

	if filledPct >= bucketGrowWatermark {
		return true
	}

	return false
}

// addSlotToBucket
func (bi *BucketInfo) addSlotToBucket() []int {

	ids := []int{}
	// Check if we are not already at the max
	if len(bi.slots) != getTotalSlots() {
		for _, slot := range bi.slots {
			ids = append(ids, slot.id)
		}

		slot, ok := getLeastUsedSlot(ids)
		if ok {
			slot.buckets += 1

			bi.slots = append(bi.slots, slot)
		}
	}

	ids = ids[:0]
	for _, slot := range bi.slots {
		ids = append(ids, slot.id)
	}
	return ids
}

// chooseSlot
func (bi BucketInfo) chooseSlot() *Slot {

	totalAvail := 0

	for _, slot := range bi.slots {
		totalAvail += slot.avail
	}

	if totalAvail == 0 {
		return bi.slots[0] // We are out of space, just return 1st slot
	}

	// Pick a random number somewhere in the total available space
	pick := rand.Intn(totalAvail)

	avail := 0

	// Return the slot corresponding
	for _, slot := range bi.slots {
		avail += slot.avail
		if pick < avail {
			return slot
		}
	}

	// Otherwise return the last
	return bi.slots[len(bi.slots)-1]
}

// CreateBucket
func CreateBucket(name string) error {

	if _, ok := Buckets[name]; ok {
		return nil // already exists, return
	}

	slot, ok := getLeastUsedSlot([]int{})
	if !ok {
		return errors.New("No space available to create bucket")
	}

	slot.buckets += 1
	fmt.Println(name, "create:", "slot ", []int{slot.id})

	bucketInfo := BucketInfo{[]*Slot{}}
	bucketInfo.slots = append(bucketInfo.slots, slot)
	Buckets[name] = &bucketInfo
	return nil
}

// CreateObject
func CreateObject(bucket string, size int) {

	bucketInfo, ok := Buckets[bucket]
	if ok {
		if bucketInfo.needToGrowBucket() {
			ids := bucketInfo.addSlotToBucket()
			fmt.Println(bucket, "  grow:", "slots", ids)
		}
		slot := bucketInfo.chooseSlot(/*size*/)
		if slot.avail - size < 0 {
			// fmt.Println("Out of available space")
			panic("Out of available space")
		}
		slot.used += size
		slot.avail -= size
	}
}