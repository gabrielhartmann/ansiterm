package winterm

// AddInRange increments a value by the passed quantity while ensuring the values
// always remain within the supplied min / max range.
func AddInRange(n SHORT, increment SHORT, min SHORT, max SHORT) SHORT {
	return ensureInRange(n+increment, min, max)
}

// ensureInRange adjusts the passed value, if necessary, to ensure it is within
// the passed min / max range.
func ensureInRange(n SHORT, min SHORT, max SHORT) SHORT {
	if n < min {
		return min
	} else if n > max {
		return max
	} else {
		return n
	}
}
