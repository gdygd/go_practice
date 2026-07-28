package main

import (
	"fmt"
	"time"

	"sigmap_convert/dbapp"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
)

func main() {
	readsignalmap3()
	fmt.Println("모든 작업이 완료되었습니다.")
}

func readsignalmap3() {
	// 6-457
	// map : 1-2

	for id := 7; id <= 457; id++ {

		for bank := 1; bank <= 2; bank++ {
			time.Sleep(time.Millisecond * 50)
			mintm, err := dbapp.CheckMintime(id, bank)
			if mintm == 0 || err != nil {
				fmt.Printf("CheckMintime err : %v, id:%d, bank:%d, mintm : %d\n ", err, id, bank, mintm)
				continue
			}

			sigmaps, err := dbapp.ReadSignalMap2(id, bank)
			if err != nil {
				fmt.Printf("ReadSignalMap2 err : %v, %d, %d", err, id, bank)
				continue
			}

			var dt []byte = []byte{}

			for _, v := range sigmaps {
				var step []byte = make([]byte, 19)
				step[0] = byte(v.SI_SSR1_VHE)
				step[1] = byte(v.SI_SSR2_VHE)
				step[2] = byte(v.SI_SSR3_VHE)
				step[3] = byte(v.SI_SSR4_VHE)
				step[4] = byte(v.SI_SSR5_VHE)
				step[5] = byte(v.SI_SSR6_VHE)
				step[6] = byte(v.SI_SSR7_VHE)
				step[7] = byte(v.SI_SSR8_VHE)

				step[8] = byte(v.SI_SSR1_PED)
				step[9] = byte(v.SI_SSR2_PED)
				step[10] = byte(v.SI_SSR3_PED)
				step[11] = byte(v.SI_SSR4_PED)
				step[12] = byte(v.SI_SSR5_PED)
				step[13] = byte(v.SI_SSR6_PED)
				step[14] = byte(v.SI_SSR7_PED)
				step[15] = byte(v.SI_SSR8_PED)

				step[16] = byte(v.SI_MINTIME)
				step[17] = byte(v.SI_MAXTIME)
				step[18] = byte(v.SI_EOP)

				dt = append(dt, step...)
				// fmt.Printf("%d, %v \n", i+1, v)
			}

			// fmt.Printf("%d \n", len(dt))
			// for i := 0; i < 32*19*2; i += 19 {
			// 	fmt.Printf("%d, %v \n", (i/19)+1, dt[i:i+19])
			// }
			err = dbapp.InsertSignalmap(id, bank, dt)
			if err != nil {
				fmt.Printf("insert err : %v, %d, %d", err, id, bank)
				continue
			}

		}
	}

}

func readsignalmap2() {

	sigmaps, err := dbapp.ReadSignalMap2(427, 2)
	if err != nil {
		fmt.Printf("err : %v", err)
		return
	}

	var dt []byte = []byte{}

	for _, v := range sigmaps {
		var step []byte = make([]byte, 19)
		step[0] = byte(v.SI_SSR1_VHE)
		step[1] = byte(v.SI_SSR2_VHE)
		step[2] = byte(v.SI_SSR3_VHE)
		step[3] = byte(v.SI_SSR4_VHE)
		step[4] = byte(v.SI_SSR5_VHE)
		step[5] = byte(v.SI_SSR6_VHE)
		step[6] = byte(v.SI_SSR7_VHE)
		step[7] = byte(v.SI_SSR8_VHE)

		step[8] = byte(v.SI_SSR1_PED)
		step[9] = byte(v.SI_SSR2_PED)
		step[10] = byte(v.SI_SSR3_PED)
		step[11] = byte(v.SI_SSR4_PED)
		step[12] = byte(v.SI_SSR5_PED)
		step[13] = byte(v.SI_SSR6_PED)
		step[14] = byte(v.SI_SSR7_PED)
		step[15] = byte(v.SI_SSR8_PED)

		step[16] = byte(v.SI_MINTIME)
		step[17] = byte(v.SI_MAXTIME)
		step[18] = byte(v.SI_EOP)

		dt = append(dt, step...)
		// fmt.Printf("%d, %v \n", i+1, v)
	}

	// fmt.Printf("%d \n", len(dt))
	// for i := 0; i < 32*19*2; i += 19 {
	// 	fmt.Printf("%d, %v \n", (i/19)+1, dt[i:i+19])
	// }
	dbapp.InsertSignalmap(6, 1, dt)
}

func readsignalmap() {
	sigmaps, err := dbapp.ReadSignalMap()
	if err != nil {
		fmt.Printf("err : %v", err)
		return
	}
	// 1-410, 1-2
	var mpsigmap map[int]map[int]map[int]dbapp.Signalmap = make(map[int]map[int]map[int]dbapp.Signalmap)
	for _, v := range sigmaps {
		if mpsigmap[v.IT_NO] == nil {
			mpsigmap[v.IT_NO] = make(map[int]map[int]dbapp.Signalmap)
		}
		if mpsigmap[v.IT_NO][v.SI_MAPNO] == nil {
			mpsigmap[v.IT_NO][v.SI_MAPNO] = make(map[int]dbapp.Signalmap)
		}

		IT_NO := v.IT_NO
		MAPNO := v.SI_MAPNO
		RING := v.SI_RINGNO

		mpsigmap[IT_NO][MAPNO][RING] = v
		// fmt.Printf("%v \n", v)
	}

	for itno := range mpsigmap {
		for mapno := range mpsigmap {
			for r := 1; r <= 2; r++ {
				v := mpsigmap[itno][mapno][r]
				fmt.Printf("%d / %d / %d \n", itno, mapno, r)
				fmt.Printf("%v \n", v)
			}
		}
	}
}
