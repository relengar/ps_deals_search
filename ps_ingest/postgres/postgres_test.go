package postgres

import (
	datatypes "ps_ingest/dataTypes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleEmbedding = []float64{0.012421737425029278, 0.009654982946813107, 0.02047879435122013, -0.054666440933942795, 0.04728856310248375, 0.06023680046200752, -0.04370197653770447, -0.0037065569777041674, 0.0034301141276955605, 0.02721296064555645, 0.034117430448532104, -0.07156721502542496, 0.013467920944094658, -0.0033158620353788137, 0.08989061415195465, -0.035150643438100815, 0.013203272596001625, 0.038510844111442566, 0.030942946672439575, -0.0735090970993042, 0.0004520770162343979, -0.03533260151743889, 0.0013799309963360429, 0.029630044475197792, -0.08792532980442047, 0.03165394067764282, -0.008342149667441845, 0.03356732055544853, 0.017092427238821983, -0.046525876969099045, 0.06603241711854935, 0.019502734765410423, -0.1225460097193718, 0.01030975766479969, 0.0013571494491770864, -0.06212208792567253, 0.011766523122787476, 0.0529390312731266, -0.040894653648138046, -0.04647260904312134, -0.04729074612259865, -0.004078480415046215, 0.003174461657181382, 0.1258150041103363, 0.11542364954948425, 0.023734452202916145, -0.08049371838569641, -0.06309617310762405, -0.0018515995470806956, 0.036001063883304596, 0.051692306995391846, -0.02154192142188549, 0.051280777901411057, 0.056036576628685, 0.09989657253026962, -0.00476260669529438, -0.043253686279058456, 0.018322085961699486, 0.10192409157752991, 0.00007272499351529405, 0.03433256596326828, -0.13866637647151947, -0.06123381853103638, -0.012047936208546162, -0.04659326374530792, -0.0055207423865795135, 0.047587230801582336, 0.007091716397553682, -0.009806780144572258, -0.05507146939635277, -0.09450715035200119, 0.003104227827861905, 0.08471163362264633, 0.019233407452702522, -0.04915211722254753, 0.029766250401735306, 0.05958031117916107, -0.03625922277569771, -0.037360887974500656, 0.010779410600662231, 0.013197309337556362, -0.04584004357457161, -0.09639465808868408, -0.08398380875587463, 0.04607284069061279, 0.011784467846155167, 0.03296612575650215, -0.06425067037343979, -0.021474679931998253, -0.004263280890882015, 0.001816318603232503, 0.008628207258880138, 0.06701745092868805, -0.03306353837251663, -0.07032701373100281, 0.03147941082715988, -0.08995728939771652, -0.005539973732084036, 0.008544083684682846, 0.0359005406498909, -0.016435444355010986, -0.0007261380669660866, 0.0566241480410099, 0.038081757724285126, 0.027058972045779228, -0.04760665446519852, 0.004372147843241692, 0.05650264769792557, 0.00121390912681818, 0.014046943746507168, -0.05478218197822571, 0.04152733460068703, -0.10004770755767822, 0.001066996599547565, 0.030594786629080772, 0.1116006150841713, -0.024873336777091026, 0.045124102383852005, 0.13720040023326874, -0.030241290107369423, 0.03483397886157036, -0.008312884718179703, 0.07161790132522583, 0.009597871452569962, -0.06358199566602707, -0.07070786505937576, 0.02377004362642765, 2.6289520192295153e-34, -0.0661875531077385, 0.018839212134480476, -0.029507320374250412, -0.022586416453123093, -0.04897892102599144, -0.03344827517867088, 0.03785064071416855, -0.0016796654090285301, -0.019679663702845573, 0.029065480455756187, -0.05241955816745758, 0.03289258852601051, -0.021221444010734558, 0.062288373708724976, 0.0494382381439209, -0.02154603786766529, -0.018689250573515892, -0.025904513895511627, 0.0025256217923015356, 0.0882740393280983, 0.013384878635406494, 0.06340386718511581, 0.04912685975432396, -0.003762702690437436, 0.01305385585874319, 0.06461197882890701, -0.021973559632897377, -0.06217455118894577, 0.11757700890302658, 0.048901431262493134, -0.10832773149013519, -0.065951868891716, 0.030330045148730278, -0.007924783043563366, 0.07699497044086456, -0.01373019628226757, -0.004008245188742876, -0.0407269187271595, 0.01080675981938839, 0.05102575942873955, 0.022099019959568977, -0.031229671090841293, -0.136829674243927, -0.039338063448667526, 0.004921028856188059, -0.05050288513302803, 0.028907835483551025, 0.006517417728900909, -0.03196919336915016, 0.09065984934568405, -0.05367995426058769, 0.006213001906871796, -0.0829385295510292, -0.07353091984987259, 0.0772358775138855, -0.06725496798753738, -0.04620254784822464, -0.01565881073474884, 0.06556228548288345, 0.042628396302461624, 0.007572839967906475, 0.007081720046699047, 0.010579456575214863, 0.0075623332522809505, -0.021176280453801155, 0.021457603201270103, 0.07553278654813766, -0.015661010518670082, -0.022081362083554268, -0.04972729831933975, -0.05806469917297363, 0.028566304594278336, 0.03516770526766777, 0.039363663643598557, -0.03835119679570198, -0.022713638842105865, -0.07514191418886185, 0.01340540498495102, 0.02573302946984768, 0.09884879738092422, -0.10178431123495102, 0.07861462980508804, 0.003491437528282404, 0.03944285959005356, 0.036009158939123154, 0.03584061190485954, 0.01059073954820633, -0.01457689143717289, 0.013810665346682072, -0.024893201887607574, -0.02823825553059578, -0.08550281822681427, -0.01731104589998722, 0.049191735684871674, 0.01965275965631008, -1.7168291521721287e-33, 0.013596777804195881, -0.009848459623754025, 0.013642171397805214, 0.011126548983156681, 0.014110278338193893, -0.010964805260300636, -0.00963451899588108, 0.04084393009543419, -0.006227591540664434, 0.0400388166308403, 0.06731162220239639, 0.03219141066074371, -0.06255465745925903, 0.018932940438389778, 0.09608536958694458, -0.04004771262407303, -0.012747244909405708, 0.01905861124396324, 0.05420773848891258, -0.01827353984117508, 0.07606122642755508, 0.04862707108259201, 0.04102201759815216, 0.08423595130443573, 0.05734824761748314, 0.04022079333662987, 0.0766567587852478, 0.08118832111358643, 0.036780185997486115, 0.01506213191896677, 0.016330352053046227, -0.005059749353677034, 0.044171933084726334, 0.04940033704042435, -0.05255935341119766, -0.02376597560942173, 0.06823404878377914, 0.0083847651258111, -0.0004654281947296113, -0.01915566436946392, -0.023062199354171753, 0.08177348971366882, -0.044566962867975235, 0.10182110965251923, -0.09166985750198364, 0.01926034316420555, 0.009661141782999039, 0.024226225912570953, 0.13435201346874237, -0.0848093181848526, 0.029945682734251022, -0.004276324063539505, -0.04074392467737198, -0.056655414402484894, -0.09768948704004288, -0.07542752474546432, -0.03639267757534981, 0.0806342214345932, 0.024542439728975296, -0.0018871640786528587, 0.0749296247959137, -0.09518462419509888, -0.01578442007303238, -0.006581583060324192, 0.02134011499583721, 0.021987635642290115, 0.052206479012966156, -0.06634269654750824, -0.0651886910200119, 0.08495290577411652, -0.11334484815597534, -0.01487865298986435, -0.05188785865902901, -0.07071381062269211, -0.03797251731157303, 0.028768563643097878, -0.024139489978551865, 0.07199695706367493, 0.08454675227403641, -0.02799566648900509, -0.02739885076880455, 0.07843151688575745, 0.027425533160567284, 0.001954078208655119, 0.03758733719587326, -0.09760979562997818, -0.04553079232573509, 0.02687402069568634, -0.05410371720790863, -0.11048538237810135, -0.0249054916203022, -0.02061118558049202, -0.04012894630432129, 0.07412160187959671, -0.0358012355864048, -2.2364233132066147e-8, 0.0805424302816391, -0.018298866227269173, -0.0627867579460144, 0.006415423471480608, 0.00601272052153945, 0.05062832310795784, -0.05060097947716713, -0.07494290173053741, 0.06763923913240433, 0.06297601759433746, -0.0002334269811399281, 0.01285359263420105, 0.046334363520145416, -0.031531888991594315, -0.007119838614016771, -0.025672558695077896, -0.06742169708013535, 0.04429321736097336, -0.01153495255857706, -0.011199289932847023, 0.08373389393091202, -0.012787955813109875, -0.02231961488723755, -0.058773573487997055, 0.0645093023777008, -0.026415850967168808, -0.03362078592181206, -0.005766150075942278, 0.04564327001571655, -0.0882694199681282, 0.0692102238535881, -0.08000500500202179, -0.06449983268976212, -0.05873788893222809, 0.0507621131837368, -0.004360276274383068, 0.0040720500983297825, 0.03477080538868904, 0.04733099415898323, -0.00181919417809695, 0.019377904012799263, -0.07495973259210587, -0.02329442650079727, 0.02033570595085621, -0.0565473809838295, -0.02333110384643078, -0.02248743362724781, -0.12225466966629028, -0.0694926530122757, -0.08372706919908524, -0.02591049298644066, 0.025006704032421112, -0.05300772190093994, -0.043432071805000305, 0.010332023724913597, 0.03369209170341492, 0.023540474474430084, -0.0352490060031414, -0.0021287056151777506, -0.05122355744242668, 0.04134267568588257, -0.08292125165462494, 0.036184292286634445, 0.026538077741861343}

func TestPostgresInsert(t *testing.T) {
	client := setup(t)

	data := []struct {
		game       datatypes.Game
		embeddings [][]float64
		shouldFail bool
	}{
		{
			game: datatypes.Game{
				Name:          "Batman",
				Description:   "Batman saves the day",
				Price:         2.3,
				OriginalPrice: 5.0,
				Rating:        1.2,
				RatingsSum:    10,
				Expiration:    time.Now(),
				URL:           "http://batmangame.com",
			},
			embeddings: [][]float64{sampleEmbedding},
			shouldFail: false,
		},
	}
	for _, v := range data {
		err := client.InsertGame(v.game, v.embeddings)
		if v.shouldFail {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func setup(t *testing.T) PgClient {
	client := client{"postgres", "supersecure", "localhost", "postgres", nil}

	err := client.Connect()
	require.Nil(t, err)

	return &client
}