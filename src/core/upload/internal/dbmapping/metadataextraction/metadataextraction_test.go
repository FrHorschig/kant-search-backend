package metadataextraction

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/testutil"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/stretchr/testify/assert"
)

func TestMetadataExtraction(t *testing.T) {
	// both texts are based on an AI generated paragraph
	fmtText := "<ks-meta-line>1</ks-meta-line><ks-fmt-bold>Immanuel Kant</ks-fmt-bold> was an 18th-century German philosopher who shaped modern thought. <ks-meta-line>2</ks-meta-line>In his <ks-fmt-bold><ks-fmt-emph>Critique of Pure Reason</ks-fmt-emph></ks-fmt-bold>, he argued that knowledge arises from both experience and the mind’s structures. <ks-meta-page>18</ks-meta-page><ks-meta-line>1</ks-meta-line>His idea of the <ks-fmt-emph>categorical imperative</ks-fmt-emph> emphasized moral duty over outcomes. <ks-meta-page>19</ks-meta-page><ks-meta-line>1</ks-meta-line>Kant remains a <ks-fmt-bold>foundational figure</ks-fmt-bold> in ethics and epistemology."
	searchText := "Immanuel Kant was an 18th-century German philosopher who shaped modern thought. In his Critique of Pure Reason, he argued that knowledge arises from both experience and the mind’s structures. His idea of the categorical imperative emphasized moral duty over outcomes. Kant remains a foundational figure in ethics and epistemology."
	content := []dbmodel.Content{{
		FmtText:    fmtText,
		SearchText: searchText,
	}}
	expContent := []dbmodel.Content{{
		FmtText:    fmtText,
		SearchText: searchText,
		PageByIndex: []dbmodel.IndexNumberPair{
			{I: 333, Num: 18},
			{I: 497, Num: 19},
		},
		LineByIndex: []dbmodel.IndexNumberPair{
			{I: 0, Num: 1},
			{I: 137, Num: 2},
			{I: 364, Num: 1},
			{I: 528, Num: 1},
		},
		WordIndexMap: map[int32]int32{
			0:   43,
			9:   52,
			14:  71,
			18:  75,
			23:  80,
			26:  83,
			34:  91,
			41:  98,
			53:  110,
			57:  114,
			64:  121,
			71:  128,
			80:  167,
			83:  170,
			87:  200,
			96:  209,
			99:  212,
			104: 217,
			112: 253,
			115: 256,
			122: 263,
			127: 268,
			137: 278,
			144: 285,
			149: 290,
			154: 295,
			165: 306,
			169: 310,
			173: 314,
			178: 319,
			180: 321,
			192: 394,
			196: 398,
			201: 403,
			204: 406,
			208: 423,
			220: 435,
			231: 460,
			242: 471,
			248: 477,
			253: 482,
			258: 487,
			268: 558,
			273: 563,
			281: 571,
			283: 586,
			296: 599,
			303: 620,
			306: 623,
			313: 630,
			317: 634,
		},
	}}

	err := ExtractMetadata(content)
	assert.False(t, err.HasError)
	testutil.AssertDbContents(t, expContent, content)
}
