package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	ts := []struct {
		n       string
		m, y, s string
	}{
		{"[ www.Speed.cd ] -Sons.of.Anarchy.S07E07.720p.HDTV.X264-DIMENSION", "Sons Of Anarchy S07E07", "", "07"},
		{"[@Difilm] The.Hot.Spot.1990.480p.BluRay.HardSub", "The Hot Spot", "1990", ""},
		{"[@MovieSpecial] Wild.Things.1998.BRRip.HardSub", "Wild Things", "1998", ""},
		{"[720pMkv.Com]_sons.of.anarchy.s05e10.480p.BluRay.x264-GAnGSteR", "Sons Of Anarchy S05e10", "", "05"},
		{"🃏@film_night🃏Venus in Fur 2013 BluRay 720p", "Venus In Fur", "2013", ""},
		{"2047 - Sights of Death (2014) 720p BrRip x264 - YIFY", "2047 Sights Of Death", "2014", ""},
		{"22 Jump Street (2014) 720p BrRip x264 - YIFY", "22 Jump Street", "2014", ""},
		{"9.Songs.2004.720p.BluRay.HardSub.Digimoviez", "9 Songs", "2004", ""},
		{"Akira (2016) - UpScaled - 720p - DesiSCR-Rip - Hindi - x264 - AC3 - 5.1 - Mafiaking - M2Tv", "Akira", "2016", ""},
		{"American.Gods.S01E01.1080p.HEVC.x265-MeGusta", "American Gods S01E01", "", "01"},
		{"american.gods.s01e02.1080p.webrip.hevc.x265-rmteam", "American Gods S01e02", "", "01"},
		{"Annabelle.2014.1080p.PROPER.HC.WEBRip.x264.AAC.2.0-RARBG", "Annabelle", "2014", ""},
		{"Annabelle.2014.HC.HDRip.XViD.AC3-juggs[ETRG]", "Annabelle", "2014", ""},
		{"Ant-Man.2015.3D.1080p.BRRip.Half-SBS.x264.AAC-m2g", "Ant-Man", "2015", ""},
		{"Ben Hur 2016 TELESYNC x264 AC3 MAXPRO", "Ben Hur", "2016", ""},
		{"Bliss.1997.DVDRip.HardSub", "Bliss", "1997", ""},
		{"Brave.2012.R5.DVDRip.XViD.LiNE-UNiQUE", "Brave", "2012", ""},
		{"breaking.bad.s01e01.720p.bluray.x264-reward", "Breaking Bad S01e01", "", "01"},
		{"Caníbal.2013.BluRay.720p.HardSub", "Caníbal", "2013", ""},
		{"Community.s02e20.rus.eng.720p.Kybik.v.Kybe", "Community S02e20", "", "02"},
		{"Dawn.Of.The.Planet.of.The.Apes.2014.1080p.WEB-DL.DD51.H264-RARBG", "Dawn Of The Planet Of The Apes", "2014", ""},
		{"Dawn.of.the.Planet.of.the.Apes.2014.HDRip.XViD-EVO", "Dawn Of The Planet Of The Apes", "2014", ""},
		{"Die.Marquise.von.Sade.1976.720p.BluRay.HardSub.Digimoviez", "Die Marquise Von Sade", "1976", ""},
		{"Dinosaur 13 2014 WEBrip XviD AC3 MiLLENiUM", "Dinosaur 13", "2014", ""},
		{"Doctor.Who.2005.8x11.Dark.Water.720p.HDTV.x264-FoV[rartv]", "Doctor Who 8x11 Dark Water", "2005", "08"},
		// {"doctor_who_2005.8x12.death_in_heaven.720p_hdtv_x264-fov", "Doctor Who 8x12 Death In Heaven", "2005", "08"},
		{"Double.Lover.2017.720p.BluRay.HardSub.Digimoviez", "Double Lover", "2017", ""},
		{"Downton Abbey 5x06 HDTV x264-FoV [eztv]", "Downton Abbey 5x06", "", "05"},
		{"Dracula.Untold.2014.TS.XViD.AC3.MrSeeN-SiMPLE", "Dracula Untold", "2014", ""},
		{"Eliza Graves (2014) Dual Audio WEB-DL 720p MKV x264", "Eliza Graves", "2014", ""},
		{"Femme.Fatale.2002.720p.BluRay.HardSub.mp4", "Femme Fatale", "2002", ""},
		{"Game of Thrones - 4x03 - Breaker of Chains", "Game Of Thrones 4x03 Breaker Of Chains", "", "04"},
		{"Girl House (2015) BluRay 720p-hardsub-(@GalleryMovies)", "Girl House", "2015", ""},
		{"Gotham.S01E05.Viper.WEB-DL.x264.AAC", "Gotham S01E05 Viper", "", "01"},
		{"Gotham.S01E07.Penguins.Umbrella.WEB-DL.x264.AAC", "Gotham S01E07 Penguins Umbrella", "", "01"},
		{"Guardians of the Galaxy (2014) Dual Audio DVDRip AVI", "Guardians Of The Galaxy", "2014", ""},
		{"Guardians Of The Galaxy 2014 R6 720p HDCAM x264-JYK", "Guardians Of The Galaxy", "2014", ""},
		{"Guardians of the Galaxy (CamRip - 2014)", "Guardians Of The Galaxy", "2014", ""},
		{"Halt.and.Catch.Fire.S04E02.Signal.to.Noise.1080p.AMZN.WEBRip.DDP5.1.x264-NTb[rarbg]", "Halt And Catch Fire S04E02 Signal To Noise", "", "04"},
		{"Halt.and.Catch.Fire.S04E06.CONVERT.1080p.WEB.h264-TBS[rarbg]", "Halt And Catch Fire S04E06", "", "04"},
		{"Halt.and.Catch.Fire.S04E10.1080p.WEB.H264-STRiFE[rarbg]", "Halt And Catch Fire S04E10", "", "04"},
		{"Hercules (2014) 1080p BrRip H264 - YIFY", "Hercules", "2014", ""},
		{"Hercules.2014.EXTENDED.1080p.WEB-DL.DD5.1.H264-RARBG", "Hercules", "2014", ""},
		{"Hercules.2014.Extended.Cut.HDRip.XViD-juggs[ETRG]", "Hercules", "2014", ""},
		{"Hercules (2014) WEBDL DVDRip XviD-MAX", "Hercules", "2014", ""},
		{"Hes.Just.Not.That.Into.You.2009,[@Intermedia]", "Hes Just Not That Into You", "2009", ""},
		{"Ice.Age.Collision.Course.2016.READNFO.720p.HDRIP.X264.AC3.TiTAN", "Ice Age Collision Course", "2016", ""},
		{"Interstellar (2014) CAM ENG x264 AAC-CPG", "Interstellar", "2014", ""},
		{"Into The Storm 2014 1080p BRRip x264 DTS-JYK", "Into The Storm", "2014", ""},
		{"Into.The.Storm.2014.1080p.WEB-DL.AAC2.0.H264-RARBG", "Into The Storm", "2014", ""},
		{"Its.Always.Sunny.In.Philadelphia.S05E02.BDRip", "Its Always Sunny In Philadelphia S05E02", "", "05"},
		{"Jack.And.The.Cuckoo-Clock.Heart.2013.BRRip XViD", "Jack And The Cuckoo-Clock Heart", "2013", ""},
		{"Last.Tango.in.Paris.1972.720p.BluRay.HardSub", "Last Tango In Paris", "1972", ""},
		{"Lets.Be.Cops.2014.BRRip.XViD-juggs[ETRG]", "Lets Be Cops", "2014", ""},
		{"Lovelace.2013.720p.BluRay-@TheMovieShare", "Lovelace", "2013", ""},
		{"Lucy 2014 Dual-Audio 720p WEBRip", "Lucy", "2014", ""},
		{"Lucy 2014 Dual-Audio WEBRip 1400Mb", "Lucy", "2014", ""},
		{"Lucy.2014.HC.HDRip.XViD-juggs[ETRG]", "Lucy", "2014", ""},
		{"Malizia.1973.480p.perSub", "Malizia", "1973", ""},
		{"Marvel's.Agents.of.S.H.I.E.L.D.S02E01.Shadows.1080p.WEB-DL.DD5.1", "Marvel'S Agents Of S H I E L D S02E01 Shadows", "", "02"},
		{"Marvels Agents of S H I E L D S02E05 HDTV x264-KILLERS [eztv]", "Marvels Agents Of S H I E L D S02E05", "", "02"},
		{"Marvels Agents of S.H.I.E.L.D. S02E06 HDTV x264-KILLERS[ettv]", "Marvels Agents Of S.H.I.E.L.D. S02E06", "", "02"},
		{"Match_Point_2005_hardsub", "Match Point", "2005", ""},
		{"Mektoub.My.Love.Canto.Uno.2017.720p.HardSub", "Mektoub My Love Canto Uno", "2017", ""},
		{"One Shot [2014] DVDRip XViD-ViCKY", "One Shot", "2014", ""},
		{"Red.Sonja.Queen.Of.Plagues.2016.BDRip.x264-W4F[PRiME]", "Red Sonja Queen Of Plagues", "2016", ""},
		{"Return.To.Snowy.River.1988.iNTERNAL.DVDRip.x264-W4F[PRiME]", "Return To Snowy River", "1988", ""},
		{"rick.and.morty.s03e01.720p.hdtv.x264-w4f", "Rick And Morty S03e01", "", "03"},
		{"Silicon.Valley.S04E04.1080p.WEB.h264-TBS", "Silicon Valley S04E04", "", "04"},
		{"Sin.City.A.Dame.to.Kill.For.2014.1080p.BluRay.x264-SPARKS", "Sin City A Dame To Kill For", "2014", ""},
		{"Sons.of.Anarchy.S01E03", "Sons Of Anarchy S01E03", "", "01"},
		{"South Park S18E05 HDTV x264-KILLERS [eztv]", "South Park S18E05", "", "18"},
		{"Teenage.Mutant.Ninja.Turtles.2014.720p.HDRip.x264.AC3.5.1-RARBG", "Teenage Mutant Ninja Turtles", "2014", ""},
		{"Teenage.Mutant.Ninja.Turtles.2014.HDRip.XviD.MP3-RARBG", "Teenage Mutant Ninja Turtles", "2014", ""},
		{"Teenage Mutant Ninja Turtles (HdRip - 2014)", "Teenage Mutant Ninja Turtles", "2014", ""},
		{"Teenage Mutant Ninja Turtles (unknown_release_type - 2014)", "Teenage Mutant Ninja Turtles", "2014", ""},
		{"Teeth_2007", "Teeth", "2007", ""},
		{"The Big Bang Theory S08E06 HDTV XviD-LOL [eztv]", "The Big Bang Theory S08E06", "", "08"},
		{"The.Boss.2016.UNRATED.720p.BRRip.x264.AAC-ETRG", "The Boss", "2016", ""},
		{"The.Dark.Side.of.the.Heart.DVDRip.HardSub", "The Dark Side Of The Heart", "", ""},
		{"The.Duke.of.Burgundy.2014.720p.BluRay.HardSub", "The Duke Of Burgundy", "2014", ""},
		{"The Flash 2014 S01E01 HDTV x264-LOL[ettv]", "The Flash S01E01", "2014", "01"},
		{"The Flash 2014 S01E03 HDTV x264-LOL[ettv]", "The Flash S01E03", "2014", "01"},
		{"The Flash 2014 S01E04 HDTV x264-FUM[ettv]", "The Flash S01E04", "2014", "01"},
		{"The Hateful Eight (2015) 720p BluRay - x265 HEVC - 999MB - ShAaN", "The Hateful Eight", "2015", ""},
		{"The.Jungle.Book.2016.3D.1080p.BRRip.SBS.x264.AAC-ETRG", "The Jungle Book", "2016", ""},
		{"The Missing 1x01 Pilot HDTV x264-FoV [eztv]", "The Missing 1x01 Pilot", "", "01"},
		{"The.Platform.2019.720p.WEB-DL.SoftSub", "The Platform", "2019", ""},
		{"The Purge: Election Year (2016) HC - 720p HDRiP - 900MB - ShAaNi", "The Purge: Election Year", "2016", ""},
		{"The.Secret.Life.of.Pets.2016.HDRiP.AAC-LC.x264-LEGi0N", "The Secret Life Of Pets", "2016", ""},
		{"These.Final.Hours.2013.WBBRip XViD", "These Final Hours", "2013", ""},
		{"The Shaukeens (2014) 1CD DvDScr Rip x264 [DDR]", "The Shaukeens", "2014", ""},
		{"The Shaukeens 2014 Hindi (1CD) DvDScr x264 AAC...Hon3y", "The Shaukeens", "2014", ""},
		{"The Simpsons S26E05 HDTV x264 PROPER-LOL [eztv]", "The Simpsons S26E05", "", "26"},
		{"The.Walking.Dead.S05E03.1080p.WEB-DL.DD5.1.H.264-Cyphanix[rartv]", "The Walking Dead S05E03", "", "05"},
		{"The Walking Dead S05E03 720p HDTV x264-ASAP[ettv]", "The Walking Dead S05E03", "", "05"},
		{"The.Wings.of.The.Dove.1997.720p.HardSub", "The Wings Of The Dove", "1997", ""},
		{"They.2017.WEBRip.1080p.YTS.Dream", "They", "2017", ""},
		{"Trainwreck", "Trainwreck", "", ""},
		{"Two and a Half Men S12E01 HDTV x264 REPACK-LOL [eztv]", "Two And A Half Men S12E01", "", "12"},
		{"UFC.179.PPV.HDTV.x264-Ebi[rartv]", "UFC 179", "", ""},
		{"War Dogs (2016) HDTS 600MB - NBY", "War Dogs", "2016", ""},
		{"Wild.Things.2.2004.720p.HardSub", "Wild Things 2", "2004", ""},
		{"WWE Hell in a Cell 2014 HDTV x264 SNHD", "WWE Hell In A Cell", "2014", ""},
		{"WWE Hell in a Cell 2014 PPV WEB-DL x264-WD -={SPARROW}=-", "WWE Hell In A Cell", "2014", ""},
		{"WWE Monday Night Raw 2014 11 10 WS PDTV x264-RKOFAN1990 -={SPARR", "WWE Monday Night Raw 11 10", "2014", ""},
		{"WWE Monday Night Raw 3rd Nov 2014 HDTV x264-Sir Paul", "WWE Monday Night Raw 3rd Nov", "2014", ""},
		{"www.torrenting.com - Silicon.Valley.S04E04.1080p.WEB.h264-TBS", "Silicon Valley S04E04", "", "04"},
		{"X-Men.Days.of.Future.Past.2014.1080p.WEB-DL.DD5.1.H264-RARBG", "X-Men Days Of Future Past", "2014", ""},
	}

	for _, tt := range ts {
		pf := &plexFile{name: tt.n, mov: movie{}}
		pf.parse()
		if pf.mov.name != tt.m {
			t.Errorf("name: %v\ngot:  %v\nwant: %v", tt.n, pf.mov.name, tt.m)
		}
		if pf.mov.year != tt.y {
			t.Errorf("year: %v\ngot:  %v\nwant: %v", tt.n, pf.mov.year, tt.y)
		}
		if pf.mov.season != tt.s {
			t.Errorf("season: %v\ngot:  %v\nwant: %v", tt.n, pf.mov.season, tt.s)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	pf := &plexFile{name: "Marvel's.Agents.of.S.H.I.E.L.D.S02E01.Shadows.1080p.WEB-DL.DD5.1", mov: movie{}}
	for i := 0; i < b.N; i++ {
		pf.parse()
	}
}
