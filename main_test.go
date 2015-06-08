package main

import "testing"

func TestExampleDecode(t *testing.T) {

	//process("/Users/dave/Dropbox (Personal)/Pictures/2005/2005-04-23 - Had Kandi @ KoKo/John/IMG_1738.JPG")

	updateDateIfWrong("/Users/dave/Dropbox (Personal)/Pictures/2007/2007-08-28 - Wood floor/CIMG0032.JPG")

	/*
		fname := "sample1.jpg"

		f, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}

		// Optionally register camera makenote data parsing - currently Nikon and
		// Canon are supported.
		exif.RegisterParsers(mknote.All...)

		x, err := exif.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		camModel, _ := x.Get(exif.Model) // normally, don't ignore errors!
		fmt.Println(camModel.StringVal())

		focal, _ := x.Get(exif.FocalLength)
		numer, denom, _ := focal.Rat2(0) // retrieve first (only) rat. value
		fmt.Printf("%v/%v", numer, denom)

		// Two convenience functions exist for date/time taken and GPS coords:
		tm, _ := x.DateTime()
		fmt.Println("Taken: ", tm)

		lat, long, _ := x.LatLong()
		fmt.Println("lat, long: ", lat, ", ", long)
	*/
}
