package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

// Die is a trivial type to represent a single die with a values between
// min and max, inclusive.
type Die struct {
	min  *big.Int
	max  *big.Int
	rmax *big.Int
}

// NewDieFromString takes a string representation of the number of faces,
// a die should have, and returns a Die with values 1..N inclusive. An
// error is returned if N cannot be converted into an integer.
func NewDieFromString(N string) (*Die, error) {
	n, err := strconv.Atoi(N)
	if err != nil {
		return nil, err
	}
	return NewDieFromMinMax(int64(1), int64(n)), nil
}

// NewDieFromMinMax creates and returns a Die from the specified
// min and max.
func NewDieFromMinMax(min, max int64) *Die {
	return &Die{
		min:  big.NewInt(min),
		max:  big.NewInt(max),
		rmax: big.NewInt(max - min + 1),
	}
}

// Roll returns a random integer between min and max inclusively.
func (d *Die) Roll() *big.Int {
	if n, err := rand.Int(rand.Reader, d.rmax); err != nil {
		panic(err)
	} else {
		return n.Add(n, d.min)
	}
}

// RollNdF creates an F-sided die, and rolls it N times,
// returning the resulting rolls as a comma-deimited string.
// If total is true, the sum of the rolls is appended to the
// result.
func RollNdF(n, f string, total bool) (string, int64) {
	die, err := NewDieFromString(f)
	if err != nil {
		// f is not an integer
		panic(err)
	}

	c, err := strconv.Atoi(n)
	if err != nil {
		// n is not an integer
		panic(err)
	}

	var (
		s   string
		sum = big.NewInt(int64(0))
	)

	// roll the die n times
	for i := 0; i < c; i++ {
		v := die.Roll()
		if i > 0 {
			s = s + ","
		}
		sum = sum.Add(sum, v)
		s = s + v.String()
	}
	if total {
		return fmt.Sprintf("%s = %s", s, sum.String()), sum.Int64()
	}
	return s, sum.Int64()
}
