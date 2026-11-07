package main

import "io"

// usable by everything for which source fulfills the io.Reader interface and
// dest fulfills the io.Writer interface
// Rule of Thumb: interfaces in, concrete types out
// Be liberal in what you accept (interfaces), and be specific in what you return (concrete types).
func copySourceToDest(source io.Reader, dest io.Writer) error {
	b, err := io.ReadAll(source)
	if err != nil {
		return err
	}
	_, err = dest.Write(b)
	return err
}
