package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	Patterns [][]int
	P0       = []int{0, 6, 12, 18, 24}
	P1       = []int{4, 8, 12, 16, 20}
	P2       = []int{0, 4, 12, 20, 24}
	P3L1     = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	P4L2     = []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	P5L3     = []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	P6L4     = []int{15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	P7L5     = []int{0, 1, 2, 3, 4, 10, 11, 12, 13, 14}
	P8L6     = []int{5, 6, 7, 8, 9, 15, 16, 17, 18, 19}
	P9L7     = []int{10, 11, 12, 13, 14, 20, 21, 22, 23, 24}
	P10L8    = []int{0, 1, 2, 3, 4, 15, 16, 17, 18, 19}
	P11L9    = []int{5, 6, 7, 8, 9, 20, 21, 22, 23, 24}
	P12L10   = []int{0, 1, 2, 3, 4, 20, 21, 22, 23, 24}
	P13E     = []int{0, 1, 2, 3, 4, 7, 12, 17, 22}
	P14FF    = []int{0, 4, 6, 8, 12, 16, 18, 20, 24}
	//TL=[]int{}
	P15H = []int{0, 1, 2, 3, 4, 5, 9, 10, 12, 14, 15, 19, 20, 21, 22, 23, 24}
)
var m = make(map[int]int, 0)
var cards = make([]*Card, 0)
var a []int
var lock sync.Mutex

type Card struct {
	Id      string
	Numbers []int
	Index   []int
	Prize   []int
}

func TestB(t *testing.T) {
	Patterns = append(Patterns, P0)
	Patterns = append(Patterns, P1)
	Patterns = append(Patterns, P2)
	Patterns = append(Patterns, P3L1)
	Patterns = append(Patterns, P4L2)
	Patterns = append(Patterns, P5L3)
	Patterns = append(Patterns, P6L4)
	Patterns = append(Patterns, P7L5)
	Patterns = append(Patterns, P8L6)
	Patterns = append(Patterns, P9L7)
	Patterns = append(Patterns, P10L8)
	Patterns = append(Patterns, P11L9)
	Patterns = append(Patterns, P12L10)
	Patterns = append(Patterns, P13E)
	Patterns = append(Patterns, P14FF)
	Patterns = append(Patterns, P15H)
	GenCard(1000000)
	fmt.Println("GenCard over")
	PrizeNumber()
	fmt.Println("PrizeNumber over")
	p := 10
	for i, i2 := range a {
		fmt.Printf("prize %d-%d \n", i, i2)
		now := time.Now()
		//for _, card := range cards {
		//	card.GenIndex(i2)
		//	card.Settle()
		//}

		group := sync.WaitGroup{}
		l := len(cards)
		count := l / p
		for j := 0; j < p; j++ {
			group.Add(1)
			i1 := j * count
			i3 := count * (j + 1)
			ints := cards[i1:i3]
			go func(in []*Card, n int) {
				for _, i4 := range in {
					i4.GenIndex(n)
					i4.Settle()
				}
				group.Done()
			}(ints, i2)
		}
		if l%count != 0 {
			ints := cards[count*p:]
			for _, card := range ints {
				card.GenIndex(i2)
				card.Settle()
			}
		}
		group.Wait()
		fmt.Println("settle time==>", time.Since(now))
		fmt.Println(m)
		time.Sleep(3 * time.Second)
	}
}

func PrizeNumber() {
	for i := 1; i <= 75; i++ {
		a = append(a, i)
	}
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	a = append(a[:45], a[75:]...)
	fmt.Println("settle=>", a)
}

func GenCard(b int) {
	for i := 0; i < b; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		var card []int
		for i := 1; i <= 5; i++ {
			a := A(card, i, r)
			card = append(card, a)
			a1 := A(card, i, r)
			card = append(card, a1)
			a2 := A(card, i, r)
			card = append(card, a2)
			a3 := A(card, i, r)
			card = append(card, a3)
			a4 := A(card, i, r)
			card = append(card, a4)
			//	fmt.Printf("%d-%d-%d-%d-%d\n", a, a1, a2, a3, a4)
		}
		cardd := Card{
			Id:      fmt.Sprintf("c-%d", i),
			Numbers: card,
		}
		cards = append(cards, &cardd)
	}
}

func (c *Card) GenIndex(i int) {
	i2 := i / 15
	if i%15 > 0 {
		i2 = i2 + 1
	}
	i2 = i2 - 1
	//0-4
	//5-9
	//10-14
	//15-19
	//20-24
	i3 := i2 * 5
	i4 := (i2+1)*5 - 1
	for j := i3; j <= i4; j++ {
		if i == c.Numbers[j] {
			//fmt.Printf("%d<=>%d\t", i, j)
			c.Index = append(c.Index, j)
		}
	}
}
func (c *Card) Settle() {
	for i, pattern := range Patterns {
		if c.S(i) {
			continue
		}
		if F(pattern, c.Index) {
			Add(i)
			c.Prize = append(c.Prize, i)
			//fmt.Println("P", i, " success")
			//for i2, number := range c.Numbers {
			//	i22 := i2 + 1
			//	if i22%5 == 0 {
			//		fmt.Printf("%d\n", number)
			//	} else {
			//		fmt.Printf("%d-", number)
			//	}
			//}
			//fmt.Println("index=", c.Index)
		}
	}
}
func (c *Card) S(i int) bool {
	if c.Prize != nil {
		for _, i3 := range c.Prize {
			if i == i3 {
				return true
			}
		}
	}
	return false
}

func F(p, c []int) bool {
	var q int
	for _, i := range p {
		if 12 == i {
			q++
			continue
		}
		for _, i2 := range c {
			if i == i2 {
				q++
				continue
			}
		}
	}
	return len(p) <= q
}

func Add(i int) {
	lock.Lock()
	m[i]++
	defer lock.Unlock()
}
func A(s []int, d int, r *rand.Rand) int {
	intn := (r.Intn(15) + 1) + (d-1)*15
	if Con(s, intn) {
		return A(s, d, r)
	} else {
		return intn
	}
}

func Con(s []int, d int) bool {
	for _, i := range s {
		if i == d {
			return true
		}
	}
	return false
}
