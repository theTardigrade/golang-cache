# golang-cache

This package allows arbitrary values to be stored and accessed using string keys.

## Support

If you use or find value in this package, please consider donating at PayPal: [https://www.paypal.me/jismithpp](https://www.paypal.me/jismithpp)

## Example

```golang
package main

import (
	"fmt"
	"time"

	cache "github.com/theTardigrade/golang-cache"
)

var (
	c = cache.NewCacheWithOptions(cache.Options{
		ExpiryDuration:         time.Hour * 24 * 7,
		CleanDuration:          time.Minute * 15,
		MaxValues:              1024,
		CleanMaxValuesPerSweep: 32,
		UnsetPreFunc: func(key string, value interface{}, setTime time.Time) {
			fmt.Printf("this function runs before the entry with the key \"%s\" is unset\n", key)
		},
		UnsetPostFunc: func(key string, value interface{}, setTime time.Time) {
			fmt.Printf("this function runs after the entry with the key \"%s\" is unset\n", key)
		},
	})
)

func main() {
	const key = "secret-number-pointer"
	x := 99

	c.Set(key, &x)

	p := c.MustGet(key).(*int)

	fmt.Println(p)  // prints pointer address
	fmt.Println(*p) // prints 99

	c.Iterate(func(key string, value interface{}, setTime time.Time) {
		fmt.Printf("%s --> %v (%s)\n", key, value, setTime.Format(time.RFC822))
	})

	c.Unset(key) // runs UnsetPreFunc, unsets the entry and runs UnsetPostFunc

	if found := c.Has(key); !found {
		fmt.Printf("entry with the key \"%s\" is no longer found\n", key) // this message should print
	}
}
```