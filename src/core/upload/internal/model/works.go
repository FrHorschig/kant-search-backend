package model

type WorkMetadata struct {
	Id           int32
	Code         string
	Abbreviation string
	Volume       int32
	Year         []string
}

// TODO frhorschig: find a better solution than hardcoding this data here
var Metadata = []WorkMetadata{
	// Band 1
	{Code: "GSK", Abbreviation: "GSK", Volume: 1, Id: 0, Year: []string{"1746", "1749"}},
	{Code: "UFE", Abbreviation: "UFE", Volume: 1, Id: 1, Year: []string{"1754"}},
	{Code: "FE", Abbreviation: "FE", Volume: 1, Id: 2, Year: []string{"1754"}},
	{Code: "NTH", Abbreviation: "NTH", Volume: 1, Id: 3, Year: []string{"1755"}},
	{Code: "Di", Abbreviation: "Di", Volume: 1, Id: 4, Year: []string{"1755"}},
	{Code: "PND", Abbreviation: "PND", Volume: 1, Id: 5, Year: []string{"1755"}},
	{Code: "VUE", Abbreviation: "VUE", Volume: 1, Id: 6, Year: []string{"1756"}},
	{Code: "GNVE", Abbreviation: "GNVE", Volume: 1, Id: 7, Year: []string{"1756"}},
	{Code: "FBZE", Abbreviation: "FBZE", Volume: 1, Id: 8, Year: []string{"1756"}},
	{Code: "MON_PH", Abbreviation: "MonPh", Volume: 1, Id: 9, Year: []string{"1756"}},
	{Code: "TW", Abbreviation: "TW", Volume: 1, Id: 10, Year: []string{"1756"}},
	{Code: "INTRO_1", Abbreviation: "", Volume: 1, Id: 11, Year: nil},

	// Band 2
	{Code: "EACG", Abbreviation: "EACG", Volume: 2, Id: 0, Year: []string{"1757"}},
	{Code: "NLBR", Abbreviation: "NLBR", Volume: 2, Id: 1, Year: []string{"1758"}},
	{Code: "VBO", Abbreviation: "VBO", Volume: 2, Id: 2, Year: []string{"1759"}},
	{Code: "GAJFF", Abbreviation: "GAJFF", Volume: 2, Id: 3, Year: []string{"1760"}},
	{Code: "DFS", Abbreviation: "DfS", Volume: 2, Id: 4, Year: []string{"1762"}},
	{Code: "BDG", Abbreviation: "BDG", Volume: 2, Id: 5, Year: []string{"1763"}},
	{Code: "NG", Abbreviation: "NG", Volume: 2, Id: 6, Year: []string{"1763"}},
	{Code: "GSE", Abbreviation: "GSE", Volume: 2, Id: 7, Year: []string{"1764"}},
	{Code: "VKK", Abbreviation: "VKK", Volume: 2, Id: 8, Year: []string{"1764"}},
	{Code: "REZ_SILBERSCHLAG_2", Abbreviation: "", Volume: 2, Id: 9, Year: []string{"1764"}},
	{Code: "UD", Abbreviation: "UD", Volume: 2, Id: 10, Year: []string{"1764"}},
	{Code: "NEV", Abbreviation: "NEV", Volume: 2, Id: 11, Year: []string{"1765"}},
	{Code: "TG", Abbreviation: "TG", Volume: 2, Id: 12, Year: []string{"1766"}},
	{Code: "GUGR", Abbreviation: "GUGR", Volume: 2, Id: 13, Year: []string{"1768"}},
	{Code: "MSI", Abbreviation: "MSI", Volume: 2, Id: 14, Year: []string{"1770"}},
	{Code: "REZ_MOSCATI", Abbreviation: "RezMoscati", Volume: 2, Id: 15, Year: []string{"1771"}},
	{Code: "VVRM", Abbreviation: "VvRM", Volume: 2, Id: 16, Year: []string{"1775"}},
	{Code: "AP", Abbreviation: "AP", Volume: 2, Id: 17, Year: []string{"1776", "1777"}},

	// Band 3
	{Code: "KRV_B", Abbreviation: "KrV B", Volume: 3, Id: 0, Year: []string{"1787"}},

	// Band 4
	{Code: "KRV_A", Abbreviation: "KrV A", Volume: 4, Id: 0, Year: []string{"1781"}},
	{Code: "PROL", Abbreviation: "Prol", Volume: 4, Id: 1, Year: []string{"1783"}},
	{Code: "GMS", Abbreviation: "GMS", Volume: 4, Id: 2, Year: []string{"1785"}},
	{Code: "MAN", Abbreviation: "MAN", Volume: 4, Id: 3, Year: []string{"1786"}},

	// Band 5
	{Code: "KPV", Abbreviation: "KpV", Volume: 5, Id: 0, Year: []string{"1788"}},
	{Code: "KU", Abbreviation: "KU", Volume: 5, Id: 1, Year: []string{"1790"}},

	// Band 6
	{Code: "RGV", Abbreviation: "RGV", Volume: 6, Id: 0, Year: []string{"1793"}},
	{Code: "MS", Abbreviation: "MS", Volume: 6, Id: 1, Year: []string{"1797"}},

	// Band 7
	{Code: "SF", Abbreviation: "SF", Volume: 7, Id: 0, Year: []string{"1798"}},
	{Code: "ANTH", Abbreviation: "Anth", Volume: 7, Id: 1, Year: []string{"1798"}},

	// Band 8
	{Code: "LAMBERT_BRIEFWECHSEL", Abbreviation: "", Volume: 8, Id: 0, Year: []string{"1782"}},
	{Code: "NACHRICHT_AERZTE", Abbreviation: "", Volume: 8, Id: 1, Year: []string{"1782"}},
	{Code: "REZ_SCHULZ", Abbreviation: "RezSchulz", Volume: 8, Id: 2, Year: []string{"1783", "1790"}},
	{Code: "IDEE_GESCHICHTE", Abbreviation: "", Volume: 8, Id: 3, Year: []string{"1784"}},
	{Code: "FRAGE_AUFKLAERUNG", Abbreviation: "", Volume: 8, Id: 4, Year: []string{"1784"}},
	{Code: "REZ_HERDER", Abbreviation: "RezHerder", Volume: 8, Id: 5, Year: []string{"1785"}},
	{Code: "VULKANE_MOND", Abbreviation: "", Volume: 8, Id: 6, Year: []string{"1785"}},
	{Code: "VUB", Abbreviation: "VUB", Volume: 8, Id: 7, Year: []string{"1785"}},
	{Code: "BEGRIFF_MENSCHENRACE", Abbreviation: "", Volume: 8, Id: 8, Year: []string{"1785"}},
	{Code: "ANFANG_MENSCHENGESCHICHTE", Abbreviation: "", Volume: 8, Id: 9, Year: []string{"1786"}},
	{Code: "REZ_HUFELAND", Abbreviation: "RezHufeland", Volume: 8, Id: 10, Year: []string{"1786"}},
	{Code: "WDO", Abbreviation: "WDO", Volume: 8, Id: 11, Year: []string{"1786"}},
	{Code: "BEM_MORGENSTUNDEN", Abbreviation: "", Volume: 8, Id: 12, Year: []string{"1786"}},
	{Code: "UEGTP", Abbreviation: "ÜGTP", Volume: 8, Id: 13, Year: []string{"1788"}},
	{Code: "UEE", Abbreviation: "ÜE", Volume: 8, Id: 14, Year: []string{"1790"}},
	{Code: "MISSLINGEN_THEODICEE", Abbreviation: "", Volume: 8, Id: 15, Year: []string{"1791"}},
	{Code: "TP", Abbreviation: "TP", Volume: 8, Id: 16, Year: []string{"1793"}},
	{Code: "MOND_WITTERUNG", Abbreviation: "", Volume: 8, Id: 17, Year: []string{"1794"}},
	{Code: "EAD", Abbreviation: "", Volume: 8, Id: 18, Year: []string{"1794"}},
	{Code: "ZEF", Abbreviation: "ZeF", Volume: 8, Id: 19, Year: []string{"1795"}},
	{Code: "VT", Abbreviation: "VT", Volume: 8, Id: 20, Year: []string{"1796"}},
	{Code: "AUSGLEICH_STREIT", Abbreviation: "", Volume: 8, Id: 21, Year: []string{"1796"}},
	{Code: "VNAEF", Abbreviation: "VNAEF", Volume: 8, Id: 22, Year: []string{"1796"}},
	{Code: "VRML", Abbreviation: "VRML", Volume: 8, Id: 23, Year: []string{"1797"}},
	{Code: "BUCHMACHEREI", Abbreviation: "", Volume: 8, Id: 24, Year: []string{"1798"}},
	{Code: "VORREDE_REL_PHIL", Abbreviation: "", Volume: 8, Id: 25, Year: []string{"1800"}},
	{Code: "NACHSCHRIFT_WOERTERBUCH", Abbreviation: "", Volume: 8, Id: 26, Year: []string{"1800"}},
	{Code: "NACHTRAG_8", Abbreviation: "", Volume: 8, Id: 27, Year: nil},
	{Code: "REZ_SILBERSCHLAG_8", Abbreviation: "", Volume: 8, Id: 28, Year: []string{"1764"}},
	{Code: "ANHANG_8", Abbreviation: "", Volume: 8, Id: 29, Year: nil},
	{Code: "REZ_ULRICH", Abbreviation: "RezUlrich", Volume: 8, Id: 30, Year: []string{"1788"}},

	// Band 9
	{Code: "LOG", Abbreviation: "Log", Volume: 9, Id: 0, Year: []string{"1800"}},
	{Code: "PG", Abbreviation: "PG", Volume: 9, Id: 1, Year: []string{"1802"}},
	{Code: "PAED", Abbreviation: "Päd", Volume: 9, Id: 2, Year: []string{"1803"}},
}
