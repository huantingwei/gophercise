package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// https://bicyclecards.com/how-to-play/blackjack/
// player: 1 face up to everyone, 1 face up to themselves
// dealer: 1 face up to everyone, 1 face down
//  "stand" or "hit"

// When the dealer has served every player, the dealers face-down card is turned up.
// If the total is 17 or more, it must stand.
// If the total is 16 or under, they must take a card.
// The dealer must continue to take cards until the total is 17 or more, at which point the dealer must stand.
// If the dealer has an ace, and counting it as 11 would bring the total to 17 or more (but not over 21),
//  the dealer must count the ace as 11 and stand

const MAX = 21
const NUM_DECK = 3

type Deck struct {
	cards []string
}

type Player struct {
	role     string
	cards    []string
	total    []int
	isBusted bool
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Deal() string {
	ret := d.cards[len(d.cards)-1]
	d.cards = append(d.cards[:len(d.cards)-1], d.cards[len(d.cards):]...)
	return ret
}

func (d *Deck) Show() {
	fmt.Printf("| ")
	for _, card := range d.cards {
		fmt.Printf("%3v", card)
	}
	fmt.Printf(" |\n")
}

// only takes in JQK and 2~10
func getPoint(card string) int {
	if card == "J" || card == "Q" || card == "K" {
		return 10
	} else {
		n, err := strconv.Atoi(card)
		if err != nil {
			panic(err)
		} else {
			return n
		}
	}
}

func (p *Player) check(threshold int) bool {
	m := 0
	for i, pos := range p.total {
		if i == 0 || pos < m {
			m = pos
		}
	}
	return m > threshold
}

func (p *Player) add(card string) {
	p.cards = append(p.cards, card)

	if card == "A" {
		if len(p.total) == 0 {
			p.total = append(p.total, 1)
			p.total = append(p.total, 11)
		} else {
			for i := 0; i < len(p.total); i++ {
				p.total[i] += 1
				p.total = append(p.total, p.total[i]+11)
			}
		}
	} else {
		if len(p.total) == 0 {
			p.total = append(p.total, getPoint(card))
		} else {
			for i := 0; i < len(p.total); i++ {
				p.total[i] += getPoint(card)
			}
		}
	}
}

func (p *Player) Init(name string) {
	p.total = make([]int, 0)
	p.cards = make([]string, 0)
	p.role = name
	p.isBusted = false
}

func (p *Player) SetName(name string) {
	p.role = name
}

func (p *Player) Hit(deck *Deck) {
	p.add(deck.Deal())
}

func (p *Player) CheckIsOver(threshold int) bool {
	isOver := p.check(threshold)
	if threshold == MAX && isOver {
		p.isBusted = true
	}
	return isOver
}

func (p *Player) Show() {
	fmt.Printf("\n%s has:\n", p.role)
	for i, card := range p.cards {
		fmt.Printf("| %2d: %2v ", i, card)
		fmt.Printf("|\n")
	}
	fmt.Printf("Total = ")
	for i, t := range p.total {
		if i != 0 {
			fmt.Printf("or %2d ", t)
		} else {
			fmt.Printf("%2d ", t)
		}
	}
	fmt.Printf("\n\n==============================\n")
	fmt.Printf("%10s's best score = %2d\n", p.role, p.BestScore())
	fmt.Printf("==============================\n")
}

func (p *Player) BestScore() int {
	if len(p.total) == 0 {
		return 0
	}

	best := p.total[0]
	for i := 1; i < len(p.total); i++ {
		if p.total[i] <= MAX && p.total[i] > best {
			best = p.total[i]
		}
	}
	return best
}

func main() {
	cards := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	rand.Seed(time.Now().UnixNano())
	var deck Deck
	for i := 0; i < 4*NUM_DECK; i++ {
		deck.cards = append(deck.cards, cards...)
	}

	var player, dealer Player
	var action string
	player.Init("Player 1")
	dealer.Init("Dealer")

	fmt.Printf("Shuffling...\n")
	deck.Shuffle()
	deck.Show()

	fmt.Printf("Dealing the first round (face down)\n")
	player.Hit(&deck)
	dealer.Hit(&deck)

	fmt.Printf("Dealing the second round (face up)\n")
	player.Hit(&deck)
	dealer.Hit(&deck)
	player.Show()
	dealer.Show()

	for {
		if !dealer.CheckIsOver(16) {
			dealer.Hit(&deck)
		} else {
			break
		}
	}
	dealer.Show()

	fmt.Printf("Player's turn!!!\n")
	for {
		if player.CheckIsOver(MAX) {
			fmt.Printf("You are busted!\n")
			break
		}

		player.Show()
		fmt.Printf("Hit or Stand? (hit or stand)\n")
		fmt.Scanf("%s", &action)

		if action == "hit" || action == "h" || action == "H" {
			player.Hit(&deck)
		} else if action == "stand" || action == "s" || action == "S" {
			break
		} else {
			fmt.Printf("Please enter `hit` or `stand`\n")
		}

	}
	fmt.Printf("Show hands!\n")
	player.Show()
	dealer.Show()
}
