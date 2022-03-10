package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// https://bicyclecards.com/how-to-play/blackjack/
// Deck
// shuffle
// deal - player
// deal - dealer
// Player
// .cards
// stand()
// hit()
// isNatural
// isBusted

// six-deck game (312 cards)
// A = 11 or 1
// JQK = 10

// every player place bets
// player: 1 face up to everyone, 1 face up to themselves
// dealer: 1 face up to everyone, 1 face down

// first 2 cards = "natural" -> blackjack
// => dealer pays 1.5x of the player's bet
// => other players pay 1x of the player's bet

//  "stand" or "hit"

// When the dealer has served every player, the dealers face-down card is turned up.
// If the total is 17 or more, it must stand.
// If the total is 16 or under, they must take a card.
// The dealer must continue to take cards until the total is 17 or more, at which point the dealer must stand.
// If the dealer has an ace, and counting it as 11 would bring the total to 17 or more (but not over 21), the dealer must count the ace as 11 and stand

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
	// go through every card again
	// better version: just check the newest card
	var possibilities []int
	for _, card := range p.cards {
		if card == "A" {
			// for each possibilities, duplicate
			if len(possibilities) == 0 {
				possibilities = append(possibilities, 1)
				possibilities = append(possibilities, 11)
			} else {
				for _, pos := range possibilities {
					pos += 1
					possibilities = append(possibilities, pos+11)
				}
			}
		} else {
			if len(possibilities) == 0 {
				possibilities = append(possibilities, getPoint(card))
			} else {
				for i := 0; i < len(possibilities); i++ {
					possibilities[i] += getPoint(card)
				}
			}
		}
	}

	// set possible totals
	p.total = possibilities
	// check if over the threshold
	m := 0
	for i, pos := range possibilities {
		if i == 0 || pos < m {
			m = pos
		}
	}
	return m > threshold
}

func (p *Player) add(card string) {
	p.cards = append(p.cards, card)
}

func (p *Player) Hit(deck *Deck) {
	p.add(deck.Deal())
}

func (p *Player) CheckIsOver(threshold int) bool {
	return p.check(threshold)
}

func (p *Player) Show() {
	fmt.Printf("%s has:\n", p.role)
	for i, card := range p.cards {
		fmt.Printf("| %2d: %2v ", i, card)
		fmt.Printf("|\n")
	}
	fmt.Println()
}

func (p *Player) SetName(name string) {
	p.role = name
}

func main() {
	cards := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	rand.Seed(time.Now().UnixNano())
	var deck Deck
	for i := 0; i < NUM_DECK; i++ {
		deck.cards = append(deck.cards, cards...)
	}

	var player, dealer Player
	var action string
	player.SetName("Player 1")
	dealer.SetName("Dealer")

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
			dealer.Show()
		} else {
			break
		}
	}

	fmt.Printf("Player's turn!!!\n")
	for {
		fmt.Printf("Hit or Stand? (hit or stand))\n")
		fmt.Scanf("%s", &action)
		fmt.Printf("%s", action)
		if action == "hit" {
			player.Hit(&deck)
			player.Show()

		} else if action == "stand" {
			break
		} else {
			fmt.Printf("Please enter hit or stand\n")
		}
		if player.CheckIsOver(MAX) {
			fmt.Printf("You are busted!\n")
			break
		}
	}
	fmt.Printf("Show hands!\n")
	player.Show()
	dealer.Show()

}
