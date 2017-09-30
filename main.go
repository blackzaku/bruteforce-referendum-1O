package main

import (
	"fmt"
	"o1st"
	"strconv"
	"sync"
	"flag"
	"strings"
	"errors"
)

const(
	MAX_DNI int = 1e5 - 1
)

var (
	MIN_YEAR int = 1965
	MAX_YEAR int = 2000

	DNI_NUM_MIN_RANGE = 0
	DNI_NUM_MAX_RANGE = MAX_DNI
	DNI_LETTER_MIN_RANGE = uint8('A')
	DNI_LETTER_MAX_RANGE = uint8('Z')

	DB_FOLDER = "db"

	ZIP_CODES []string
)

var (
	dni string
	minDni string
	maxDni string
	zip string
	province string
	year int
)

func init() {
	flag.StringVar(&DB_FOLDER, "db", "db", "Sets the folder from which the database is loaded")
	flag.StringVar(&dni, "dni", "", "Sets exact DNI")
	flag.IntVar(&DNI_NUM_MIN_RANGE, "min-dni-num", 0, "Sets minimum DNI number")
	flag.IntVar(&DNI_NUM_MAX_RANGE, "max-dni-num", 99999, "Sets maximum DNI number")
	flag.StringVar(&minDni, "min-dni-letter", "A", "Sets minimum DNI letter")
	flag.StringVar(&maxDni, "max-dni-letter", "Z", "Sets maximum DNI letter")
	flag.StringVar(&zip, "zip", "", "Restricts to a single zip code")
	flag.StringVar(&province, "province", "all", "Sets province for the zip codes")
	flag.IntVar(&MIN_YEAR, "min-year", 1965, "Sets minimum year")
	flag.IntVar(&MAX_YEAR, "max-year", 2000, "Sets maximum year")
	flag.IntVar(&year, "year", 0, "Sets exact year")
}

func processArguments() {
	if dni != "" {
		DNI_LETTER_MIN_RANGE = dni[len(dni) - 1]
		DNI_LETTER_MAX_RANGE = dni[len(dni) - 1]
		dni = dni[:len(dni) - 1]
		parsed,_ := strconv.ParseInt(dni, 10, 0)
		DNI_NUM_MIN_RANGE = int(parsed % 1e5)
		DNI_NUM_MAX_RANGE = DNI_NUM_MIN_RANGE
	} else {
		if DNI_NUM_MAX_RANGE >= 1e5 {
			if DNI_NUM_MAX_RANGE - DNI_NUM_MIN_RANGE >= 1e5 {
				// Full range
				DNI_NUM_MAX_RANGE = 0
				DNI_NUM_MAX_RANGE = MAX_DNI
			} else {
				DNI_NUM_MAX_RANGE = DNI_NUM_MAX_RANGE % 1e5
				DNI_NUM_MIN_RANGE = DNI_NUM_MIN_RANGE % 1e5
				if DNI_NUM_MAX_RANGE < DNI_NUM_MIN_RANGE {
					// Change order
					DNI_NUM_MAX_RANGE, DNI_NUM_MIN_RANGE = DNI_NUM_MIN_RANGE, DNI_NUM_MAX_RANGE
				}
			}

		}
		DNI_LETTER_MIN_RANGE = minDni[0]
		DNI_LETTER_MAX_RANGE = maxDni[0]
	}

	if zip != "" {
		ZIP_CODES = append(ZIP_CODES, zip)
	} else {
		switch strings.ToLower(province) {
		case "all", "catalunya", "cataluña":
			ZIP_CODES = append(ZIP_CODES, o1st.BARCELONA...)
			ZIP_CODES = append(ZIP_CODES, o1st.GIRONA...)
			ZIP_CODES = append(ZIP_CODES, o1st.LLEIDA...)
			ZIP_CODES = append(ZIP_CODES, o1st.TARRAGONA...)
		case "barcelona":
			ZIP_CODES = o1st.BARCELONA
		case "girona", "gerona":
			ZIP_CODES = o1st.GIRONA
		case "lleida", "lerida", "lérida":
			ZIP_CODES = o1st.LLEIDA
		case "tarragona":
			ZIP_CODES = o1st.TARRAGONA
		default:
			panic(errors.New("Unknown province " + province))
		}
	}

	if year != 0 {
		MIN_YEAR = year
		MAX_YEAR = year + 1
	}
	// Read the database
	o1st.ReadData("db")
}

func pad2(number int64) string {
	if number >= 10 {
		return strconv.FormatInt(number, 10)
	} else {
		return "0" + strconv.FormatInt(number, 10)
	}
}
func pad5(number int64) string {
	if number < 1e1 {
		return "0000" + strconv.FormatInt(number, 10)
	} else if number < 1e2 {
		return "000" + strconv.FormatInt(number, 10)
	} else if number < 1e3 {
		return "00" + strconv.FormatInt(number, 10)
	} else if number < 1e4 {
		return "0" + strconv.FormatInt(number, 10)
	} else {
		return strconv.FormatInt(number, 10)
	}
}

func nextDni(dniN int, dniCc uint8) (string, int, uint8) {
    if dniCc < DNI_LETTER_MAX_RANGE {
        dniCc += 1
    } else if dniN < DNI_NUM_MAX_RANGE {
        dniCc = DNI_LETTER_MIN_RANGE
        dniN += 1
    } else {
        return "", dniN, dniCc
    }
    return pad5(int64(dniN)) + string(dniCc), dniN, dniCc
}

func nextZip(zipI int) (string, int) {
	zipI += 1
    if zipI < len(ZIP_CODES) {
        return ZIP_CODES[zipI], zipI
    } else {
        return "", zipI
    }
}

func nextDate(year, month, day int) (string, int, int, int) {
    day += 1;
    if day == 32 {
        month += 1
        day = 1
    }
    if month == 13 {
        year += 1
        month = 1
    }
    if year == MAX_YEAR {
        return "", year, month, day
    } else {
        return strconv.FormatInt(int64(year), 10) + pad2(int64(month)) + pad2(int64(day)), year, month, day
    }
}

func outputFinding(info, dni, date, zip string) {
	fmt.Println("#################################")
	fmt.Println("# FOUND !")
	fmt.Println("# DNI:", dni)
	fmt.Println("# DATE:", date)
	fmt.Println("# ZIP:", zip)
	fmt.Println("# INFO:", info)
	fmt.Println("#################################")

}

func loopZip(dni, date string) {
	// fmt.Println(date)
	var info string
	zip, zipI := nextZip(-1)
	for zip != "" {
		info = o1st.Check(dni, date, zip);
		if info != "" {
			outputFinding(info, dni, date, zip)
		}
		zip, zipI = nextZip(zipI)
	}
}

func loopDates(dni string) {
	date, year, month, day := nextDate(MIN_YEAR - 1, 12, 31)
	var wg sync.WaitGroup

	for date != "" {
		wg.Add(1)
		go func(dni, date string) {
			defer wg.Done()
			loopZip(dni, date)
		}(dni, date)
		date, year, month, day = nextDate(year, month, day)
	}
	wg.Wait()
}

func loopDnis() {
	dni, dniN, dniCc := nextDni(DNI_NUM_MIN_RANGE - 1, DNI_LETTER_MAX_RANGE)
	// Loop through DNIs
	for dni != "" {
		fmt.Println("Checking DNI", dni);
		loopDates(dni)
		dni, dniN, dniCc = nextDni(dniN, dniCc)
	}
}

func main() {
	// Parse arguments
	flag.Parse()
	// Process arguments
	processArguments()
	// Run
	loopDnis()
}

