package main

import (
	"math"
	"math/rand"
	"tempo-loger/pkg/outlook"
	"testing"
	"time"
)

func TestGenerateRandom(t *testing.T) {
	for i := 0; i < 1; i++ {
		c := generateRandomInt(60 * 60 * 20)
		//time.Sleep(100)
		t.Logf("Generated: %v", c/60/60)
	}

}

func TestGenerateRandom1(t *testing.T) {
	neededSpentSeconds := 60 * 60 * 1
	min, max := 1, neededSpentSeconds
	rand.Seed(time.Date(2021, 12, 15, 18, 40, 00, 00, time.Now().Location()).UnixNano())
	t.Logf("Generated: %v", math.Round(float64(rand.Intn(max-min)+min)/60/60))
	t.Log("----------")
	rand.Seed(time.Date(2021, 12, 16, 18, 40, 00, 00, time.Now().Location()).UnixNano())
	t.Logf("Generated: %v", math.Round(float64(rand.Intn(max-min)+min)/60/60))
	t.Log("----------")
	rand.Seed(time.Date(2021, 12, 17, 18, 40, 00, 00, time.Now().Location()).UnixNano())
	t.Logf("Generated: %v", math.Round(float64(rand.Intn(max-min)+min)/60/60))
	t.Log("----------")
	rand.Seed(time.Date(2021, 12, 18, 18, 40, 00, 00, time.Now().Location()).UnixNano())
	t.Logf("Generated: %v", math.Round(float64(rand.Intn(max-min)+min)/60/60))
	t.Log("----------")
	rand.Seed(time.Date(2021, 12, 19, 18, 40, 00, 00, time.Now().Location()).UnixNano())
	t.Logf("Generated: %v", math.Round(float64(rand.Intn(max-min)+min)/60/60))
	t.Log("----------")
	rand.Seed(time.Date(2021, 12, 20, 18, 40, 00, 00, time.Now().Location()).UnixNano())
	t.Logf("Generated: %v", math.Round(float64(rand.Intn(max-min)+min)/60/60))
	t.Log("----------")
	rand.Seed(time.Date(2021, 12, 21, 18, 40, 00, 00, time.Now().Location()).UnixNano())
	t.Logf("Generated: %v", math.Round(float64(rand.Intn(max-min)+min)/60/60))
	t.Log("----------")
	rand.Seed(time.Date(2021, 12, 22, 18, 40, 00, 00, time.Now().Location()).UnixNano())
	t.Logf("Generated: %v", math.Round(float64(rand.Intn(max-min)+min)/60/60))
	t.Log("----------")
	rand.Seed(time.Date(2022, 01, 31, 18, 27, 00, 00, time.Now().Location()).UnixNano())
	t.Logf("Generated: %v", math.Round(float64(rand.Intn(max-min)+min)/60/60))

	//s1 := rand.NewSource(time.Now().UnixNano())
	//r1 := rand.New(s1)
	/*	if max-min <= 0 {
			t.Logf("Generated: %v", 1*60*60)
		} else {
			t.Logf("Generated: %v", (rand.Intn(max-min)+min)*60*60)
		}*/
}

func TestReqNtlm(t *testing.T) {
	o := outlook.New("i.nazmutdinov", "Naz01.19", "https://mail.stoloto.ru", "api/v2.0/me/calendarview")
	year, month, day := time.Now().Date()
	t.Log(year, month, day)
	rs, _ := o.GetEvents(time.Date(2022, time.February, 10, 00, 00, 00, 0000, time.UTC), time.Date(2022, time.February, 10, 23, 59, 59, 0, time.UTC))
	t.Log(rs.Value[0].Subject)
}
