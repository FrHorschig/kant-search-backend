package model

type WorkMetadata struct {
	Id           int32
	Code         string
	Abbreviation string
	Volume       int32
}

// TODO frhorschig: find a better solution than hardcoding this data here
var Metadata = []WorkMetadata{
	// Band 1
	{Code: "GSK", Abbreviation: "GSK", Volume: 1, Id: 0},
	{Code: "UFE", Abbreviation: "UFE", Volume: 1, Id: 1},
	{Code: "FE", Abbreviation: "FE", Volume: 1, Id: 2},
	{Code: "NTH", Abbreviation: "NTH", Volume: 1, Id: 3},
	{Code: "Di", Abbreviation: "Di", Volume: 1, Id: 4},
	{Code: "PND", Abbreviation: "PND", Volume: 1, Id: 5},
	{Code: "VUE", Abbreviation: "VUE", Volume: 1, Id: 6},
	{Code: "GNVE", Abbreviation: "GNVE", Volume: 1, Id: 7},
	{Code: "FBZE", Abbreviation: "FBZE", Volume: 1, Id: 8},
	{Code: "MON_PH", Abbreviation: "MonPh", Volume: 1, Id: 9},
	{Code: "TW", Abbreviation: "TW", Volume: 1, Id: 10},
	{Code: "INTRO_1", Abbreviation: "", Volume: 1, Id: 11},

	// Band 2
	{Code: "EACG", Abbreviation: "EACG", Volume: 2, Id: 0},
	{Code: "NLBR", Abbreviation: "NLBR", Volume: 2, Id: 1},
	{Code: "VBO", Abbreviation: "VBO", Volume: 2, Id: 2},
	{Code: "GAJFF", Abbreviation: "GAJFF", Volume: 2, Id: 3},
	{Code: "DFS", Abbreviation: "DfS", Volume: 2, Id: 4},
	{Code: "BDG", Abbreviation: "BDG", Volume: 2, Id: 5},
	{Code: "NG", Abbreviation: "NG", Volume: 2, Id: 6},
	{Code: "GSE", Abbreviation: "GSE", Volume: 2, Id: 7},
	{Code: "VKK", Abbreviation: "VKK", Volume: 2, Id: 8},
	{Code: "REZ_SILBERSCHLAG_2", Abbreviation: "", Volume: 2, Id: 9},
	{Code: "UD", Abbreviation: "UD", Volume: 2, Id: 10},
	{Code: "NEV", Abbreviation: "NEV", Volume: 2, Id: 11},
	{Code: "TG", Abbreviation: "TG", Volume: 2, Id: 12},
	{Code: "GUGR", Abbreviation: "GUGR", Volume: 2, Id: 13},
	{Code: "MSI", Abbreviation: "MSI", Volume: 2, Id: 14},
	{Code: "REZ_MOSCATI", Abbreviation: "RezMoscati", Volume: 2, Id: 15},
	{Code: "VVRM", Abbreviation: "VvRM", Volume: 2, Id: 16},
	{Code: "AP", Abbreviation: "AP", Volume: 2, Id: 17},

	// Band 3
	{Code: "KRV_B", Abbreviation: "KrV B", Volume: 3, Id: 0},

	// Band 4
	{Code: "KRV_A", Abbreviation: "KrV A", Volume: 4, Id: 0},
	{Code: "PROL", Abbreviation: "Prol", Volume: 4, Id: 1},
	{Code: "GMS", Abbreviation: "GMS", Volume: 4, Id: 2},
	{Code: "MAN", Abbreviation: "MAN", Volume: 4, Id: 3},

	// Band 5
	{Code: "KPV", Abbreviation: "KpV", Volume: 5, Id: 0},
	{Code: "KU", Abbreviation: "KU", Volume: 5, Id: 1},

	// Band 6
	{Code: "RGV", Abbreviation: "RGV", Volume: 6, Id: 0},
	{Code: "MS", Abbreviation: "MS", Volume: 6, Id: 1},

	// Band 7
	{Code: "SF", Abbreviation: "SF", Volume: 7, Id: 0},
	{Code: "ANTH", Abbreviation: "Anth", Volume: 7, Id: 1},

	// Band 8
	{Code: "LAMBERT_BRIEFWECHSEL", Abbreviation: "", Volume: 8, Id: 0},
	{Code: "NACHRICHT_AERZTE", Abbreviation: "", Volume: 8, Id: 1},
	{Code: "REZ_SCHULZ", Abbreviation: "RezSchulz", Volume: 8, Id: 2},
	{Code: "IDEE_GESCHICHTE", Abbreviation: "", Volume: 8, Id: 3},
	{Code: "FRAGE_AUFKLAERUNG", Abbreviation: "", Volume: 8, Id: 4},
	{Code: "REZ_HERDER", Abbreviation: "RezHerder", Volume: 8, Id: 5},
	{Code: "VULKANE_MOND", Abbreviation: "", Volume: 8, Id: 6},
	{Code: "VUB", Abbreviation: "VUB", Volume: 8, Id: 7},
	{Code: "BEGRIFF_MENSCHENRACE", Abbreviation: "", Volume: 8, Id: 8},
	{Code: "ANFANG_MENSCHENGESCHICHTE", Abbreviation: "", Volume: 8, Id: 9},
	{Code: "REZ_HUFELAND", Abbreviation: "RezHufeland", Volume: 8, Id: 10},
	{Code: "WDO", Abbreviation: "WDO", Volume: 8, Id: 11},
	{Code: "BEM_MORGENSTUNDEN", Abbreviation: "", Volume: 8, Id: 12},
	{Code: "UEGTP", Abbreviation: "ÜGTP", Volume: 8, Id: 13},
	{Code: "UEE", Abbreviation: "ÜE", Volume: 8, Id: 14},
	{Code: "MISSLINGEN_THEODICEE", Abbreviation: "", Volume: 8, Id: 15},
	{Code: "TP", Abbreviation: "TP", Volume: 8, Id: 16},
	{Code: "MOND_WITTERUNG", Abbreviation: "", Volume: 8, Id: 17},
	{Code: "EAD", Abbreviation: "", Volume: 8, Id: 18},
	{Code: "ZEF", Abbreviation: "ZeF", Volume: 8, Id: 19},
	{Code: "VT", Abbreviation: "VT", Volume: 8, Id: 20},
	{Code: "AUSGLEICH_STREIT", Abbreviation: "", Volume: 8, Id: 21},
	{Code: "VNAEF", Abbreviation: "VNAEF", Volume: 8, Id: 22},
	{Code: "VRML", Abbreviation: "VRML", Volume: 8, Id: 23},
	{Code: "BUCHMACHEREI", Abbreviation: "", Volume: 8, Id: 24},
	{Code: "VORREDE_REL_PHIL", Abbreviation: "", Volume: 8, Id: 25},
	{Code: "NACHSCHRIFT_WOERTERBUCH", Abbreviation: "", Volume: 8, Id: 26},
	{Code: "NACHTRAG_8", Abbreviation: "", Volume: 8, Id: 27},
	{Code: "REZ_SILBERSCHLAG_8", Abbreviation: "", Volume: 8, Id: 28},
	{Code: "ANHANG_8", Abbreviation: "", Volume: 8, Id: 29},
	{Code: "REZ_ULRICH", Abbreviation: "RezUlrich", Volume: 8, Id: 30},

	// Band 9
	{Code: "LOG", Abbreviation: "Log", Volume: 9, Id: 0},
	{Code: "PG", Abbreviation: "PG", Volume: 9, Id: 1},
	{Code: "PAED", Abbreviation: "Päd", Volume: 9, Id: 2},
}
