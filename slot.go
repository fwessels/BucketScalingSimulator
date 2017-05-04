package BucketScalingSimulator

import (
	"sort"
	"fmt"
)

type Slot struct {
	id     int  // ID for this slot
	used   int  // Used space in this slot
	avail  int  // Available space in this slot
	buckets int // Number of buckets for this slot
	active bool // Indicate whether this bucket is already available (or offline until later)
}


var Slots []Slot

func init() {
	Slots = []Slot{}

	Slots = append(Slots, Slot{1, 0,100*1024*1024,0, true})
	Slots = append(Slots, Slot{2, 0,100*1024*1024,0, true})
	Slots = append(Slots, Slot{3, 0,100*1024*1024,0, true})
	Slots = append(Slots, Slot{4, 0,100*1024*1024,0, true})
}

func ExpandSlots(slotsToAdd, capacity int) {

	idStart := len(Slots)
	for i := 0; i < slotsToAdd; i++ {
		Slots = append(Slots, Slot{idStart+1+i, 0,capacity,0, true})
	}
}

func PrintSlots() {

	used, avail := 0, 0

	fmt.Println()
	fmt.Println("      Size", "      Used", "     Avail", "Use%")
	for _, s := range Slots {
		fmt.Printf("%10d %10d %10d %3.0f%%\n", s.used+s.avail, s.used, s.avail, float64(s.used*100)/float64(s.used+s.avail))
		used += s.used
		avail += s.avail
	}
	fmt.Println("----------", "----------", "----------", "----")
	fmt.Printf("%10d %10d %10d %3.0f%%\n", used+avail, used, avail, float64(used*100)/float64(used+avail))
	fmt.Println()
}

func getTotalSlots() int {

	return len(Slots)
}

func findId(id int, ids []int) bool {
	for _, i := range ids {
		if id == i {
			return true
		}
	}
	return false
}

// getLeastUsedSlot returns the slot that
// is used the least (has most free space)
func getLeastUsedSlot(existingIds []int) *Slot {

	// make copy for various sort operations below
	slots := make([]Slot, 0, len(Slots)-len(existingIds))
	for _, s := range Slots {
		if findId(s.id, existingIds) {
			continue // skip entries for already existing ids
		}
		slots = append(slots, s)
	}

	// sort slots on available size
	sort.Sort(bySlotMostAvail(slots))

	maxAvail := slots[0].avail

	islot := 0
	var slot Slot
	for islot, slot = range slots {
		if slot.avail < maxAvail*99/100 {
			break
		}
	}
	// keep slots within 99% of availability
	slots = slots[:islot+1]

	// sort on number of buckets
	sort.Sort(bySlotLeastBuckets(slots))

	// find and return slot corresponding to first entry in sorted aray
	for i, s := range Slots {
		if s.id == slots[0].id {
			return &Slots[i]
		}
	}
	return &Slots[0]
}

// bySlotMostAvail is a collection satisfying sort.Interface.
type bySlotMostAvail []Slot

func (s bySlotMostAvail) Len() int           { return len(s) }
func (s bySlotMostAvail) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s bySlotMostAvail) Less(i, j int) bool { return s[i].avail > s[j].avail }

// bySlotLeastBuckets is a collection satisfying sort.Interface.
type bySlotLeastBuckets []Slot

func (s bySlotLeastBuckets) Len() int           { return len(s) }
func (s bySlotLeastBuckets) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s bySlotLeastBuckets) Less(i, j int) bool { return s[i].buckets < s[j].buckets }
